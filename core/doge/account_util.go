package doge

import (
	"errors"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var dogeCfg = chaincfg.MainNetParams

func getDogeCfg() *chaincfg.Params {
	dogeCfg.PubKeyHashAddrID = 0x1e
	dogeCfg.ScriptHashAddrID = 0x16
	dogeCfg.PrivateKeyID = 0x9e
	return &dogeCfg
}

type Util struct {
}

func NewUtilWithChainnet() *Util {
	return &Util{}
}

// MARK - Implement the protocol Util

// @param publicKey can start with 0x or not.
func (u *Util) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return EncodePublicKeyToAddress(publicKey)
}

func (u *Util) EncodePublicDataToAddress(public []byte) (string, error) {
	return EncodePublicDataToAddress(public)
}

// Warning: Dogecoin cannot support decode address to public key
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	return "", errors.New("Dogecoin cannot support decode address to public key")
}

func (u *Util) IsValidAddress(address string) bool {
	return IsValidAddress(address)
}

// MARK - like Util

// @param publicKey can start with 0x or not.
func EncodePublicKeyToAddress(publicKey string) (string, error) {
	pubData, err := types.HexDecodeString(publicKey)
	if err != nil {
		return "", err
	}
	return EncodePublicDataToAddress(pubData)
}

func EncodePublicDataToAddress(public []byte) (string, error) {
	address, err := btcutil.NewAddressPubKey(public, getDogeCfg())
	// address, err := btcutil.NewAddressPubKey(public, getDogeCfg()) 		  			// DJhF8ahvTfGhqcLEn7sN4gJMJVVbmfwxkU
	// address, err := btcutil.NewAddressPubKeyHash(public[1:21], getDogeCfg()) 		// DLNp3hjacupjoLMEawnwPcfggRoRD1RykH
	// address, err := btcutil.NewAddressScriptHash(public, getDogeCfg()) 				// A5zRFiKcnDZhJ9DYamCpBg84GTS453LN8cs
	// address, err := btcutil.NewAddressScriptHashFromHash(public[1:21], getDogeCfg()) // A7fzAqMGwU7jFsEYPb8PWcVPePjsXbQHkv
	// address, err := btcutil.NewAddressTaproot(public[1:33], getDogeCfg())
	// address, err := btcutil.NewAddressWitnessPubKeyHash(public[1:21], getDogeCfg()) 	// bc1q5uslzuqy8k40mc86jfdtdjh4624umtwj9lcr8a
	// address, err := btcutil.NewAddressWitnessScriptHash(public[1:33], getDogeCfg()) 	// bc1q5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9yskyzcg5
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

func IsValidAddress(address string) bool {
	_, err := btcutil.DecodeAddress(address, getDogeCfg())
	if err != nil {
		return false
	}
	return true
}
