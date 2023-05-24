package solana

import (
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/portto/solana-go-sdk/pkg/hdwallet"
	solana "github.com/portto/solana-go-sdk/types"
	"github.com/tyler-smith/go-bip39"
)

type Account struct {
	account *solana.Account
}

func NewAccountWithMnemonic(mnemonic string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	path := "m/44'/501'/0'/0'"
	derivedKey, err := hdwallet.Derived(path, seed)
	account, err := solana.AccountFromSeed(derivedKey.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &Account{&account}, nil
}

// rename for support android.
// Android cant support both NewAccountWithMnemonic(string) and NewAccountWithPrivateKey(string)
func AccountWithPrivateKey(prikey string) (*Account, error) {
	prikey = strings.TrimPrefix(prikey, "0x")
	account, err := solana.AccountFromHex(prikey)
	if err != nil {
		return nil, err
	}
	return &Account{&account}, nil
}

// MARK - Implement the protocol Account

// @return privateKey data
func (a *Account) PrivateKey() ([]byte, error) {
	return a.account.PrivateKey, nil
}

// @return privateKey string that will start with 0x.
func (a *Account) PrivateKeyHex() (string, error) {
	return types.HexEncodeToString(a.account.PrivateKey), nil
}

// @return publicKey data
func (a *Account) PublicKey() []byte {
	return a.account.PublicKey.Bytes()
}

// @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	return types.HexEncodeToString(a.account.PublicKey.Bytes())
}

func (a *Account) Address() string {
	return a.account.PublicKey.ToBase58()
}

func (a *Account) Sign(message []byte, password string) ([]byte, error) {
	return a.account.Sign(message), nil
}

func (a *Account) SignHex(messageHex string, password string) (*base.OptionalString, error) {
	message, err := types.HexDecodeString(messageHex)
	if err != nil {
		return nil, err
	}
	res := types.HexEncodeToString(a.account.Sign(message))
	return &base.OptionalString{Value: res}, nil
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

func AsSolanaAccount(account base.Account) *Account {
	if r, ok := account.(*Account); ok {
		return r
	} else {
		return nil
	}
}
