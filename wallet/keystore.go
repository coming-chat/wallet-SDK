package wallet

import (
	"bytes"
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/coming-chat/merlin"
	"github.com/gtank/ristretto255"
	"golang.org/x/crypto/scrypt"
	"wallet-SDK/crypto"
	"wallet-SDK/u8util"
)

var (
	pkcs8Divider = []byte{161, 35, 3, 33, 0}
	pkcs8Header  = []byte{48, 83, 2, 1, 1, 48, 5, 6, 3, 43, 101, 112, 4, 34, 4, 32}
	seedOffset   = len(pkcs8Header)
	divOffset    = seedOffset + secLength
)

const (
	keyLength  = 32
	pubLength  = 32
	saltLength = 32
	secLength  = 64
	seedLength = 32

	scryptLength = 32 + (3 * 4)
	nonceLength  = 24

	defaultN int64 = 1 << 15
	defaultP int64 = 1
	defaultR int64 = 8
)

func decodePolkaKeystore(passphrase *string, encrypted []byte) ([]byte, []byte, error) {
	var (
		naclPub   []byte
		naclPriv  []byte
		tmpSecret [32]byte
		tmpNonce  [24]byte
		err       error
		password  []byte
	)

	if encrypted == nil || len(encrypted) == 0 {
		return naclPub, naclPriv, ErrNoEncrypted
	}

	encoded := encrypted
	if passphrase != nil {
		if len(encrypted) < 24 {
			return naclPub, naclPriv, ErrEncryptedLength
		}

		salt := encrypted[:saltLength]

		N := u8util.ToBN(encrypted[32+0:32+4], true).Int64()
		p := u8util.ToBN(encrypted[32+4:32+8], true).Int64()
		r := u8util.ToBN(encrypted[32+8:32+12], true).Int64()
		if N != defaultN || p != defaultP || r != defaultR {
			return nil, nil, ErrInvalidParams
		}
		password, err = scrypt.Key([]byte(*passphrase), salt, int(N), int(r), int(p), 64)
		if err != nil {
			return nil, nil, err
		}
		encrypted = encrypted[scryptLength:]
		secret := u8util.FixLength(password, 256, true)
		if len(secret) != 32 {
			return naclPub, naclPriv, ErrSecretLength
		}
		copy(tmpSecret[:], secret)
		copy(tmpNonce[:], encrypted[0:nonceLength])
		encoded, err = crypto.NaclDecrypt(encrypted[nonceLength:], tmpNonce, tmpSecret)
		if err != nil {
			return nil, nil, err
		}
	}

	if encoded == nil || len(encoded) == 0 {
		return naclPub, naclPriv, ErrEncoded
	}
	header := encoded[:seedOffset]
	if string(header) != string(pkcs8Header) {
		return naclPub, naclPriv, ErrPkcs8Header
	}
	// note: check encoded lengths?
	secretKey := encoded[seedOffset : seedOffset+secLength]
	divider := encoded[divOffset : divOffset+len(pkcs8Divider)]
	if !bytes.Equal(divider, pkcs8Divider) {
		divOffset = seedOffset + seedLength
		secretKey = encoded[seedOffset:divOffset]
		divider = encoded[divOffset : divOffset+len(pkcs8Divider)]
		if !bytes.Equal(divider, pkcs8Divider) {
			return naclPub, naclPriv, ErrPkcs8Divider
		}
	}
	pubOffset := divOffset + len(pkcs8Divider)
	publicKey := encoded[pubOffset : pubOffset+pubLength]
	return publicKey, secretKey, nil
}

func signingContext(msg []byte) *merlin.Transcript {
	tml := merlin.NewTranscript("SigningContext")
	tml.AppendMessage([]byte(""), []byte("substrate"))
	tml.AppendMessage([]byte("sign-bytes"), msg)
	return tml
}

func Sign(sk []byte, t *merlin.Transcript, publicKey []byte) (*schnorrkel.Signature, error) {
	var (
		pubKey [32]byte
		rByte  [64]byte
		sck    [32]byte
	)

	copy(pubKey[:], publicKey[:])
	pubByte := schnorrkel.NewPublicKey(pubKey).Compress()

	t.AppendMessage([]byte("proto-name"), []byte("Schnorr-sig"))
	t.AppendMessage([]byte("sign:pk"), pubByte[:])

	_, err := t.BuildRNG().ReKeyWithWitnessBytes([]byte("signing"), sk[32:]).Read(rByte[:])
	if err != nil {
		return nil, err
	}

	r := ristretto255.NewScalar().FromUniformBytes(rByte[:])
	R := ristretto255.NewElement().ScalarBaseMult(r)
	t.AppendMessage([]byte("sign:R"), R.Encode([]byte{}))

	// form k
	k := ristretto255.NewScalar().FromUniformBytes(t.ExtractBytes([]byte("sign:c"), 64))
	//k.

	// form scalar from secret key x
	key := divideScalarByCofactor(sk[:32])
	copy(sck[:], key[:])

	x, err := schnorrkel.ScalarFromBytes(sck)
	if err != nil {
		return nil, err
	}

	// s = kx + r
	s := x.Multiply(x, k).Add(x, r)

	return &schnorrkel.Signature{R: R, S: s}, nil
}

func divideScalarByCofactor(s []byte) []byte {
	l := len(s) - 1
	low := byte(0)
	for i := range s {
		r := s[l-i] & 0x07 // remainder
		s[l-i] >>= 3
		s[l-i] += low
		low = r << 5
	}

	return s
}
