package starknet

import (
	"regexp"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/dontpanicdao/caigo/types"
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

// Warning: starknet cannot support decode address to public key
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	return "", base.ErrUnsupportedFunction
}

// Check if address is 40 hexadecimal characters
func (u *Util) IsValidAddress(address string) bool {
	return IsValidAddress(address)
}

// MARK - like wallet.Util

func EncodePublicKeyToAddress(publicKey string) (string, error) {
	return encodePublicKeyToAddressArgentX(publicKey)
}

func encodePublicKeyToAddressArgentX(publicKey string) (string, error) {
	txn, err := newDeployAccountTransaction(publicKey, 0)
	if err != nil {
		return "", err
	}
	return types.BigToHex(txn.ContractAddress), nil
}

// Warning: starknet cannot support decode address to public key
func DecodeAddressToPublicKey(address string) (string, error) {
	return "", base.ErrUnsupportedFunction
}

func IsValidAddress(address string) bool {
	reg := regexp.MustCompile(`^(0x|0X)?[0-9a-fA-F]{1,64}$`)
	return reg.MatchString(address)
}
