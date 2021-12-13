package wallet

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/scrypt"
	"log"
	"wallet-SDK/crypto"
	"wallet-SDK/u8util"
)

var (
	// DEFAULT_PKCS8_DIVIDER ...
	DEFAULT_PKCS8_DIVIDER = []byte{161, 35, 3, 33, 0}
	// DEFAULT_PKCS8_HEADER ...
	DEFAULT_PKCS8_HEADER = []byte{48, 83, 2, 1, 1, 48, 5, 6, 3, 43, 101, 112, 4, 34, 4, 32}
	// DEFAULT_KEY_LENGTH ...
	DEFAULT_KEY_LENGTH = 32
	// DEFAULT_SEED_OFFSET ...
	DEFAULT_SEED_OFFSET = len(DEFAULT_PKCS8_HEADER)
	// DEFAULT_DIV_OFFSET ...
	DEFAULT_DIV_OFFSET = DEFAULT_SEED_OFFSET + SEC_LENGTH
	// DEFAULT_PUBLIC_OFFSET ...
	DEFAULT_PUBLIC_OFFSET = DEFAULT_SEED_OFFSET + DEFAULT_KEY_LENGTH + len(DEFAULT_PKCS8_DIVIDER)

	PUB_LENGTH  = 32
	SALT_LENGTH = 32
	SEC_LENGTH  = 64
	SEED_LENGTH = 32

	SCRYPT_LENGTH = 32 + (3 * 4)
	NONCE_LENGTH  = 24

	defaultN int64 = 1 << 15
	defaultP int64 = 1
	defaultR int64 = 8
)

// Decode ...
func Decode(passphrase *string, encrypted []byte) ([]byte, []byte, error) {
	var (
		naclPub  []byte
		naclPriv []byte
	)

	if encrypted == nil || len(encrypted) == 0 {
		return naclPub, naclPriv, errors.New("no encrypted data to decode")
	}

	encoded := encrypted
	if passphrase != nil {
		//if len(encrypted) < 24 {
		//	return naclPub, naclPriv, errors.New("encrypted length is less than 24")
		//}

		salt := encrypted[:32]

		N := u8util.ToBN(encrypted[32+0:32+4], true).Int64()
		p := u8util.ToBN(encrypted[32+4:32+8], true).Int64()
		r := u8util.ToBN(encrypted[32+8:32+12], true).Int64()
		if N != defaultN || p != defaultP || r != defaultR {
			return nil, nil, errors.New("invalid injected scrypt params found")
		}
		var (
			tmpSecret [32]byte
			tmpNonce  [24]byte
			err       error
			password  []byte
		)
		password, err = scrypt.Key([]byte(*passphrase), salt, int(N), int(r), int(p), 64)
		if err != nil {
			return nil, nil, err
		}
		encrypted = encrypted[SCRYPT_LENGTH:]
		secret := u8util.FixLength(password, 256, true)
		if len(secret) != 32 {
			log.Println(secret, len(secret))
			return naclPub, naclPriv, errors.New("secret length is not 32")
		}
		copy(tmpSecret[:], secret)
		copy(tmpNonce[:], encrypted[0:NONCE_LENGTH])
		encoded, err = crypto.NaclDecrypt(encrypted[NONCE_LENGTH:], tmpNonce, tmpSecret)
		if err != nil {
			return nil, nil, err
		}
	}

	if encoded == nil || len(encoded) == 0 {
		return naclPub, naclPriv, errors.New("unable to decode")
	}
	header := encoded[:DEFAULT_SEED_OFFSET]
	if string(header) != string(DEFAULT_PKCS8_HEADER) {
		return naclPub, naclPriv, errors.New("invalid Pkcs8 header found in body")
	}
	// note: check encoded lengths?
	secretKey := encoded[DEFAULT_SEED_OFFSET : DEFAULT_SEED_OFFSET+SEC_LENGTH]
	divider := encoded[DEFAULT_DIV_OFFSET : DEFAULT_DIV_OFFSET+len(DEFAULT_PKCS8_DIVIDER)]
	if !bytes.Equal(divider, DEFAULT_PKCS8_DIVIDER) {
		divOffset := DEFAULT_SEED_OFFSET + SEED_LENGTH
		secretKey = encoded[DEFAULT_SEED_OFFSET:divOffset]
		divider = encoded[divOffset : divOffset+len(DEFAULT_PKCS8_DIVIDER)]
		if !bytes.Equal(divider, DEFAULT_PKCS8_DIVIDER) {
			return naclPub, naclPriv, errors.New("invalid Pkcs8 divider found in body")
		}
	}
	pubOffset := DEFAULT_DIV_OFFSET + len(DEFAULT_PKCS8_DIVIDER)
	publicKey := encoded[pubOffset : pubOffset+PUB_LENGTH]
	return publicKey, secretKey, nil
}
