package eth

import (
	"github.com/ethereum/go-ethereum/common"
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

func (u *Util) IsValidAddress(address string) bool {
	return IsValidAddress(address)
}

func IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}
