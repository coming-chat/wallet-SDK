package wallet

import (
	"fmt"

	"github.com/coming-chat/wallet-SDK/core/btc"
	"github.com/coming-chat/wallet-SDK/core/cosmos"
	"github.com/coming-chat/wallet-SDK/core/eth"
	"github.com/coming-chat/wallet-SDK/core/polka"
)

type Wallet struct {
	Mnemonic string

	Keystore string
	password string

	// cache
	polkaAccounts   map[int]*polka.Account
	bitcoinAccounts map[string]*btc.Account
	ethereumAccount *eth.Account
	cosmosAccounts  map[string]*cosmos.Account
}

func NewWalletWithMnemonic(mnemonic string) (*Wallet, error) {
	if !IsValidMnemonic(mnemonic) {
		return nil, ErrInvalidMnemonic
	}
	return &Wallet{Mnemonic: mnemonic}, nil
}

// Only support Polka keystore.
func NewWalletWithKeyStore(keyStoreJson string, password string) (*Wallet, error) {
	// check keystore's password
	err := polka.CheckKeystorePassword(keyStoreJson, password)
	if err != nil {
		return nil, err
	}
	return &Wallet{
		Keystore: keyStoreJson,
		password: password,
	}, nil
}

// Get or create the polka account with specified network.
func (w *Wallet) GetOrCreatePolkaAccount(network int) (*polka.Account, error) {
	if w.polkaAccounts == nil {
		w.polkaAccounts = make(map[int]*polka.Account)
	}

	cache := w.polkaAccounts[network]
	if cache != nil {
		return cache, nil
	}

	var account *polka.Account
	var err error
	if len(w.Mnemonic) > 0 {
		account, err = polka.NewAccountWithMnemonic(w.Mnemonic, network)
	} else if len(w.Keystore) > 0 {
		account, err = polka.NewAccountWithKeystore(w.Keystore, w.password, network)
	}
	if err != nil {
		return nil, err
	}
	// save to cache
	w.polkaAccounts[network] = account
	return account, err
}

// Get or create the bitcoin account with specified chainnet.
func (w *Wallet) GetOrCreateBitcoinAccount(chainnet string) (*btc.Account, error) {
	if w.bitcoinAccounts == nil {
		w.bitcoinAccounts = make(map[string]*btc.Account)
	}

	cache := w.bitcoinAccounts[chainnet]
	if cache != nil {
		return cache, nil
	}

	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := btc.NewAccountWithMnemonic(w.Mnemonic, chainnet)
	if err != nil {
		return nil, err
	}

	// save to cache
	w.bitcoinAccounts[chainnet] = account
	return account, err
}

// Get or create the ethereum account.
func (w *Wallet) GetOrCreateEthereumAccount() (*eth.Account, error) {
	cache := w.ethereumAccount
	if cache != nil {
		return cache, nil
	}
	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := eth.NewAccountWithMnemonic(w.Mnemonic)
	if err != nil {
		return nil, err
	}
	// save to cache
	w.ethereumAccount = account
	return account, err
}

// Get or create a wallet account based on cosmos architecture.
func (w *Wallet) GetOrCreateCosmosTypeAccount(cointype int64, addressPrefix string) (*cosmos.Account, error) {
	key := fmt.Sprintf("%d-%s", cointype, addressPrefix)
	if w.cosmosAccounts == nil {
		w.cosmosAccounts = make(map[string]*cosmos.Account)
	}

	cache := w.cosmosAccounts[key]
	if cache != nil {
		return cache, nil
	}

	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := cosmos.NewAccountWithMnemonic(w.Mnemonic, cointype, addressPrefix)
	if err != nil {
		return nil, err
	}

	// save to cache
	w.cosmosAccounts[key] = account
	return account, err
}

// Get or create cosmos chain account
func (w *Wallet) GetOrCreateCosmosAccount() (*cosmos.Account, error) {
	return w.GetOrCreateCosmosTypeAccount(cosmos.CosmosCointype, cosmos.CosmosPrefix)
}

// check keystore password
func (w *Wallet) CheckPassword(password string) (bool, error) {
	err := polka.CheckKeystorePassword(w.Keystore, w.password)
	return err == nil, err
}

// Deprecated: Sign is deprecated. Please use wallet.PolkaAccount(network).Sign() instead
func (w *Wallet) Sign(message []byte, password string) (b []byte, err error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return nil, err
	}
	return account.Sign(message, password)
}

// Deprecated: SignFromHex is deprecated. Please use wallet.PolkaAccount(network).SignHex() instead
func (w *Wallet) SignFromHex(messageHex string, password string) ([]byte, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return nil, err
	}
	return account.SignHex(messageHex, password)
}

// Deprecated: GetPublicKey is deprecated. Please use wallet.PolkaAccount(network).PublicKey() instead
func (w *Wallet) GetPublicKey() ([]byte, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return nil, err
	}
	return account.PublicKey(), nil
}

// Deprecated: GetPublicKeyHex is deprecated. Please use wallet.PolkaAccount(network).PublicKey() instead
func (w *Wallet) GetPublicKeyHex() (string, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return "", err
	}
	return account.PublicKeyHex(), nil
}

// Deprecated: GetAddress is deprecated. Please use wallet.PolkaAccount(network).Address() instead
func (w *Wallet) GetAddress(network int) (string, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return "", err
	}
	return account.Address(), nil
}

// Deprecated: GetPrivateKeyHex is deprecated. Please use wallet.PolkaAccount(network).PrivateKey() instead
func (w *Wallet) GetPrivateKeyHex() (string, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return "", err
	}
	return account.PrivateKeyHex()
}
