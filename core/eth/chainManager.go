package eth

import (
	"sync"
)

var chainConnections = make(map[string]*EthChain)
var lock sync.RWMutex

// 通过 rpcUrl, 获取 eth 的连接对象
func GetConnection(rpcUrl string) (*EthChain, error) {
	return getConnectionWithTimeout(rpcUrl, 0)
}

// @param timeout time unit millisecond, zero instead use default.
func getConnectionWithTimeout(rpcUrl string, timeout int64) (*EthChain, error) {
	chain, ok := chainConnections[rpcUrl]
	if ok {
		return chain, nil
	}

	// 通过加锁范围
	lock.Lock()
	defer lock.Unlock()

	// 再判断一次
	chain, ok = chainConnections[rpcUrl]
	if ok {
		return chain, nil
	}

	// 创建并存储
	chain, err := NewEthChain().CreateRemoteWithTimeout(rpcUrl, timeout)
	if err != nil {
		return nil, err
	}

	chainConnections[rpcUrl] = chain
	return chain, nil
}
