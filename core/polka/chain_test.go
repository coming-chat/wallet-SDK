package polka

type rpcInfo struct {
	url  string
	scan string
	net  int
}
type rpcConfigs struct {
	chainxProd  rpcInfo
	chainxTest  rpcInfo
	minixProd   rpcInfo
	minixTest   rpcInfo
	sherpaxProd rpcInfo
	sherpaxTest rpcInfo
	polkadot    rpcInfo
	kusama      rpcInfo
}

var rpcs = rpcConfigs{
	chainxProd:  rpcInfo{"https://mainnet.chainx.org/rpc", "https://multiscan-api.coming.chat", 44},
	chainxTest:  rpcInfo{"https://testnet3.chainx.org/rpc", "https://multiscan-api-pre.coming.chat", 44},
	minixProd:   rpcInfo{"https://minichain-mainnet.coming.chat/rpc", "https://multiscan-api.coming.chat", 44},
	minixTest:   rpcInfo{"https://rpc-minichain.coming.chat", "https://multiscan-api-pre.coming.chat", 44},
	sherpaxProd: rpcInfo{"https://mainnet.sherpax.io/rpc", "https://multiscan-api.coming.chat", 44},
	sherpaxTest: rpcInfo{"https://sherpax-testnet.chainx.org/rpc", "https://multiscan-api-pre.coming.chat", 44},
	polkadot:    rpcInfo{"https://polkadot.api.onfinality.io/public-rpc", "", 0},
	kusama:      rpcInfo{"https://kusama.api.onfinality.io/public-rpc", "", 2},
}

func (n *rpcInfo) Chain() (*Chain, error) {
	return NewChainWithRpc(n.url, n.scan, n.net)
}

func (n *rpcInfo) MetadataString() (string, error) {
	c, err := NewChainWithRpc(n.url, n.scan, n.net)
	if err != nil {
		return "", err
	}
	return c.GetMetadataString()
}
