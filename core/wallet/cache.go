package wallet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/coming-chat/wallet-SDK/core/base"
)

var walletCache = sync.Map{}

func SaveWallet(wallet *Wallet) {
	key := "wallet-" + wallet.WalletId
	walletCache.Store(key, wallet)
}

func GetWallet(walletId string) *Wallet {
	key := "wallet-" + walletId
	if c, ok := walletCache.Load(key); ok {
		return c.(*Wallet)
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
	wallet          *Wallet
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
	default:
		return nil, ErrWalletInfoNotExist
	}
	if err != nil {
		return nil, err
	}
	i.loadAndSaveCacheIfNotFound(account)
	return account, nil
}

func (i *AccountInfo) PrivateKeyHex() (string, error) {
	account, err := i.Account()
	if err != nil {
		return "", err
	}
	return account.PrivateKeyHex()
}

func (i *AccountInfo) PublicKeyHex() string {
	if i.publicKey == nil {
		new, err := i.loadAndSaveCacheIfNotFound(nil)
		if err != nil {
			return ""
		}
		i.publicKey = new.publicKey
	}
	return ByteToHex(i.publicKey)
}

func (i *AccountInfo) Address() string {
	if i.address == "" {
		new, err := i.loadAndSaveCacheIfNotFound(nil)
		if err != nil {
			return ""
		}
		i.address = new.address
	}
	return i.address
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
