package starknet

import (
	"regexp"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/xiang-xx/starknet.go/account"
	"github.com/xiang-xx/starknet.go/utils"
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
	addr, err := defaultEncodeContractAddress(publicKey)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

// Warning: starknet cannot support decode address to public key
func DecodeAddressToPublicKey(address string) (string, error) {
	return "", base.ErrUnsupportedFunction
}

func IsValidAddress(address string) bool {
	reg := regexp.MustCompile(`^(0x|0X)?[0-9a-fA-F]{1,64}$`)
	return reg.MatchString(address)
}

// MARK - Contract Address Generate

type deployParam struct {
	Caller    *felt.Felt
	Pubkey    felt.Felt
	ClassHash *felt.Felt
	CallData  []*felt.Felt
}

func (p *deployParam) ComputeContractAddress() (*felt.Felt, error) {
	acc := account.Account{}
	return acc.PrecomputeAddress(p.Caller, &p.Pubkey, p.ClassHash, p.CallData)
}

func defaultEncodeContractAddress(pub string) (*felt.Felt, error) {
	pubFelt, err := utils.HexToFelt(pub)
	if err != nil {
		return nil, err
	}
	p := defaultDeployParam(*pubFelt)
	acc := account.Account{}
	return acc.PrecomputeAddress(p.Caller, &p.Pubkey, p.ClassHash, p.CallData)
}

// var defaultDeployParam = deployParamForBraavos
var defaultDeployParam = deployParamForArgentX

func deployParamForArgentX(pub felt.Felt) deployParam {
	return deployParam{
		Caller:    &felt.Zero,
		Pubkey:    pub,
		ClassHash: mustFelt("01a736d6ed154502257f02b1ccdf4d9d1089f80811cd6acad48e6b6a9d1f2003"),
		CallData: []*felt.Felt{
			&pub,
			&felt.Zero,
		},
	}
}

func deployParamForBraavos(pub felt.Felt) deployParam {
	return deployParam{
		Caller:    &felt.Zero,
		Pubkey:    pub,
		ClassHash: mustFelt("03131fa018d520a037686ce3efddeab8f28895662f019ca3ca18a626650f7d1e"),
		CallData: []*felt.Felt{
			mustFelt("5aa23d5bb71ddaa783da7ea79d405315bafa7cf0387a74f4593578c3e9e6570"),
			mustFelt("2dd76e7ad84dbed81c314ffe5e7a7cacfb8f4836f01af4e913f275f89a3de1a"),
			mustFelt("1"),
			&pub,
		},
	}
}
