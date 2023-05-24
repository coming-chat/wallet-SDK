package btc

import (
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/tyler-smith/go-bip39"
)

type Account struct {
	privateKey []byte
	publicKey  []byte
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

	address, err := EncodePublicDataToAddress(pubData, chainnet)
	if err != nil {
		return nil, err
	}

	return &Account{
		privateKey: priData,
		publicKey:  pubData,
		address:    address,
		Chainnet:   chainnet,
	}, nil
}

func AccountWithPrivateKey(prikey string, chainnet string) (*Account, error) {
	seed, err := types.HexDecodeString(prikey)
	if err != nil {
		return nil, err
	}

	pri, pub := btcec.PrivKeyFromBytes(seed)
	priData := pri.Serialize()
	pubData := pub.SerializeUncompressed()

	address, err := EncodePublicDataToAddress(pubData, chainnet)
	if err != nil {
		return nil, err
	}

	return &Account{
		privateKey: priData,
		publicKey:  pubData,
		address:    address,
		Chainnet:   chainnet,
	}, nil
}

func (a *Account) DeriveAccountAt(chainnet string) (*Account, error) {
	address, err := EncodePublicDataToAddress(a.publicKey, chainnet)
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

// MARK - Implement the protocol Account

// @return privateKey data
func (a *Account) PrivateKey() ([]byte, error) {
	return a.privateKey, nil
}

// @return privateKey string that will start with 0x.
func (a *Account) PrivateKeyHex() (string, error) {
	return types.HexEncodeToString(a.privateKey), nil
}

// @return publicKey data
func (a *Account) PublicKey() []byte {
	return a.publicKey
}

// @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	return types.HexEncodeToString(a.publicKey)
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
func (a *Account) SignHex(messageHex string, password string) (*base.OptionalString, error) {
	return nil, errors.New("TODO: function not implement yet.")
}

// MARK - Implement the protocol AddressUtil

// @param publicKey can start with 0x or not.
func (a *Account) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return EncodePublicKeyToAddress(publicKey, a.Chainnet)
}

// @return publicKey that will start with 0x.
func (a *Account) DecodeAddressToPublicKey(address string) (string, error) {
	return "", ErrDecodeAddress
}

func (a *Account) IsValidAddress(address string) bool {
	return IsValidAddress(address, a.Chainnet)
}

func AsBitcoinAccount(account base.Account) *Account {
	if r, ok := account.(*Account); ok {
		return r
	} else {
		return nil
	}
}
