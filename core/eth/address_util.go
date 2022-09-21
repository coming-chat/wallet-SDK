package eth

import (
	"errors"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Util struct {
}

func NewUtil() *Util {
	return &Util{}
}

// MARK - Implement the protocol wallet.Util

func (u *Util) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return EncodePublicKeyToAddress(publicKey)
}

// Warning: eth cannot support decode address to public key
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	return "", errors.New("eth cannot support decode address to public key")
}

// Check if address is 40 hexadecimal characters
func (u *Util) IsValidAddress(address string) bool {
	return IsValidAddress(address)
}

// MARK - like wallet.Util

func EncodePublicKeyToAddress(publicKey string) (string, error) {
	bytes, err := types.HexDecodeString(publicKey)
	if err != nil {
		return "", err
	}
	publicKeyECDSA, err := crypto.UnmarshalPubkey(bytes)
	if err != nil {
		return "", err
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address.String(), nil
}

// Warning: eth cannot support decode address to public key
func DecodeAddressToPublicKey(address string) (string, error) {
	return "", errors.New("eth cannot support decode address to public key")
}

// Check if address is 40 hexadecimal characters
func IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

// It will check based on eip55 rules
func IsValidEIP55Address(address string) bool {
	if !IsValidAddress(address) {
		return false
	}
	eip55Address := TransformEIP55Address(address)
	return strings.HasSuffix(eip55Address, address)
}

func TransformEIP55Address(address string) string {
	address = strings.TrimPrefix(address, "0x")
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

	return "0x" + string(addressBytes)
}
