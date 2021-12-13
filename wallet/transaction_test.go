package wallet

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

const (
	// PublicKeySize is the size, in bytes, of public keys as used in this package.
	PublicKeySize = 32
	// PrivateKeySize is the size, in bytes, of private keys as used in this package.
	PrivateKeySize = 64
	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize = 64
	// SeedSize is the size, in bytes, of private key seeds. These are the private key representations used by RFC 8032.
	SeedSize = 32
)

func TestPolDecode(t *testing.T) {
	decodeByte, err := base64.RawStdEncoding.DecodeString("q5FkOBeJtiNT/KpwBAqGh1zYN767e87A14tKl+n+M4oAgAAAAQAAAAgAAADbeVjFCVVlKoDxJhGY0nrSJgx4C2Insm2UXPZUpA6eQyj2ETpa4U5pA2FwlJVcrCMwup1+Sghk3a07osZrspnj722OyhtNJSpthiNDsicHNeBWO17WTxAIA8n6URQa87hTDDk6nQmMjSYwgl5eoKaGD77mRxe5ROdb1XdX143Vu2n47iSTZ+LInIrJUn/Bk4ajKmQZ5kYXbS5HE88S")
	var s = "111"
	publicKey, secretKeys, err := Decode(&s, decodeByte)
	if err != nil {
		return
	}
	sr, _ := hex.DecodeString(testSecretSeed)
	t.Log(hex.EncodeToString(publicKey), hex.EncodeToString(secretKeys[32:]), hex.EncodeToString(secretKeys[:32]), sr)
	w, err := NewWallet(Minix, "0x"+hex.EncodeToString(secretKeys), 44)
	if err != nil {
		t.Error(err)
	}
	t.Log(w)
	//if err != nil {
	//	return
	//}
	//salt := decodeByte[:32]
	////N := decodeByte[32:36]
	////Nhx := hex.EncodeToString(N)
	////p := decodeByte[36:40]
	////r := decodeByte[40:44]
	//N := 1 << 15
	//p := 1
	//r := 8
	//pwBytes, err := scrypt.Key([]byte("111"), salt, int(math.Log2(float64(N))), r, p, 32)
	//if err != nil {
	//	return
	//}
	//privateKey, err := pkcs8.
	//(testSecretSeed, "111", tt.opts)
	//if err != nil {
	//	return
	//}
	//key, n, err := pkcs8.ParsePrivateKey(decodeByte, []byte("111"))
	//if err != nil {
	//	return
	//}
	////if err != nil {
	////	return
	////}
	//t.Log(key, n)

}

//func Sign(data []byte, privateKeyURI string) ([]byte, error) {
//	// if data is longer than 256 bytes, hash it first
//	if len(data) > 256 {
//		h := blake2b.Sum256(data)
//		data = h[:]
//	}
//
//	scheme := sr25519.Scheme{}
//	kyr, err := subkey.DeriveKeyPair(scheme, privateKeyURI)
//	if err != nil {
//		return nil, err
//	}
//
//	signature, err := kyr.Sign(data)
//	if err != nil {
//		return nil, err
//	}
//
//	return signature, nil
//}
//
//func sign(signature, privateKey, message []byte) {
//	if l := len(privateKey); l != PrivateKeySize {
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
