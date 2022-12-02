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

// 旧的钱包对象在内存里面会缓存 **助记词**，还有很多链的账号 **私钥**，可能会有用户钱包被盗的风险
// 因此 SDK 里面不能缓存 **助记词** 、**私钥** 、还有 **keystore 密码**
//
// 考虑到每次都导入助记词生成账号，而仅仅是为了获取账号地址或者公钥，可能会影响钱包的性能和体验
// 因此新提供了这个可以缓存 *账号地址* 和 *公钥* 这种不敏感信息的钱包
//
// 使用方法：
//  1. 先提供一个遵循 `WalletInfoProvider` 协议的对象, sdk会在需要的时候从该对象读取 **助记词/keystore/私钥** 等信息
//     InfoProvider = `your wallet info provider`
//  2. 通过 walletId 创建钱包
//     var	wallet = NewCacheWallet("wallet1")
//  3. 调用相应链的账号方法，获取账号信息
//     var accountInfo = wallet.PolkaAccountInfo(0)
//     var accountInfo = wallet.EthereumAccountInfo()
//     var accountInfo = wallet.SolanaAccountInfo()
//     var accountInfo = ...
//  4. 通过账号信息获取 账号对象/私钥(每次都会取助记词等信息，用完销毁)、公钥/地址(如果有缓存，会直接读取缓存)
//     var account = accountInfo.Account()
//     var privateKey = accountInfo.PrivateKeyHex()
//     var publicKey = accountInfo.PublickKeyHex()
//     var address = accountInfo.Address()
//
// 其他：
//
//	获取钱包类型
//	var walletType = wallet.WalletType()
//	枚举值见 `WalletType` (Mnemonic / Keystore / PrivateKey / Watch / Error)
type CacheWallet struct {
	WalletId   string
	walletType WalletType

	watchAddress string
}

func NewCacheWallet(walletId string) *CacheWallet {
	return &CacheWallet{WalletId: walletId}
}

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

func (w *CacheWallet) WalletType() WalletType {
	if typ, ok := w.checkWalletType(); ok {
		return typ
	}
	w.walletType, _ = readTypeAndValue(w.WalletId)
	SaveWallet(w)
	return w.walletType
}

func (w *CacheWallet) checkWalletType() (typ WalletType, ok bool) {
	if w.walletType >= WalletTypeMnemonic && w.walletType <= WalletTypeWatch {
		return w.walletType, true
	}
	if cache := GetWallet(w.WalletId); cache != nil &&
		cache.walletType >= WalletTypeMnemonic && cache.walletType <= WalletTypeWatch {
		w.walletType = cache.walletType
		return cache.walletType, true
	}
	return WalletTypeError, false
}

func (w *CacheWallet) PolkaAccountInfo(network int) *AccountInfo {
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

func (w *CacheWallet) BitcoinAccountInfo(chainnet string) *AccountInfo {
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

func (w *CacheWallet) EthereumAccountInfo() *AccountInfo {
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

func (w *CacheWallet) CosmosAccountInfo(cointype int64, prefix string) *AccountInfo {
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

func (w *CacheWallet) DogecoinAccountInfo(chainnet string) *AccountInfo {
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

func (w *CacheWallet) SolanaAccountInfo() *AccountInfo {
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

func (w *CacheWallet) AptosAccountInfo() *AccountInfo {
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

func (w *CacheWallet) SuiAccountInfo() *AccountInfo {
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

func (w *CacheWallet) StarcoinAccountInfo() *AccountInfo {
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
