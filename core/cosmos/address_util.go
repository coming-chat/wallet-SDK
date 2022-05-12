package cosmos

import (
	"errors"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

type Util struct {
	AddressPrefix string
}

func NewUtilWithPrefix(addressPrefix string) *Util {
	return &Util{addressPrefix}
}

// MARK - Implement the protocol wallet.Util

// @param publicKey can start with 0x or not.
func (u *Util) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return EncodePublicKeyToAddress(publicKey, u.AddressPrefix)
}

// @throw decode address to public key is not supported
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	return "", errors.New("decode address to public key is not supported")
}

func (u *Util) IsValidAddress(address string) bool {
	return IsValidAddress(address, u.AddressPrefix)
}

// MARK - like wallet.Util

// @param publicKey can start with 0x or not.
func EncodePublicKeyToAddress(publicKey string, addressPrefix string) (string, error) {
	pubBytes, err := types.HexDecodeString(publicKey)
	if err != nil {
		return "", err
	}
	pubKey := secp256k1.PubKey{Key: pubBytes}
	return Bech32FromAccAddress(pubKey.Address().Bytes(), addressPrefix)
}

// @param chainnet chain name
func IsValidAddress(address string, addressPrefix string) bool {
	_, err := AccAddressFromBech32(address, addressPrefix)
	return err == nil
}

func AccAddressFromBech32(address string, addressPrefix string) (sdk.AccAddress, error) {
	if len(strings.TrimSpace(address)) == 0 {
		return nil, errors.New("empty address string is not allowed")
	}

	bz, err := sdk.GetFromBech32(address, addressPrefix)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.AccAddress(bz), nil
}

func Bech32FromAccAddress(address []byte, addressPrefix string) (string, error) {
	return bech32.ConvertAndEncode(addressPrefix, address)
}
