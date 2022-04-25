package polka

type Chain struct {
	*Util
	RpcUrl  string
	ScanUrl string
}

// @param rpcUrl will be used to get metadata, query balance, estimate fee, send signed tx.
// @param scanUrl will be used to query transaction details
func NewChainWithRpc(rpcUrl, scanUrl string, network int) (*Chain, error) {
	util := NewUtilWithNetwork(network)
	return &Chain{
		Util:    util,
		RpcUrl:  rpcUrl,
		ScanUrl: scanUrl,
	}, nil
}
