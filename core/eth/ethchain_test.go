package eth

const (
	ethRpcUrl = "https://data-seed-prebsc-1-s1.binance.org:8545"
	// ethRpcUrl           = "https://bsc-dataseed.binance.org"
)

func testEthChain() *EthChain {
	chain := NewEthChain()
	_, _ = chain.CreateRemote(ethRpcUrl)
	return chain
}
