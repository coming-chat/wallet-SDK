package wallet

import (
	"sync"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func xx() {
	x, _ := gsrpc.NewSubstrateAPI("")
	x.RPC.State.GetMetadataLatest()
}

type polkaclient struct {
	api      *gsrpc.SubstrateAPI
	metadata *types.Metadata
	rpcUrl   string
}

func newPolkaClient(rpcUrl string, metadataString string) (*polkaclient, error) {
	api, err := gsrpc.NewSubstrateAPI(rpcUrl)
	if err != nil {
		return nil, err
	}

	var metadata *types.Metadata
	if metadataString == "" {
		metadata, err = api.RPC.State.GetMetadataLatest()
		if err != nil {
			return nil, err
		}
	} else {
		var meta types.Metadata
		err = types.DecodeFromHexString(metadataString, &meta)
		if err != nil {
			return nil, ErrWrongMetadata
		}
		metadata = &meta
	}

	return &polkaclient{
		api:      api,
		metadata: metadata,
		rpcUrl:   rpcUrl,
	}, nil
}

func (c *polkaclient) ReloadMetadata() error {
	meta, err := c.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}
	c.metadata = meta
	return nil
}

func (c *polkaclient) LoadMetadataIfNotExists() error {
	if c.metadata != nil {
		return nil
	}
	return c.ReloadMetadata()
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

// 通过 rpcUrl, 获取 eth 的连接对象
func getPolkaClient(rpcUrl string) (*polkaclient, error) {
	return getPolkaClientWithMetadata(rpcUrl, "")
}

func getPolkaClientWithMetadata(rpcUrl, metadata string) (*polkaclient, error) {
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
