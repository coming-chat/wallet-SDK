package sui

import (
	"encoding/hex"
	"errors"
	"regexp"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"golang.org/x/crypto/sha3"
)

type Util struct {
}

func NewUtil() (*Util, error) {
	return &Util{}, nil
}

// MARK - Implement the protocol Util

// @param publicKey can start with 0x or not.
func (u *Util) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return EncodePublicKeyToAddress(publicKey)
}

// Warning: Sui cannot support decode address to public key
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	return DecodeAddressToPublicKey(address)
}

func (u *Util) IsValidAddress(address string) bool {
	return IsValidAddress(address)
}

// MARK - like Util

// @param publicKey can start with 0x or not.
func EncodePublicKeyToAddress(publicKey string) (string, error) {
	publicBytes, err := types.HexDecodeString(publicKey)
	if err != nil {
		return "", err
	}

	tmp := []byte{0x00}
	tmp = append(tmp, publicBytes...)
	addrBytes := sha3.Sum256(tmp)
	return "0x" + hex.EncodeToString(addrBytes[:])[:40], nil
}

func DecodeAddressToPublicKey(address string) (string, error) {
	return "", errors.New("Sui cannot support decode address to public key")
}

// @param chainnet chain name
func IsValidAddress(address string) bool {
	reg := regexp.MustCompile(`^(0x|0X)?[0-9a-fA-F]{1,40}$`)
	return reg.MatchString(address)
}
