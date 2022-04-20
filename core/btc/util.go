package btc

import (
	"errors"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

const (
	ChainSignet  = "signet"
	ChainMainnet = "mainnet"
	ChainBitcoin = "bitcoin" // ComingChat local name
)

type Util struct {
	Chainnet string
}

func NewUtilWithChainnet(chainnet string) (*Util, error) {
	switch chainnet {
	case ChainSignet, ChainMainnet, ChainBitcoin:
		return &Util{Chainnet: chainnet}, nil
	default:
		return nil, ErrUnsupportedChain
	}
}

// MARK - Implement the protocol wallet.Util

// @param publicKey can start with 0x or not.
func (u *Util) EncodePublicKeyToAddress(publicKey string) (string, error) {
	pubData, err := types.HexDecodeString(publicKey)
	if err != nil {
		return "", err
	}
	return u.EncodePublicDataToAddress(pubData)
}

func (u *Util) EncodePublicDataToAddress(public []byte) (string, error) {
	segwitAddress, err := btcutil.NewAddressTaproot(public[1:33], u.net())
	if err != nil {
		return "", err
	}
	return segwitAddress.EncodeAddress(), nil
}

// Warning: Btc cannot support decode address to public key
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	return "", errors.New("Btc cannot support decode address to public key")
}

func (u *Util) IsValidAddress(address string) bool {
	_, err := btcutil.DecodeAddress(address, u.net())
	return err == nil
}

func (u *Util) net() *chaincfg.Params {
	switch u.Chainnet {
	case ChainSignet:
		return &chaincfg.SigNetParams
	case ChainMainnet, ChainBitcoin:
		return &chaincfg.MainNetParams
	}
	return nil
}

// @param chainnet chain name
func IsValidAddress(address, chainnet string) bool {
	u, err := NewUtilWithChainnet(chainnet)
	if err != nil {
		return false
	}
	return u.IsValidAddress(address)
}
