package wallet

import (
	"errors"
	"fmt"

	"github.com/coming-chat/wallet-SDK/core/aptos"
	"github.com/coming-chat/wallet-SDK/core/base"
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

type WalletType = base.SDKEnumInt

const (
	WalletTypeMnemonic   WalletType = 1
	WalletTypeKeystore   WalletType = 2
	WalletTypePrivateKey WalletType = 3
	WalletTypeWatch      WalletType = 4
	WalletTypeError      WalletType = 5
)

func isValidWalletType(typ WalletType) bool {
	return typ >= WalletTypeMnemonic && typ <= WalletTypeWatch
}

// 任意一个类可以遵循这个协议, 通过创建一个 `CacheWallet` 来计算钱包公钥、地址
// SDK 会在需要的时候从该对象读取 **助记词/keystore/私钥** 等信息
type SDKWalletSecretInfo interface {
	// 需要提供一个 key, 可以缓存钱包地址公钥信息
	SDKCacheKey() string
	// 如果是一个助记词钱包，返回助记词
	SDKMnemonic() string
	// 如果是一个 keystore 钱包，返回 keystore
	SDKKeystore() string
	// 如果是一个 keystore 钱包，返回密码
	SDKPassword() string
	// 如果是一个私钥钱包，返回私钥
	SDKPrivateKey() string
	// 如果是一个观察者钱包，返回观察地址
	SDKWatchAddress() string
}

// 旧的钱包对象在内存里面会缓存 **助记词**，还有很多链的账号 **私钥**，可能会有用户钱包被盗的风险
// 因此 SDK 里面不能缓存 **助记词** 、**私钥** 、还有 **keystore 密码**
//
// 考虑到每次都导入助记词生成账号，而仅仅是为了获取账号地址或者公钥，可能会影响钱包的性能和体验
// 因此新提供了这个可以缓存 *账号地址* 和 *公钥* 这种不敏感信息的钱包
type CacheWallet struct {
	walletType   WalletType
	watchAddress string

	WalletInfo SDKWalletSecretInfo
}

func NewCacheWallet(info SDKWalletSecretInfo) *CacheWallet {
	return &CacheWallet{WalletInfo: info}
}

// Create a watch wallet
func NewWatchWallet(address string) (*CacheWallet, error) {
	chainType := ChainTypeFrom(address)
	if chainType.Count() == 0 {
		return nil, errors.New("Invalid wallet address")
	}
	return &CacheWallet{
		walletType:   WalletTypeWatch,
		watchAddress: address,
	}, nil
}

func (w *CacheWallet) key() string {
	if w.WalletInfo == nil {
		return ""
	}
	return w.WalletInfo.SDKCacheKey()
}

// 获取钱包类型
// 枚举值见 `WalletType` (Mnemonic / Keystore / PrivateKey / Watch / Error)
func (w *CacheWallet) WalletType() WalletType {
	if typ, ok := w.checkWalletType(); ok {
		return typ
	}
	w.walletType, _ = w.readTypeAndValue()
	saveWallet(w)
	return w.walletType
}

func (w *CacheWallet) checkWalletType() (typ WalletType, ok bool) {
	if isValidWalletType(w.walletType) {
		return w.walletType, true
	}
	if cache := getWallet(w.key()); cache != nil && isValidWalletType(cache.walletType) {
		w.walletType = cache.walletType
		return cache.walletType, true
	}
	return WalletTypeError, false
}

func (w *CacheWallet) readTypeAndValue() (WalletType, string) {
	if w.WalletInfo == nil {
		return WalletTypeError, ""
	}
	if m := w.WalletInfo.SDKMnemonic(); len(m) > 24 {
		return WalletTypeMnemonic, m
	} else if k := w.WalletInfo.SDKKeystore(); len(k) > 0 {
		return WalletTypeKeystore, k
	} else if p := w.WalletInfo.SDKPrivateKey(); len(p) > 0 {
		return WalletTypePrivateKey, p
	} else if a := w.WalletInfo.SDKWatchAddress(); len(a) > 0 {
		return WalletTypeWatch, a
	} else {
		return WalletTypeError, ""
	}
}

func (w *CacheWallet) readValue(typ WalletType) (WalletType, string) {
	if w.WalletInfo == nil {
		return WalletTypeError, ""
	}
	switch typ {
	case WalletTypeMnemonic:
		return typ, w.WalletInfo.SDKMnemonic()
	case WalletTypeKeystore:
		return typ, w.WalletInfo.SDKKeystore()
	case WalletTypePrivateKey:
		return typ, w.WalletInfo.SDKPrivateKey()
	case WalletTypeWatch:
		return typ, w.WalletInfo.SDKWatchAddress()
	default:
		return WalletTypeError, ""
	}
}

func (w *CacheWallet) PolkaAccountInfo(network int) *AccountInfo {
	return &AccountInfo{
		Wallet:   w,
		cacheKey: fmt.Sprintf("polka-%v", network),
		mnemonicCreator: func(val string) (base.Account, error) {
			return polka.NewAccountWithMnemonic(val, network)
		},
		keystoreCreator: func(val string) (base.Account, error) {
			pwd := w.WalletInfo.SDKPassword()
			return polka.NewAccountWithKeystore(val, pwd, network)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return polka.AccountWithPrivateKey(val, network)
		},
	}
}

func (w *CacheWallet) BitcoinAccountInfo(chainnet string, addressType btc.AddressType) *AccountInfo {
	return &AccountInfo{
		Wallet:   w,
		cacheKey: fmt.Sprintf("bitcoin-%v", chainnet),
		mnemonicCreator: func(val string) (base.Account, error) {
			acc, err := btc.NewAccountWithMnemonic(val, chainnet)
			if err != nil {
				return nil, err
			}
			acc.AddressType = addressType
			return acc, nil
		},
		privkeyCreator: func(val string) (base.Account, error) {
			acc, err := btc.AccountWithPrivateKey(val, chainnet)
			if err != nil {
				return nil, err
			}
			acc.AddressType = addressType
			return acc, nil
		},
	}
}

func (w *CacheWallet) EthereumAccountInfo() *AccountInfo {
	return &AccountInfo{
		Wallet:   w,
		cacheKey: "ethereum",
		mnemonicCreator: func(val string) (base.Account, error) {
			return eth.NewAccountWithMnemonic(val)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return eth.AccountWithPrivateKey(val)
		},
	}
}

func (w *CacheWallet) CosmosAccountInfo(cointype int64, prefix string) *AccountInfo {
	return &AccountInfo{
		Wallet:   w,
		cacheKey: fmt.Sprintf("cosmos-%v-%v", cointype, prefix),
		mnemonicCreator: func(val string) (base.Account, error) {
			return cosmos.NewAccountWithMnemonic(val, cointype, prefix)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return cosmos.AccountWithPrivateKey(val, cointype, prefix)
		},
	}
}

func (w *CacheWallet) DogecoinAccountInfo(chainnet string) *AccountInfo {
	return &AccountInfo{
		Wallet:   w,
		cacheKey: fmt.Sprintf("dogecoin-%v", chainnet),
		mnemonicCreator: func(val string) (base.Account, error) {
			return doge.NewAccountWithMnemonic(val, chainnet)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return doge.AccountWithPrivateKey(val, chainnet)
		},
	}
}

func (w *CacheWallet) SolanaAccountInfo() *AccountInfo {
	return &AccountInfo{
		Wallet:   w,
		cacheKey: "solana",
		mnemonicCreator: func(val string) (base.Account, error) {
			return solana.NewAccountWithMnemonic(val)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return solana.AccountWithPrivateKey(val)
		},
	}
}

func (w *CacheWallet) AptosAccountInfo() *AccountInfo {
	return &AccountInfo{
		Wallet:   w,
		cacheKey: "aptos",
		mnemonicCreator: func(val string) (base.Account, error) {
			return aptos.NewAccountWithMnemonic(val)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return aptos.AccountWithPrivateKey(val)
		},
	}
}

func (w *CacheWallet) SuiAccountInfo() *AccountInfo {
	return &AccountInfo{
		Wallet:   w,
		cacheKey: "sui",
		mnemonicCreator: func(val string) (base.Account, error) {
			return sui.NewAccountWithMnemonic(val)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return sui.AccountWithPrivateKey(val)
		},
	}
}

func (w *CacheWallet) StarcoinAccountInfo() *AccountInfo {
	return &AccountInfo{
		Wallet:   w,
		cacheKey: "starcoin",
		mnemonicCreator: func(val string) (base.Account, error) {
			return starcoin.NewAccountWithMnemonic(val)
		},
		privkeyCreator: func(val string) (base.Account, error) {
			return starcoin.AccountWithPrivateKey(val)
		},
	}
}

func (w *CacheWallet) StarknetAccountInfo(isCairo0 bool) *AccountInfo {
	return &AccountInfo{
		Wallet:   w,
		cacheKey: fmt.Sprintf("starknet-%v", isCairo0),
		mnemonicCreator: func(val string) (base.Account, error) {
			acc, err := starknet.NewAccountWithMnemonic(val)
			if err != nil {
				return nil, err
			}
			acc.Cairo0 = isCairo0
			return acc, nil
		},
		privkeyCreator: func(val string) (base.Account, error) {
			acc, err := starknet.AccountWithPrivateKey(val)
			if err != nil {
				return nil, err
			}
			acc.Cairo0 = isCairo0
			return acc, nil
		},
	}
}
