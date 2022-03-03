package core

var (
	ETHChainName        = "Ethereum"
	ETHMainRpcUrl       = "https://eth-mainnet.coming.chat"
	ETHRenkeyTestRpcUrl = "https://eth-testnet.coming.chat"
)

type ChainType int

const (
	ETHMainName       ChainType = 0
	ETHRenkeyTestName ChainType = 1
	BscMainName       ChainType = 2
	BscTestName       ChainType = 3
	SherpaxMainName   ChainType = 4
	SherpaxTestName   ChainType = 5
)
