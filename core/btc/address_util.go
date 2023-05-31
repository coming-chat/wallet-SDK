package btc

import (
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type Util struct {
	Chainnet string
}

func NewUtilWithChainnet(chainnet string) (*Util, error) {
	if isValidChain(chainnet) {
		return &Util{chainnet}, nil
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

// Warning: Btc cannot support decode address to public key
func (u *Util) DecodeAddressToPublicKey(address string) (string, error) {
	return "", ErrDecodeAddress
}

func (u *Util) IsValidAddress(address string) bool {
	return IsValidAddress(address, u.Chainnet)
}

// MARK - like Util

// @param publicKey can start with 0x or not.
func EncodePublicKeyToAddress(publicKey, chainnet string) (string, error) {
	pubData, err := types.HexDecodeString(publicKey)
	if err != nil {
		return "", err
	}
	return EncodePublicDataToAddress(pubData, chainnet)
}

func EncodePublicDataToAddress(public []byte, chainnet string) (string, error) {
	params, err := netParamsOf(chainnet)
	if err != nil {
		return "", err
	}
	segwitAddress, err := btcutil.NewAddressTaproot(public[1:33], params)
	if err != nil {
		return "", err
	}
	return segwitAddress.EncodeAddress(), nil
}

// @param chainnet chain name
func IsValidAddress(address, chainnet string) bool {
	for _, ch := range []byte(address) {
		valid := (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
		if !valid {
			return false
		}
	}
	params, err := netParamsOf(chainnet)
	if err != nil {
		return false
	}
	_, err = btcutil.DecodeAddress(address, params)
	if err != nil {
		return false
	}
	return true
}

func IsValidPrivateKey(prikey string) bool {
	if strings.HasPrefix(prikey, "0x") || strings.HasPrefix(prikey, "0X") {
		prikey = prikey[2:] // remove 0x prefix
	}
	valid, length := base.IsValidHexString(prikey)
	if valid && length == 64 {
		return true
	}
	_, err := btcutil.DecodeWIF(prikey)
	return err == nil
}
