package btc

import (
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/go-bip39"
)

type Account struct {
	privateKey string
	publicKey  string
	address    string
	util       *Util
}

func NewAccountWithMnemonic(mnemonic, chainnet string) (*Account, error) {
	util, err := NewUtilWithChainnet(chainnet)
	if err != nil {
		return nil, err
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	pri, pub := btcec.PrivKeyFromBytes(seed)
	pubData := pub.SerializeUncompressed()
	privateKey := types.HexEncodeToString(pri.Serialize())
	publicKey := types.HexEncodeToString(pubData)
	address, err := util.EncodePublicDataToAddress(pubData)
	if err != nil {
		return nil, err
	}

	return &Account{
		privateKey: privateKey,
		publicKey:  publicKey,
		address:    address,
		util:       util,
	}, nil
}

// MARK - Implement the protocol wallet.Account

// If the account generated using keystore, it will return empty
// @return privateKey that will start with 0x.
func (a *Account) PrivateKey() string {
	return a.privateKey
}

// @return publicKey that will start with 0x.
func (a *Account) PublicKey() string {
	return a.publicKey
}

func (a *Account) Address() string {
	return a.address
}

// TODO: function not implement yet.
func (a *Account) SignData(data []byte, password string) (string, error) {
	return "", errors.New("TODO: function not implement yet.")
}

// TODO: function not implement yet.
func (a *Account) SignHexData(hex string, password string) (string, error) {
	return "", errors.New("TODO: function not implement yet.")
}

// Only available to accounts generated with keystore.
// @return If the password is correct, will return nil
func (a *Account) CheckPassword(password string) error {
	return nil
}

// MARK - Implement the protocol wallet.Util

// @param publicKey can start with 0x or not.
func (a *Account) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return a.util.EncodePublicKeyToAddress(publicKey)
}

// Warning: Btc cannot support decode address to public key
func (a *Account) DecodeAddressToPublicKey(address string) (string, error) {
	return a.util.DecodeAddressToPublicKey(address)
}

func (a *Account) IsValidAddress(address string) bool {
	return a.util.IsValidAddress(address)
}
