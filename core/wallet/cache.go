package wallet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/coming-chat/wallet-SDK/core/base"
)

var walletCache = sync.Map{}

func SaveWallet(wallet *CacheWallet) {
	key := "wallet-" + wallet.WalletId
	walletCache.Store(key, wallet)
}

func GetWallet(walletId string) *CacheWallet {
	key := "wallet-" + walletId
	if c, ok := walletCache.Load(key); ok {
		return c.(*CacheWallet)
	}
	return nil
}

func SaveAccountInfo(walletId string, info *AccountInfo) {
	key := fmt.Sprintf("account-%v-%v", walletId, info.cacheKey)
	walletCache.Store(key, info)
}

func GetAccountInfo(walletId, cacheKey string) *AccountInfo {
	key := fmt.Sprintf("account-%v-%v", walletId, cacheKey)
	if c, ok := walletCache.Load(key); ok {
		return c.(*AccountInfo)
	}
	return nil
}

type accountCreator = func(val string) (base.Account, error)
type AccountInfo struct {
	wallet          *CacheWallet
	cacheKey        string
	mnemonicCreator accountCreator
	keystoreCreator accountCreator
	privkeyCreator  accountCreator

	// Chain string
	publicKey []byte
	address   string
}

func (i *AccountInfo) Account() (base.Account, error) {
	typ, ok := i.wallet.checkWalletType()
	var val = ""
	if ok {
		typ, val = readValue(i.wallet.WalletId, typ)
	} else {
		typ, val = readTypeAndValue(i.wallet.WalletId)
	}
	var account base.Account
	var err error
	switch typ {
	case WalletTypeMnemonic:
		if i.mnemonicCreator == nil {
			return nil, errors.New("Lose function of mnemonic account creator")
		}
		account, err = i.mnemonicCreator(val)
	case WalletTypeKeystore:
		if i.keystoreCreator == nil {
			return nil, ErrUnsupportKeystore
		}
		account, err = i.keystoreCreator(val)
	case WalletTypePrivateKey:
		if i.privkeyCreator == nil {
			return nil, errors.New("Lose function of private key account creator")
		}
		account, err = i.privkeyCreator(val)
	case WalletTypeWatch:
		return &WatchAccount{address: val}, nil
	default:
		return nil, ErrWalletInfoNotExist
	}
	if err != nil {
		return nil, err
	}
	i.loadAndSaveCacheIfNotFound(account)
	return account, nil
}

func (i *AccountInfo) PrivateKeyHex() (*base.OptionalString, error) {
	account, err := i.Account()
	if err != nil {
		return nil, err
	}
	h, err := account.PrivateKeyHex()
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: h}, nil
}

func (i *AccountInfo) PublicKeyHex() (*base.OptionalString, error) {
	if i.publicKey == nil {
		new, err := i.loadAndSaveCacheIfNotFound(nil)
		if err != nil {
			return nil, err
		}
		i.publicKey = new.publicKey
	}
	return &base.OptionalString{Value: ByteToHex(i.publicKey)}, nil
}

func (i *AccountInfo) Address() (*base.OptionalString, error) {
	if i.address == "" {
		new, err := i.loadAndSaveCacheIfNotFound(nil)
		if err != nil {
			return nil, err
		}
		i.address = new.address
	}
	return &base.OptionalString{Value: i.address}, nil
}

func (i *AccountInfo) loadAndSaveCacheIfNotFound(account base.Account) (*AccountInfo, error) {
	if cache := GetAccountInfo(i.wallet.WalletId, i.cacheKey); cache != nil {
		return cache, nil
	}
	if account == nil {
		var err error
		account, err = i.Account()
		if err != nil {
			return nil, err
		}
	}
	i.publicKey = account.PublicKey()
	i.address = account.Address()
	SaveAccountInfo(i.wallet.WalletId, i)
	return i, nil
}
