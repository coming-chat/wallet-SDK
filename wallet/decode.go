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
	defaultPkcs8Divider = []byte{161, 35, 3, 33, 0}
	defaultPkcs8Header  = []byte{48, 83, 2, 1, 1, 48, 5, 6, 3, 43, 101, 112, 4, 34, 4, 32}
	defaultSeedOffset   = len(defaultPkcs8Header)
	defaultDivOffset    = defaultSeedOffset + secLength
	defaultPublicOffset = defaultSeedOffset + defaultKeyLength + len(defaultPkcs8Divider)
)

const (
	defaultKeyLength = 32
	pubLength        = 32
	saltLength       = 32
	secLength        = 64
	seedLength       = 32

	scryptLength = 32 + (3 * 4)
	nonceLength  = 24

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
		encrypted = encrypted[scryptLength:]
		secret := u8util.FixLength(password, 256, true)
		if len(secret) != 32 {
			log.Println(secret, len(secret))
			return naclPub, naclPriv, errors.New("secret length is not 32")
		}
		copy(tmpSecret[:], secret)
		copy(tmpNonce[:], encrypted[0:nonceLength])
		encoded, err = crypto.NaclDecrypt(encrypted[nonceLength:], tmpNonce, tmpSecret)
		if err != nil {
			return nil, nil, err
		}
	}

	if encoded == nil || len(encoded) == 0 {
		return naclPub, naclPriv, errors.New("unable to decode")
	}
	header := encoded[:defaultSeedOffset]
	if string(header) != string(defaultPkcs8Header) {
		return naclPub, naclPriv, errors.New("invalid Pkcs8 header found in body")
	}
	// note: check encoded lengths?
	secretKey := encoded[defaultSeedOffset : defaultSeedOffset+seedLength]
	divider := encoded[defaultDivOffset : defaultDivOffset+len(defaultPkcs8Divider)]
	if !bytes.Equal(divider, defaultPkcs8Divider) {
		divOffset := defaultSeedOffset + seedLength
		secretKey = encoded[defaultSeedOffset:divOffset]
		divider = encoded[divOffset : divOffset+len(defaultPkcs8Divider)]
		if !bytes.Equal(divider, defaultPkcs8Divider) {
			return naclPub, naclPriv, errors.New("invalid Pkcs8 divider found in body")
		}
	}
	pubOffset := defaultDivOffset + len(defaultPkcs8Divider)
	publicKey := encoded[pubOffset : pubOffset+pubLength]
	return publicKey, secretKey, nil
}
