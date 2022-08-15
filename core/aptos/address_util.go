package aptos

import (
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

// Warning: Aptos cannot support decode address to public key
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
	publicBytes = append(publicBytes, 0x00)
	addressBytes := sha3.Sum256(publicBytes)
	return types.HexEncodeToString(addressBytes[:]), nil
}

func DecodeAddressToPublicKey(address string) (string, error) {
	return "", errors.New("Aptos cannot support decode address to public key")
}

// @param chainnet chain name
func IsValidAddress(address string) bool {
	reg := regexp.MustCompile(`^(0x|0X)?[0-9a-fA-F]{1,64}$`)
	return reg.MatchString(address)
}
