package doge

import (
	"errors"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Util struct {
	Chainnet string
}

func NewUtilWithChainnet(chainnet string) (*Util, error) {
	if isValidChain(chainnet) {
		return &Util{Chainnet: chainnet}, nil
	} else {
		return nil, ErrUnsupportedChain
	}
}

// MARK - Implement the protocol Util

// @param publicKey can start with 0x or not.
func (u *Util) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return EncodePublicKeyToAddress(publicKey, u.Chainnet)
}

func (u *Util) EncodePublicDataToAddress(public []byte) (string, error) {
	return EncodePublicDataToAddress(public, u.Chainnet)
}

// Warning: Dogecoin cannot support decode address to public key
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	return "", errors.New("Dogecoin cannot support decode address to public key")
}

func (u *Util) IsValidAddress(address string) bool {
	return IsValidAddress(address, u.Chainnet)
}

// MARK - like Util

// @param publicKey can start with 0x or not.
func EncodePublicKeyToAddress(publicKey string, chainnet string) (string, error) {
	pubData, err := types.HexDecodeString(publicKey)
	if err != nil {
		return "", err
	}
	return EncodePublicDataToAddress(pubData, chainnet)
}

func EncodePublicDataToAddress(public []byte, chainnet string) (string, error) {
	net, err := netParamsOf(chainnet)
	if err != nil {
		return "", err
	}
	address, err := btcutil.NewAddressPubKey(public, net)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

func IsValidAddress(address, chainnet string) bool {
	net, err := netParamsOf(chainnet)
	if err != nil {
		return false
	}
	_, err = btcutil.DecodeAddress(address, net)
	if err != nil {
		return false
	}
	return true
}
