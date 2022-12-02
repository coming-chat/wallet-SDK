package wallet

import (
	"sync"

	"github.com/coming-chat/wallet-SDK/core/base"
)

// 旧的钱包对象在内存里面会缓存 **助记词**，还有很多链的账号 **私钥**，这是比较危险的，可能会有用户钱包被盗的风险
// 因此 SDK 里面不能缓存 **助记词** 、**私钥** 、还有 **keystore 密码**
// 也建议客户端每次使用完带有私钥的账号后，不要缓存这些账号，而是销毁它们
//
// 考虑到每次都导入助记词生成账号，而仅仅是为了获取账号地址或者公钥，可能会影响钱包的性能和体验
// 因此这里会提供一个可以缓存 *账号地址* 和 *公钥* 这种不敏感信息的工具
type AccountCache struct {
	Store sync.Map
}

func NewAccountCache() *AccountCache {
	return &AccountCache{}
}

func (c *AccountCache) GetAccountInfo(walletName, chainName string) *AccountInfo {
	key := walletName + "###" + chainName
	if val, ok := c.Store.Load(key); ok {
		if info, ok := val.(*AccountInfo); ok && info.Chain == chainName {
			return info
		}
	}
	return nil
}

func (c *AccountCache) SaveAccountInfo(walletName, chainName string, info *AccountInfo) {
	key := walletName + "###" + chainName
	if info == nil {
		c.Store.Delete(key)
	} else {
		info.Chain = chainName
		c.Store.Store(key, info)
	}
}

// 缓存账号信息，该账号的私钥不会被缓存
func (c *AccountCache) SaveAccount(walletName, chainName string, account base.Account) {
	key := walletName + "###" + chainName
	if account == nil {
		c.Store.Delete(key)
		return
	}
	info := &AccountInfo{
		PublicKey: account.PublicKey(),
		Address:   account.Address(),
		Chain:     chainName,
	}
	c.Store.Store(key, info)
}

func (c *AccountCache) Get(key string) *AccountInfo {
	if val, ok := c.Store.Load(key); ok {
		if info, ok := val.(*AccountInfo); ok {
			return info
		}
	}
	return nil
}

func (c *AccountCache) Save(key string, info *AccountInfo) {
	if info == nil {
		c.Store.Delete(key)
	} else {
		c.Store.Store(key, info)
	}
}

func (c *AccountCache) Delete(key string) {
	c.Store.Delete(key)
}

func (c *AccountCache) Clean() {
	c.Store = sync.Map{}
}

type AccountInfo struct {
	PublicKey []byte
	Address   string
	Chain     string
}

func (i *AccountInfo) PublicKeyHex() string {
	return ByteToHex(i.PublicKey)
}
