package wallet

import (
	"sync"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type polkaclient struct {
	api      *gsrpc.SubstrateAPI
	metadata *types.Metadata
	rpcUrl   string
}

func newPolkaClient(rpcUrl, metadataString string) (*polkaclient, error) {
	if len(metadataString) > 0 {
		var metadata types.Metadata
		err := types.DecodeFromHexString(metadataString, &metadata)
		if err != nil {
			return nil, ErrWrongMetadata
		}
		return &polkaclient{
			rpcUrl:   rpcUrl,
			metadata: &metadata,
		}, nil
	} else {
		return &polkaclient{
			rpcUrl: rpcUrl,
		}, nil
	}
}

func (c *polkaclient) connectApiIfNeeded() error {
	if c.api == nil {
		api, err := gsrpc.NewSubstrateAPI(c.rpcUrl)
		if err != nil {
			return err
		}
		c.api = api
	}
	return nil
}

func (c *polkaclient) ReloadMetadata() error {
	err := c.connectApiIfNeeded()
	if err != nil {
		return err
	}
	meta, err := c.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}
	c.metadata = meta
	return nil
}

func (c *polkaclient) LoadMetadataIfNotExists() error {
	if c.metadata == nil {
		return c.ReloadMetadata()
	}
	return nil
}

func (c *polkaclient) MetadataString() (string, error) {
	err := c.LoadMetadataIfNotExists()
	if err != nil {
		return "", err
	}
	return types.EncodeToHexString(c.metadata)
}

// MARK: - client manager

var clientConnections = make(map[string]*polkaclient)
var lock sync.RWMutex

// 通过 rpcUrl, 获取 polka 的连接对象
func getConnectedPolkaClient(rpcUrl string) (*polkaclient, error) {
	c, err := getOrCreatePolkaClient(rpcUrl, "")
	if err != nil {
		return nil, err
	}
	err = c.connectApiIfNeeded()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getOrCreatePolkaClient(rpcUrl, metadata string) (*polkaclient, error) {
	// 我们不会考虑覆盖 metadata
	// 假设用户只能从 sdk 里面获取到 metadata
	// 那么无论如何从外部传入的 metadata 肯定不会比 sdk 内部的 metadata 更新
	chain, ok := clientConnections[rpcUrl]
	if ok {
		return chain, nil
	}

	// 通过加锁范围
	lock.Lock()
	defer lock.Unlock()

	// 再判断一次
	chain, ok = clientConnections[rpcUrl]
	if ok {
		return chain, nil
	}

	// 创建并存储
	chain, err := newPolkaClient(rpcUrl, metadata)
	if err != nil {
		return nil, err
	}

	clientConnections[rpcUrl] = chain
	return chain, nil
}
