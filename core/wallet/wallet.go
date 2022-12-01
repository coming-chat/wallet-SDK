package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/coming-chat/wallet-SDK/core/aptos"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/btc"
	"github.com/coming-chat/wallet-SDK/core/cosmos"
	"github.com/coming-chat/wallet-SDK/core/doge"
	"github.com/coming-chat/wallet-SDK/core/eth"
	"github.com/coming-chat/wallet-SDK/core/polka"
	"github.com/coming-chat/wallet-SDK/core/solana"
	"github.com/coming-chat/wallet-SDK/core/starcoin"
	"github.com/coming-chat/wallet-SDK/core/sui"
)

type WalletType = base.SDKEnumInt

const (
	WalletTypeMnemonic   WalletType = 1
	WalletTypeKeystore   WalletType = 2
	WalletTypePrivateKey WalletType = 3
	WalletTypeError      WalletType = 4
)

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

	WatchAddress string

	WalletId   string
	walletType WalletType
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

// check keystore password
func (w *Wallet) CheckPassword(password string) (bool, error) {
	err := polka.CheckKeystorePassword(w.Keystore, w.password)
	return err == nil, err
}

// Deprecated: Sign is deprecated. Please use wallet.PolkaAccount(network).Sign() instead
func (w *Wallet) Sign(message []byte, password string) (b []byte, err error) {
	account, err := w.PolkaAccountInfo(44).Account()
	if err != nil {
		return nil, err
	}
	return account.Sign(message, password)
}

// Deprecated: SignFromHex is deprecated. Please use wallet.PolkaAccount(network).SignHex() instead
func (w *Wallet) SignFromHex(messageHex string, password string) ([]byte, error) {
	account, err := w.PolkaAccountInfo(44).Account()
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
	account, err := w.PolkaAccountInfo(44).Account()
	if err != nil {
		return nil, err
	}
	return account.PublicKey(), nil
}

// Deprecated: GetPublicKeyHex is deprecated. Please use wallet.PolkaAccount(network).PublicKey() instead
func (w *Wallet) GetPublicKeyHex() (string, error) {
	account, err := w.PolkaAccountInfo(44).Account()
	if err != nil {
		return "", err
	}
	return account.PublicKeyHex(), nil
}

// Deprecated: GetAddress is deprecated. Please use wallet.PolkaAccount(network).Address() instead
func (w *Wallet) GetAddress(network int) (string, error) {
	account, err := w.PolkaAccountInfo(44).Account()
	if err != nil {
		return "", err
	}
	return account.Address(), nil
}

// Deprecated: GetPrivateKeyHex is deprecated. Please use wallet.PolkaAccount(network).PrivateKey() instead
func (w *Wallet) GetPrivateKeyHex() (string, error) {
	account, err := w.PolkaAccountInfo(44).Account()
	if err != nil {
		return "", err
	}
	return account.PrivateKeyHex()
}

func (w *Wallet) WalletType() WalletType {
	if typ, ok := w.checkWalletType(); ok {
		return typ
	}
	w.walletType, _ = readTypeAndValue(w.WalletId)
	SaveWallet(w)
	return w.walletType
}

func (w *Wallet) checkWalletType() (typ WalletType, ok bool) {
	if w.walletType >= WalletTypeMnemonic && w.walletType <= WalletTypePrivateKey {
		return w.walletType, true
	}
	if cache := GetWallet(w.WalletId); cache != nil &&
		cache.walletType >= WalletTypeMnemonic && cache.walletType <= WalletTypePrivateKey {
		w.walletType = cache.walletType
		return cache.walletType, true
	}
	return WalletTypeError, false
}

func (w *Wallet) PolkaAccountInfo(network int) *AccountInfo {
	walletId := w.WalletId
	return &AccountInfo{
		wallet:   w,
		cacheKey: fmt.Sprintf("polka-%v", network),
		mnemonicCreator: func(val string) (base.Account, error) {
			return polka.NewAccountWithMnemonic(val, network)
		},
		keystoreCreator: func(val string) (base.Account, error) {
			pwd := InfoProvider.Password(walletId)
			return polka.NewAccountWithKeystore(val, pwd, network)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return polka.AccountWithPrivateKey(val, network)
		},
	}
}

func (w *Wallet) BitcoinAccountInfo(chainnet string) *AccountInfo {
	return &AccountInfo{
		wallet:   w,
		cacheKey: fmt.Sprintf("bitcoin-%v", chainnet),
		mnemonicCreator: func(val string) (base.Account, error) {
			return btc.NewAccountWithMnemonic(val, chainnet)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return btc.AccountWithPrivateKey(val, chainnet)
		},
	}
}

func (w *Wallet) EthereumAccountInfo() *AccountInfo {
	return &AccountInfo{
		wallet:   w,
		cacheKey: "ethereum",
		mnemonicCreator: func(val string) (base.Account, error) {
			return eth.NewAccountWithMnemonic(val)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return eth.AccountWithPrivateKey(val)
		},
	}
}

func (w *Wallet) CosmosAccountInfo(cointype int64, prefix string) *AccountInfo {
	return &AccountInfo{
		wallet:   w,
		cacheKey: fmt.Sprintf("cosmos-%v-%v", cointype, prefix),
		mnemonicCreator: func(val string) (base.Account, error) {
			return cosmos.NewAccountWithMnemonic(val, cointype, prefix)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return cosmos.AccountWithPrivateKey(val, cointype, prefix)
		},
	}
}

func (w *Wallet) DogecoinAccountInfo(chainnet string) *AccountInfo {
	return &AccountInfo{
		wallet:   w,
		cacheKey: fmt.Sprintf("bitcoin-%v", chainnet),
		mnemonicCreator: func(val string) (base.Account, error) {
			return doge.NewAccountWithMnemonic(val, chainnet)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return doge.AccountWithPrivateKey(val, chainnet)
		},
	}
}

func (w *Wallet) SolanaAccountInfo() *AccountInfo {
	return &AccountInfo{
		wallet:   w,
		cacheKey: "solana",
		mnemonicCreator: func(val string) (base.Account, error) {
			return solana.NewAccountWithMnemonic(val)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return solana.AccountWithPrivateKey(val)
		},
	}
}

func (w *Wallet) AptosAccountInfo() *AccountInfo {
	return &AccountInfo{
		wallet:   w,
		cacheKey: "aptos",
		mnemonicCreator: func(val string) (base.Account, error) {
			return aptos.NewAccountWithMnemonic(val)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return aptos.AccountWithPrivateKey(val)
		},
	}
}

func (w *Wallet) SuiAccountInfo() *AccountInfo {
	return &AccountInfo{
		wallet:   w,
		cacheKey: "sui",
		mnemonicCreator: func(val string) (base.Account, error) {
			return sui.NewAccountWithMnemonic(val)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return sui.AccountWithPrivateKey(val)
		},
	}
}

func (w *Wallet) StarcoinAccountInfo() *AccountInfo {
	return &AccountInfo{
		wallet:   w,
		cacheKey: "starcoin",
		mnemonicCreator: func(val string) (base.Account, error) {
			return starcoin.NewAccountWithMnemonic(val)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return starcoin.AccountWithPrivateKey(val)
		},
	}
}
