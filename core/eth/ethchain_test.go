package eth

const (
	ethRpcUrl = "https://data-seed-prebsc-1-s1.binance.org:8545"
	// ethRpcUrl           = "https://bsc-dataseed.binance.org"
)

var ethChain, _ = NewEthChain().CreateRemote(ethRpcUrl)
