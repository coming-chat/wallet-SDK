package solana

import (
	"github.com/blocto/solana-go-sdk/common"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
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
	return DecodeAddressToPublicKey(address)
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
	pubKey := common.PublicKeyFromBytes(bytes)
	return pubKey.ToBase58(), nil
}

func DecodeAddressToPublicKey(address string) (string, error) {
	pubKey := common.PublicKeyFromString(address)
	return types.HexEncodeToString(pubKey.Bytes()), nil
}

func IsValidAddress(address string) bool {
	if len(address) != 44 {
		return false
	}
	pub := common.PublicKeyFromString(address)
	addr2 := pub.ToBase58()
	return address == addr2
}
