package sui

import (
	"encoding/hex"
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/go-sui/account"
	"github.com/coming-chat/go-sui/sui_types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type Account struct {
	account *account.Account
}

func NewAccountWithMnemonic(mnemonic string) (*Account, error) {
	account, err := account.NewAccountWithMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	return &Account{account: account}, nil
}

// rename for support android.
// Android cannot support both NewAccountWithMnemonic(string) and NewAccountWithPrivateKey(string)
func AccountWithPrivateKey(prikey string) (*Account, error) {
	seed, err := types.HexDecodeString(prikey)
	if err != nil {
		return nil, err
	}
	scheme, err := sui_types.NewSignatureScheme(0)
	if err != nil {
		return nil, err
	}
	account := account.NewAccount(scheme, seed)
	return &Account{account: account}, nil
}

// MARK - Implement the protocol Account

// @return privateKey data
func (a *Account) PrivateKey() ([]byte, error) {
	return a.account.KeyPair.PrivateKey()[:32], nil
}

// @return privateKey string that will start with 0x.
func (a *Account) PrivateKeyHex() (string, error) {
	return types.HexEncodeToString(a.account.KeyPair.PrivateKey()[:32]), nil
}

// @return publicKey data
func (a *Account) PublicKey() []byte {
	return a.account.KeyPair.PublicKey()
}

// @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	return types.HexEncodeToString(a.account.KeyPair.PublicKey())
}

func (a *Account) Address() string {
	return a.account.Address
}

func (a *Account) Sign(message []byte, password string) ([]byte, error) {
	return a.account.Sign(message), nil
}

func (a *Account) SignHex(messageHex string, password string) (*base.OptionalString, error) {
	msg, err := types.HexDecodeString(messageHex)
	if err != nil {
		return nil, errors.New("Invalid message hex string")
	}
	signature := a.account.Sign(msg)
	signString := hex.EncodeToString(signature)
	return &base.OptionalString{Value: signString}, nil
}

// MARK - Implement the protocol AddressUtil

// @param publicKey can start with 0x or not.
func (a *Account) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return EncodePublicKeyToAddress(publicKey)
}

// @return publicKey that will start with 0x.
func (a *Account) DecodeAddressToPublicKey(address string) (string, error) {
	return DecodeAddressToPublicKey(address)
}

func (a *Account) IsValidAddress(address string) bool {
	return IsValidAddress(address)
}

func AsSuiAccount(account base.Account) *Account {
	if r, ok := account.(*Account); ok {
		return r
	} else {
		return nil
	}
}
