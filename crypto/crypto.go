package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math"

	"github.com/pierrec/xxHash/xxHash64"
	"golang.org/x/crypto/blake2b"

	//"github.com/agl/ed25519"
	//naclauth "golang.org/x/crypto/nacl/auth"
	naclsecret "golang.org/x/crypto/nacl/secretbox"
	naclsign "golang.org/x/crypto/nacl/sign"
)

// Hash ...
type Hash [sha256.Size]uint8

// Blake2b256Hash ...
type Blake2b256Hash [blake2b.Size256]uint8

// Blake2b512Hash ...
type Blake2b512Hash [blake2b.Size]uint8

// NewSHA256 ...
func NewSHA256(data []byte) *Hash {
	var hash Hash
	hash = sha256.Sum256(data)
	return &hash
}

// NewBlake2b256 ...
func NewBlake2b256(data []byte) *Blake2b256Hash {
	var hash Blake2b256Hash
	hash = blake2b.Sum256(data)
	return &hash
}

// NewBlake2b512 ...
func NewBlake2b512(data []byte) *Blake2b512Hash {
	var hash Blake2b512Hash
	hash = blake2b.Sum512(data)
	return &hash
}

// NewXXHash ...
func NewXXHash(data []byte, bitLength int64) []byte {
	return newXXHash(data, uint(bitLength))
}

// NewXXHash64 ...
func NewXXHash64(data []byte) [8]byte {
	var hash [8]byte
	copy(hash[:], newXXHash(data, 64))
	return hash
}

// NewXXHash128 ...
func NewXXHash128(data []byte) [16]byte {
	var hash [16]byte
	copy(hash[:], newXXHash(data, 128))
	return hash
}

// NewXXHash256 ...
func NewXXHash256(data []byte) [32]byte {
	var hash [32]byte
	copy(hash[:], newXXHash(data, 256))
	return hash
}

func newXXHash(data []byte, bitLength uint) []byte {
	byteLength := int64(math.Ceil(float64(bitLength) / float64(8)))
	iterations := int64(math.Ceil(float64(bitLength) / float64(64)))
	var hash = make([]byte, byteLength)

	for seed := int64(0); seed < iterations; seed++ {
		digest := xxHash64.New(uint64(seed))
		digest.Write(data)
		copy(hash[seed*8:], digest.Sum(nil))
	}

	return hash
}

// NewNaclKeyPair ...
func NewNaclKeyPair() ([32]byte, [64]byte, error) {
	var (
		pub  *[32]byte
		priv *[64]byte
		err  error
	)

	pub, priv, err = naclsign.GenerateKey(rand.Reader)
	if err != nil {
		return [32]byte{}, [64]byte{}, err
	}
	if pub == nil || priv == nil {
		return [32]byte{}, [64]byte{}, errors.New("nil keys")
	}

	return *pub, *priv, nil
}

// NewNaclKeyPairFromSeed ...
// note: return pointers???
func NewNaclKeyPairFromSeed(seed []byte) ([32]byte, [64]byte, error) {
	var (
		pub  *[32]byte
		priv *[64]byte
		err  error
	)

	reader := bytes.NewBuffer(seed)
	pub, priv, err = naclsign.GenerateKey(reader)
	if err != nil {
		return [32]byte{}, [64]byte{}, err
	}
	if pub == nil || priv == nil {
		return [32]byte{}, [64]byte{}, errors.New("nil keys")
	}

	return *pub, *priv, nil
}

//// NaclVerify returns true if signature is a valid signature of digest by public key. Using the 'agl' library because native nacl library doesn't support detached signatures.
//func NaclVerify(digest []byte, signature []byte, publicKey [naclauth.KeySize]byte) bool {
//	var sig64 [64]byte
//	copy(sig64[:], signature)
//	return ed25519.Verify(&publicKey, digest, &sig64)
//}
//
//// NaclSign ...
//// Using the 'agl' library because native nacl library doesn't support detached signatures.
//func NaclSign(secret [64]byte, message []byte) ([]byte, error) {
//	if message == nil || len(message) == 0 {
//		return nil, errors.New("cannot sign nil message")
//	}
//
//	sig := ed25519.Sign(&secret, message)
//	if sig == nil {
//		return nil, errors.New("could not sign message")
//	}
//
//	return sig[:], nil
//}

// NaclEncrypt ...
// note: use pointers???
func NaclEncrypt(message []byte, nonce [24]byte, secret [32]byte) ([]byte, error) {
	if message == nil || len(message) == 0 {
		return nil, errors.New("cannot encrypt nil message")
	}

	var out []byte
	out = naclsecret.Seal(out, message, &nonce, &secret)
	return out, nil
}

// NaclDecrypt ...
// note: use pointers???
func NaclDecrypt(box []byte, nonce [24]byte, secret [32]byte) ([]byte, error) {
	if box == nil || len(box) == 0 {
		return nil, errors.New("cannot decrypt a nil message")
	}

	var (
		out []byte
		ok  bool
	)

	out, ok = naclsecret.Open(out, box, &nonce, &secret)
	if !ok {
		return nil, errors.New("could not decrypt message")
	}

	return out, nil
}

// NewBlake2b256Sig ...
func NewBlake2b256Sig(key, data []byte) ([]byte, error) {
	hash, err := blake2b.New256(key)
	if err != nil {
		return nil, err
	}

	hash.Write(data)
	return hash.Sum(nil), nil
}

// NewBlake2b512Sig ...
func NewBlake2b512Sig(key, data []byte) ([]byte, error) {
	hash, err := blake2b.New512(key)
	if err != nil {
		return nil, err
	}

	hash.Write(data)
	return hash.Sum(nil), nil
}
