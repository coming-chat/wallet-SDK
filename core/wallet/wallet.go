package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/coming-chat/wallet-SDK/core/aptos"
	"github.com/coming-chat/wallet-SDK/core/btc"
	"github.com/coming-chat/wallet-SDK/core/cosmos"
	"github.com/coming-chat/wallet-SDK/core/doge"
	"github.com/coming-chat/wallet-SDK/core/eth"
	"github.com/coming-chat/wallet-SDK/core/polka"
	"github.com/coming-chat/wallet-SDK/core/solana"
	"github.com/coming-chat/wallet-SDK/core/starcoin"
	"github.com/coming-chat/wallet-SDK/core/starknet"
	"github.com/coming-chat/wallet-SDK/core/sui"
)

// Deprecated: 这个钱包对象缓存了助记词、密码、私钥等信息，继续使用有泄露资产的风险 ⚠️
type Wallet struct {
	Mnemonic string

	Keystore string
	password string

	// cache
	multiAccounts   sync.Map
	ethereumAccount *eth.Account
	solanaAccount   *solana.Account
	aptosAccount    *aptos.Account
	suiAccount      *sui.Account
	starcoinAccount *starcoin.Account
	starknetAccount *starknet.Account

	WatchAddress string
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

func WatchWallet(address string) (*Wallet, error) {
	chainType := ChainTypeFrom(address)
	if chainType.Count() == 0 {
		return nil, errors.New("Invalid wallet address")
	}
	return &Wallet{WatchAddress: address}, nil
}

func (w *Wallet) IsMnemonicWallet() bool {
	return len(w.Mnemonic) > 0
}

func (w *Wallet) IsKeystoreWallet() bool {
	return len(w.Keystore) > 0
}

func (w *Wallet) IsWatchWallet() bool {
	return len(w.WatchAddress) > 0
}

func (w *Wallet) GetWatchWallet() *WatchAccount {
	return &WatchAccount{address: w.WatchAddress}
}

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

// Get or create the bitcoin account with specified chainnet.
func (w *Wallet) GetOrCreateBitcoinAccount(chainnet string, addressType btc.AddressType) (*btc.Account, error) {
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
	account.AddressType = addressType

	// save to cache
	w.multiAccounts.Store(key, account)
	return account, nil
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
	return account, nil
}

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

// Get or create cosmos chain account
func (w *Wallet) GetOrCreateCosmosAccount() (*cosmos.Account, error) {
	return w.GetOrCreateCosmosTypeAccount(cosmos.CosmosCointype, cosmos.CosmosPrefix)
}

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

// Get or create the starknet account.
func (w *Wallet) GetOrCreateStarknetAccount(isCairo0 bool) (*starknet.Account, error) {
	cache := w.starknetAccount
	if cache != nil {
		cache.Cairo0 = isCairo0
		return cache, nil
	}
	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := starknet.NewAccountWithMnemonic(w.Mnemonic)
	if err != nil {
		return nil, err
	}
	account.Cairo0 = isCairo0
	// save to cache
	w.starknetAccount = account
	return account, nil
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
	bytes, err := hex.DecodeString(messageHex)
	if err != nil {
		return nil, err
	}
	return account.Sign(bytes, password)
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
