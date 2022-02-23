package wallet

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/coming-chat/merlin"
	"github.com/coming-chat/wallet-SDK/crypto"
	"github.com/coming-chat/wallet-SDK/u8util"
	"github.com/gtank/ristretto255"
	"golang.org/x/crypto/scrypt"
	"log"
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

type keystore struct {
	Encoded  string    `json:"encoded"`
	Encoding *encoding `json:"encoding"`
	Address  string    `json:"address"`
}

type encoding struct {
	Content []string `json:"content"`
	Type    []string `json:"type"`
	Version string   `json:"version"`
}

type keyring struct {
	privateKey [64]byte
	PublicKey  [32]byte
}

func (k *keystore) checkPassword(password string) error {
	_, err := decodeKeystore(k, password)
	return err
}

func (k *keystore) Sign(msg []byte, password string) ([]byte, error) {
	kr, err := decodeKeystore(k, password)
	if err != nil {
		return nil, err
	}
	signature, err := kr.sign(signingContext(msg))
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func decodeKeystore(ks *keystore, password string) (*keyring, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			return
		}
	}()
	var (
		privateKey [64]byte
		publicKey  [32]byte
	)

	if ks.Encoding == nil || ks.Encoding.Version != "3" || len(ks.Encoding.Content) < 2 || ks.Encoding.Content[0] != "pkcs8" || ks.Encoding.Content[1] != "sr25519" {
		return nil, ErrNonPkcs8
	}

	encrypted, err := base64.RawStdEncoding.DecodeString(ks.Encoded)
	if err != nil {
		return nil, err
	}
	pubKey, secretKey, err := decodePolkaKeystoreEncoded(&password, encrypted, ks.Encoding)
	if err != nil {
		return nil, err
	}

	copy(publicKey[:], pubKey[:])
	copy(privateKey[:], secretKey[:])
	addrPubKey, err := AddressToPublicKey(ks.Address)
	if err != nil {
		return nil, err
	}
	if addrPubKey != "0x"+hex.EncodeToString(publicKey[:]) {
		return nil, ErrKeystore
	}

	return &keyring{
		privateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

func decodePolkaKeystoreEncoded(passphrase *string, encrypted []byte, encodeType *encoding) ([]byte, []byte, error) {
	var (
		tmpSecret [32]byte
		tmpNonce  [24]byte
		err       error
		password  []byte
	)

	if len(encodeType.Type) < 2 || encodeType.Type[1] != "xsalsa20-poly1305" {
		return nil, nil, ErrNoEncryptedData
	}

	if passphrase == nil {
		return nil, nil, ErrNilPassword
	}

	if encrypted == nil || len(encrypted) == 0 {
		return nil, nil, ErrNoEncrypted
	}

	encoded := encrypted

	if len(encrypted) < 24 {
		return nil, nil, ErrEncryptedLength
	}

	if encodeType.Type[0] == "scrypt" {
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

	} else {
		password = []uint8(*passphrase)
	}

	secret := u8util.FixLength(password, 256, true)
	if len(secret) != 32 {
		return nil, nil, ErrSecretLength
	}

	copy(tmpSecret[:], secret)
	copy(tmpNonce[:], encrypted[0:nonceLength])

	encoded, err = crypto.NaclDecrypt(encrypted[nonceLength:], tmpNonce, tmpSecret)
	if err != nil {
		return nil, nil, err
	}

	if encoded == nil || len(encoded) == 0 {
		return nil, nil, ErrEncoded
	}
	header := encoded[:seedOffset]
	if string(header) != string(pkcs8Header) {
		return nil, nil, ErrPkcs8Header
	}
	// note: check encoded lengths?
	secretKey := encoded[seedOffset : seedOffset+secLength]
	divider := encoded[divOffset : divOffset+len(pkcs8Divider)]
	if !bytes.Equal(divider, pkcs8Divider) {
		divOffset = seedOffset + seedLength
		secretKey = encoded[seedOffset:divOffset]
		divider = encoded[divOffset : divOffset+len(pkcs8Divider)]
		if !bytes.Equal(divider, pkcs8Divider) {
			return nil, nil, ErrPkcs8Divider
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

func (kg *keyring) sign(t *merlin.Transcript) ([]byte, error) {
	var (
		pubKey [32]byte
		rByte  [64]byte
		sck    [32]byte
	)

	copy(pubKey[:], kg.PublicKey[:])
	pubByte := schnorrkel.NewPublicKey(pubKey).Compress()

	t.AppendMessage([]byte("proto-name"), []byte("Schnorr-sig"))
	t.AppendMessage([]byte("sign:pk"), pubByte[:])

	_, err := t.BuildRNG().ReKeyWithWitnessBytes([]byte("signing"), kg.privateKey[32:]).Read(rByte[:])
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
	key := divideScalarByCofactor(kg.privateKey[:32])
	copy(sck[:], key[:])

	x, err := schnorrkel.ScalarFromBytes(sck)
	if err != nil {
		return nil, err
	}

	// s = kx + r
	s := x.Multiply(x, k).Add(x, r)

	signature := &schnorrkel.Signature{R: R, S: s}
	signatureByte := signature.Encode()

	return signatureByte[:], nil
}

//func (kg *keyring) signEd25519(t *merlin.Transcript, signature, message []byte) {
//	if l := len(kg.privateKey); l != PrivateKeySize {
//		panic("ed25519: bad private key length: " + strconv.Itoa(l))
//	}
//	seed, publicKey := privateKey[:SeedSize], privateKey[SeedSize:]
//
//	h := sha512.Sum512(seed)
//	s := edwards25519.NewScalar().SetBytesWithClamping(h[:32])
//	prefix := h[32:]
//
//	mh := sha512.New()
//	mh.Write(prefix)
//	mh.Write(message)
//	messageDigest := make([]byte, 0, sha512.Size)
//	messageDigest = mh.Sum(messageDigest)
//	r := edwards25519.NewScalar().SetUniformBytes(messageDigest)
//
//	R := (&edwards25519.Point{}).ScalarBaseMult(r)
//
//	kh := sha512.New()
//	kh.Write(R.Bytes())
//	kh.Write(publicKey)
//	kh.Write(message)
//	hramDigest := make([]byte, 0, sha512.Size)
//	hramDigest = kh.Sum(hramDigest)
//	k := edwards25519.NewScalar().SetUniformBytes(hramDigest)
//
//	S := edwards25519.NewScalar().MultiplyAdd(k, s, r)
//
//	copy(signature[:32], R.Bytes())
//	copy(signature[32:], S.Bytes())
//}

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
