package eth

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Util struct {
}

func NewUtil() *Util {
	return &Util{}
}

// MARK - Implement the protocol wallet.Util

// The ethereum public key is same as address in coming
func (u *Util) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return publicKey, nil
}

// The ethereum public key is same as address in coming
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	return address, nil
}

// It will check based on eip55 rules
func (u *Util) IsValidAddress(address string) bool {
	return IsValidAddress(address)
}

// MARK - like wallet.Util

// The ethereum public key is same as address in coming
func EncodePublicKeyToAddress(publicKey string) (string, error) {
	return publicKey, nil
}

// The ethereum public key is same as address in coming
func DecodeAddressToPublicKey(address string) (string, error) {
	return address, nil
}

// It will check based on eip55 rules
func IsValidAddress(address string) bool {
	eip55Address, err := TransformEIP55Address(address)
	if err != nil {
		return false
	}

	return strings.HasSuffix(eip55Address, address)
}

func TransformEIP55Address(address string) (string, error) {
	address = strings.TrimPrefix(address, "0x")
	if !common.IsHexAddress(address) {
		return "", errors.New("Invalid hex address")
	}

	addressBytes := []byte(strings.ToLower(address))
	checksumBytes := crypto.Keccak256(addressBytes)

	for i, c := range addressBytes {
		if c >= '0' && c <= '9' {
			continue
		} else {
			checksum := checksumBytes[i/2]
			bitcode := byte(0x80) >> ((i % 2) * 4)
			if checksum&bitcode > 0 { // to Upper
				addressBytes[i] -= 32
			}
		}
	}

	return "0x" + string(addressBytes), nil
}
