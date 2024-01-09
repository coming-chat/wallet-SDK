package starknet

import (
	"errors"
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

// EncodePublicKeyToAddress
// - return: encode cairo version 1.0 address
func EncodePublicKeyToAddress(publicKey string) (string, error) {
	return encodePublicKeyToAddressArgentX(publicKey, false)
}

func EncodePublicKeyToAddressCairo0(publicKey string) (string, error) {
	return encodePublicKeyToAddressArgentX(publicKey, true)
}

func encodePublicKeyToAddressArgentX(publicKey string, isCairo0 bool) (string, error) {
	pubFelt, err := utils.HexToFelt(publicKey)
	if err != nil {
		return "", err
	}
	p := deployParamForArgentXWithVersion(*pubFelt, isCairo0)
	addr, err := p.ComputeContractAddress()
	if err != nil {
		return "", err
	}
	return fullString(*addr), nil
}

// Warning: starknet cannot support decode address to public key
func DecodeAddressToPublicKey(address string) (string, error) {
	return "", base.ErrUnsupportedFunction
}

func IsValidAddress(address string) bool {
	reg := regexp.MustCompile(`^(0x|0X)?[0-9a-fA-F]{40,64}$`)
	return reg.MatchString(address)
}

// CheckCairoVersion
// - return address version, 0 for cairo0, 1 for cairo1.0
func CheckCairoVersion(address, pubkey string) (*base.OptionalInt, error) {
	addrFelt, err := utils.HexToFelt(address)
	if err != nil {
		return nil, base.ErrInvalidAddress
	}
	pubFelt, err := utils.HexToFelt(pubkey)
	if err != nil {
		return nil, base.ErrInvalidPublicKey
	}
	version, err := CheckCairoVersionFelt(addrFelt, pubFelt)
	if err != nil {
		return nil, err
	}
	return base.NewOptionalInt(version), nil
}

func CheckCairoVersionFelt(address, pubkey *felt.Felt) (int, error) {
	p1 := deployParamForArgentX(*pubkey)
	addr, err := p1.ComputeContractAddress()
	if err == nil {
		if addr.Cmp(address) == 0 {
			return 1, nil
		}
	}
	p0 := deployParamForArgentXCairo0(*pubkey)
	addr, err = p0.ComputeContractAddress()
	if err == nil {
		if addr.Cmp(address) == 0 {
			return 0, nil
		}
	}
	return -1, errors.New("the address and public key mismatch")
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

// var defaultDeployParam = deployParamForBraavos
var defaultDeployParam = deployParamForArgentX

func deployParamForArgentXWithVersion(pub felt.Felt, isCairo0 bool) deployParam {
	if isCairo0 {
		return deployParamForArgentXCairo0(pub)
	} else {
		return deployParamForArgentX(pub)
	}
}

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

func deployParamForArgentXCairo0(pub felt.Felt) deployParam {
	return deployParam{
		Caller:    &felt.Zero,
		Pubkey:    pub,
		ClassHash: mustFelt("25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918"),
		CallData: []*felt.Felt{
			mustFelt("33434ad846cdd5f23eb73ff09fe6fddd568284a0fb7d1be20ee482f044dabe2"),
			mustFelt("79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463"),
			mustFelt("2"),
			&pub,
			&felt.Zero,
		},
	}
}
