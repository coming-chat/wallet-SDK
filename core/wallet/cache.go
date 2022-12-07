package wallet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/coming-chat/wallet-SDK/core/base"
)

var walletCache = sync.Map{}

func saveWallet(wallet *CacheWallet) {
	key := "wallet-" + wallet.key()
	walletCache.Store(key, wallet)
}

func getWallet(walletKey string) *CacheWallet {
	key := "wallet-" + walletKey
	if c, ok := walletCache.Load(key); ok {
		return c.(*CacheWallet)
	}
	return nil
}

func saveAccountInfo(walletKey string, info *AccountInfo) {
	key := fmt.Sprintf("account-%v-%v", walletKey, info.cacheKey)
	walletCache.Store(key, info)
}

func getAccountInfo(walletKey, cacheKey string) *AccountInfo {
	key := fmt.Sprintf("account-%v-%v", walletKey, cacheKey)
	if c, ok := walletCache.Load(key); ok {
		return c.(*AccountInfo)
	}
	return nil
}

type accountCreator = func(val string) (base.Account, error)
type AccountInfo struct {
	Wallet          *CacheWallet
	cacheKey        string
	mnemonicCreator accountCreator
	keystoreCreator accountCreator
	privkeyCreator  accountCreator

	// Chain string
	publicKey []byte
	address   string
}

// 获取账号对象，该对象可以用来签名. 账号不会缓存，每次都会重新生成
func (i *AccountInfo) Account() (base.Account, error) {
	typ, ok := i.Wallet.checkWalletType()
	var val = ""
	if ok {
		typ, val = i.Wallet.readValue(typ)
	} else {
		typ, val = i.Wallet.readTypeAndValue()
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
	i.saveCache(account)
	return account, nil
}

// 获取账号私钥，私钥不会缓存，每次都会重新生成
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

// 获取公钥，该方法会优先读取缓存
func (i *AccountInfo) PublicKeyHex() (*base.OptionalString, error) {
	if i.publicKey == nil {
		new, err := i.loadCache()
		if err != nil {
			return nil, err
		}
		i.publicKey = new.publicKey
	}
	return &base.OptionalString{Value: ByteToHex(i.publicKey)}, nil
}

// 获取地址，该方法会优先读取缓存
func (i *AccountInfo) Address() (*base.OptionalString, error) {
	if i.address == "" {
		new, err := i.loadCache()
		if err != nil {
			return nil, err
		}
		i.address = new.address
	}
	return &base.OptionalString{Value: i.address}, nil
}

func (i *AccountInfo) loadCache() (*AccountInfo, error) {
	if cache := getAccountInfo(i.Wallet.key(), i.cacheKey); cache != nil {
		return cache, nil
	}
	account, err := i.Account()
	if err != nil {
		return nil, err
	}
	i.saveCache(account)
	return i, nil
}

func (i *AccountInfo) saveCache(account base.Account) {
	i.publicKey = account.PublicKey()
	i.address = account.Address()
	saveAccountInfo(i.Wallet.key(), i)
}
