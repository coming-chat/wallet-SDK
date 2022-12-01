package wallet

import (
	"fmt"

	"github.com/coming-chat/wallet-SDK/core/aptos"
	"github.com/coming-chat/wallet-SDK/core/btc"
	"github.com/coming-chat/wallet-SDK/core/cosmos"
	"github.com/coming-chat/wallet-SDK/core/doge"
	"github.com/coming-chat/wallet-SDK/core/eth"
	"github.com/coming-chat/wallet-SDK/core/polka"
	"github.com/coming-chat/wallet-SDK/core/solana"
	"github.com/coming-chat/wallet-SDK/core/starcoin"
	"github.com/coming-chat/wallet-SDK/core/sui"
)

// Deprecated: 这个方法会缓存助记词、密码、私钥等信息，继续使用有泄露资产的风险 ⚠️
func NewWalletWithMnemonic(mnemonic string) (*Wallet, error) {
	if !IsValidMnemonic(mnemonic) {
		return nil, ErrInvalidMnemonic
	}
	return &Wallet{Mnemonic: mnemonic}, nil
}

// Deprecated: 这个方法会缓存助记词、密码、私钥等信息，继续使用有泄露资产的风险 ⚠️
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

// Deprecated: use `PolkaAccountInfo(network)`
// Get or create the polka account with specified network.
func (w *Wallet) GetOrCreatePolkaAccount(network int) (*polka.Account, error) {
	key := fmt.Sprintf("polka-%v", network)
	if cache, ok := w.multiAccounts.Load(key); ok {
		if acc, ok := cache.(*polka.Account); ok {
			return acc, nil
		}
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
	w.multiAccounts.Store(key, account)
	return account, nil
}

// Deprecated: use `BitcoinAccountInfo(chainnet)`
// Get or create the bitcoin account with specified chainnet.
func (w *Wallet) GetOrCreateBitcoinAccount(chainnet string) (*btc.Account, error) {
	key := "bitcoin-" + chainnet
	if cache, ok := w.multiAccounts.Load(key); ok {
		if acc, ok := cache.(*btc.Account); ok {
			return acc, nil
		}
	}

	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}
	account, err := btc.NewAccountWithMnemonic(w.Mnemonic, chainnet)
	if err != nil {
		return nil, err
	}

	// save to cache
	w.multiAccounts.Store(key, account)
	return account, nil
}

// Deprecated: use `EthereumAccountInfo()`
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
	return account, nil
}

// Deprecated: use `CosmosAccountInfo(cointype, prefix)`
// Get or create a wallet account based on cosmos architecture.
func (w *Wallet) GetOrCreateCosmosTypeAccount(cointype int64, addressPrefix string) (*cosmos.Account, error) {
	key := fmt.Sprintf("cosmos-%d-%s", cointype, addressPrefix)
	if cache, ok := w.multiAccounts.Load(key); ok {
		if acc, ok := cache.(*cosmos.Account); ok {
			return acc, nil
		}
	}

	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}
	account, err := cosmos.NewAccountWithMnemonic(w.Mnemonic, cointype, addressPrefix)
	if err != nil {
		return nil, err
	}

	// save to cache
	w.multiAccounts.Store(key, account)
	return account, nil
}

// Deprecated: use `DogecoinAccountInfo(chainnet)`
func (w *Wallet) GetOrCreateDogeAccount(chainnet string) (*doge.Account, error) {
	key := "Dogecoin" + chainnet
	if cache, ok := w.multiAccounts.Load(key); ok {
		if acc, ok := cache.(*doge.Account); ok {
			return acc, nil
		}
	}
	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}
	account, err := doge.NewAccountWithMnemonic(w.Mnemonic, chainnet)
	if err != nil {
		return nil, err
	}
	// save to cache
	w.multiAccounts.Store(key, account)
	return account, nil
}

// Deprecated: use `CosmosAccountInfo(cointype, prefix)`
// Get or create cosmos chain account
func (w *Wallet) GetOrCreateCosmosAccount() (*cosmos.Account, error) {
	return w.GetOrCreateCosmosTypeAccount(cosmos.CosmosCointype, cosmos.CosmosPrefix)
}

// Deprecated: use `SolanaAccountInfo()`
// Get or create the solana account.
func (w *Wallet) GetOrCreateSolanaAccount() (*solana.Account, error) {
	cache := w.solanaAccount
	if cache != nil {
		return cache, nil
	}
	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := solana.NewAccountWithMnemonic(w.Mnemonic)
	if err != nil {
		return nil, err
	}
	// save to cache
	w.solanaAccount = account
	return account, nil
}

// Deprecated: use `AptosAccountInfo()`
// Get or create the aptos account.
func (w *Wallet) GetOrCreateAptosAccount() (*aptos.Account, error) {
	cache := w.aptosAccount
	if cache != nil {
		return cache, nil
	}
	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := aptos.NewAccountWithMnemonic(w.Mnemonic)
	if err != nil {
		return nil, err
	}
	// save to cache
	w.aptosAccount = account
	return account, nil
}

// Deprecated: use `SuiAccountInfo()`
// Get or create the sui account.
func (w *Wallet) GetOrCreateSuiAccount() (*sui.Account, error) {
	cache := w.suiAccount
	if cache != nil {
		return cache, nil
	}
	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := sui.NewAccountWithMnemonic(w.Mnemonic)
	if err != nil {
		return nil, err
	}
	// save to cache
	w.suiAccount = account
	return account, nil
}

// Deprecated: use `StarcoinAccountInfo()`
// Get or create the starcoin account.
func (w *Wallet) GetOrCreateStarcoinAccount() (*starcoin.Account, error) {
	cache := w.starcoinAccount
	if cache != nil {
		return cache, nil
	}
	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := starcoin.NewAccountWithMnemonic(w.Mnemonic)
	if err != nil {
		return nil, err
	}
	// save to cache
	w.starcoinAccount = account
	return account, nil
}
