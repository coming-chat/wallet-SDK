package polka

import (
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
	return EncodePublicKeyToAddress(publicKey, u.Network)
}

// @return publicKey that will start with 0x.
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	return DecodeAddressToPublicKey(address)
}

func (u *Util) IsValidAddress(address string) bool {
	return IsValidAddress(address)
}

// MARK - like wallet.Util

// @param publicKey can start with 0x or not.
func EncodePublicKeyToAddress(publicKey string, network int) (string, error) {
	address := ss58.Encode(publicKey, network)
	if len(address) == 0 {
		return "", ErrPublicKey
	}
	return address, nil
}

// @return publicKey that will start with 0x.
func DecodeAddressToPublicKey(address string) (string, error) {
	ss58Format := base58.Decode(address)
	if len(ss58Format) == 0 {
		return "", ErrAddress
	}
	publicKey := ss58.Decode(address, int(ss58Format[0]))
	if len(publicKey) == 0 {
		return "", ErrAddress
	}
	return "0x" + publicKey, nil
}

// @param chainnet chain name
func IsValidAddress(address string) bool {
	_, err := DecodeAddressToPublicKey(address)
	return err == nil
}
