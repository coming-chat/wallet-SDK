package polka

import (
	"github.com/coming-chat/wallet-SDK/core/wallet"
	"github.com/decred/base58"
	"github.com/itering/subscan/util/ss58"
)

type Util struct {
	Network int
}

func NewUtilWithNetwork(network int) *Util {
	return &Util{Network: network}
}

// MARK - Implement the protocol wallet.Util

// @param publicKey can start with 0x or not.
func (u *Util) EncodePublicKeyToAddress(publicKey string) (string, error) {
	address := ss58.Encode(publicKey, u.Network)
	if len(address) == 0 {
		return "", wallet.ErrPublicKey
	}
	return address, nil
}

// @return publicKey that will start with 0x.
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	ss58Format := base58.Decode(address)
	if len(ss58Format) == 0 {
		return "", wallet.ErrAddress
	}
	publicKey := ss58.Decode(address, int(ss58Format[0]))
	if len(publicKey) == 0 {
		return "", wallet.ErrAddress
	}
	return "0x" + publicKey, nil
}

func (u *Util) IsValidAddress(address string) bool {
	_, err := u.DecodeAddressToPublicKey(address)
	return err == nil
}
