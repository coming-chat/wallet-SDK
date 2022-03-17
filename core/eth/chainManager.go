package eth

var chainConnections = make(map[string]*EthChain)

// 通过 rpcUrl, 获取 eth 的连接对象
func GetConnection(rpcUrl string) (*EthChain, error) {
	chain, ok := chainConnections[rpcUrl]
	if ok {
		return chain, nil
	}

	chain, err := NewEthChain().CreateRemote(rpcUrl)
	if err != nil {
		return nil, err
	}

	chainConnections[rpcUrl] = chain
	return chain, nil
}
