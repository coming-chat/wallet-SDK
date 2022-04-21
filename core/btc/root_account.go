package btc

import (
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/go-bip39"
)

type RootAccount struct {
	privateKey []byte
	publicKey  string
	address    string
}

func NewRootAccountWithMnemonic(mnemonic string) (*RootAccount, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	pri, pub := btcec.PrivKeyFromBytes(seed)
	priData := pri.Serialize()
	pubData := pub.SerializeUncompressed()
	pubHex := types.HexEncodeToString(pubData)

	// default is mainnet address
	segwitAddress, err := btcutil.NewAddressTaproot(pubData[1:33], &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	address := segwitAddress.EncodeAddress()

	return &RootAccount{
		privateKey: priData,
		publicKey:  pubHex,
		address:    address,
	}, nil
}

// MARK - Implement the protocol wallet.Account

// @return privateKey data
func (a *RootAccount) PrivateKeyData() ([]byte, error) {
	return a.privateKey, nil
}

// @return privateKey string that will start with 0x.
func (a *RootAccount) PrivateKey() (string, error) {
	return types.HexEncodeToString(a.privateKey), nil
}

// @return publicKey string that will start with 0x.
func (a *RootAccount) PublicKey() string {
	return a.publicKey
}

// @return default is the mainnet address
func (a *RootAccount) Address() string {
	return a.address
}

// TODO: function not implement yet.
func (a *RootAccount) SignData(data []byte, password string) (string, error) {
	return "", errors.New("TODO: function not implement yet.")
}

// TODO: function not implement yet.
func (a *RootAccount) SignHexData(hex string, password string) (string, error) {
	return "", errors.New("TODO: function not implement yet.")
}
