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
	txn, err := newDeployAccountTransactionForArgentX(publicKey, 0)
	if err != nil {
		return "", err
	}
	return types.BigToHex(txn.ContractAddress), nil
}

// func encodePublicKeyToAddressBraavos(publicKey string) (string, error) {
// 	pub, err := new(felt.Felt).SetString(publicKey)
// 	if err != nil {
// 		return "", base.ErrInvalidPublicKey
// 	}
// 	callerAddress, _ := new(felt.Felt).SetString("0x0000000000000000000000000000000000000000")
// 	classHash, _ := new(felt.Felt).SetString("0x03131fa018d520a037686ce3efddeab8f28895662f019ca3ca18a626650f7d1e")
// 	data1, _ := new(felt.Felt).SetString("0x5aa23d5bb71ddaa783da7ea79d405315bafa7cf0387a74f4593578c3e9e6570")
// 	data2, _ := new(felt.Felt).SetString("0x2dd76e7ad84dbed81c314ffe5e7a7cacfb8f4836f01af4e913f275f89a3de1a")
// 	data3, _ := new(felt.Felt).SetString("0x1")
// 	data4 := pub
// 	callData := []*felt.Felt{
// 		data1,
// 		data2,
// 		data3,
// 		data4,
// 	}
// 	address := core.ContractAddress(callerAddress, classHash, pub, callData)
// 	return address.String(), nil
// }

// Warning: starknet cannot support decode address to public key
func DecodeAddressToPublicKey(address string) (string, error) {
	return "", base.ErrUnsupportedFunction
}

func IsValidAddress(address string) bool {
	reg := regexp.MustCompile(`^(0x|0X)?[0-9a-fA-F]{1,64}$`)
	return reg.MatchString(address)
}
