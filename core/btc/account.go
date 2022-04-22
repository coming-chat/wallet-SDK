package btc

import (
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/go-bip39"
)

type Account struct {
	privateKey []byte
	publicKey  string
	address    string
	Chainnet   string
}

func NewAccountWithMnemonic(mnemonic, chainnet string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	pri, pub := btcec.PrivKeyFromBytes(seed)
	priData := pri.Serialize()
	pubData := pub.SerializeUncompressed()
	pubHex := types.HexEncodeToString(pubData)

	address, err := EncodePublicDataToAddress(pubData, chainnet)
	if err != nil {
		return nil, err
	}

	return &Account{
		privateKey: priData,
		publicKey:  pubHex,
		address:    address,
		Chainnet:   chainnet,
	}, nil
}

func (a *Account) DeriveAccountAt(chainnet string) (*Account, error) {
	address, err := EncodePublicKeyToAddress(a.publicKey, chainnet)
	if err != nil {
		return nil, err
	}
	return &Account{
		privateKey: a.privateKey,
		publicKey:  a.publicKey,
		address:    address,
		Chainnet:   chainnet,
	}, nil
}

// MARK - Implement the protocol wallet.Account

// @return privateKey data
func (a *Account) PrivateKeyData() ([]byte, error) {
	return a.privateKey, nil
}

// @return privateKey string that will start with 0x.
func (a *Account) PrivateKey() (string, error) {
	return types.HexEncodeToString(a.privateKey), nil
}

// @return publicKey string that will start with 0x.
func (a *Account) PublicKey() string {
	return a.publicKey
}

// @return default is the mainnet address
func (a *Account) Address() string {
	return a.address
}

// TODO: function not implement yet.
func (a *Account) Sign(message []byte, password string) ([]byte, error) {
	return nil, errors.New("TODO: function not implement yet.")
}

// TODO: function not implement yet.
func (a *Account) SignHex(messageHex string, password string) ([]byte, error) {
	return nil, errors.New("TODO: function not implement yet.")
}
