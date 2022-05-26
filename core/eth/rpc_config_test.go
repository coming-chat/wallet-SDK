package eth

// Test Basic rpc config

type contracts struct {
	USDT string
	BUSD string
	USDC string
	USB  string
}

type rpcInfo struct {
	url       string
	scan      string
	contracts contracts
}
type rpcConfig struct {
	ethereumProd rpcInfo
	rinkeby      rpcInfo
	binanceTest  rpcInfo
	binanceProd  rpcInfo
	sherpaxTest  rpcInfo
	sherpaxProd  rpcInfo
}

var rpcs = rpcConfig{
	ethereumProd: rpcInfo{
		"https://geth-mainnet.coming.chat",
		"https://etherscan.io",
		contracts{
			USDT: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		}},
	rinkeby: rpcInfo{
		"https://rinkeby.infura.io/v3/161645ea57d2494d996c4d2de2489419",
		"",
		contracts{},
	},
	binanceTest: rpcInfo{
		"https://data-seed-prebsc-1-s1.binance.org:8545",
		"https://testnet.bscscan.com",
		contracts{
			USDT: "0x6cd2Bf22B3CeaDfF6B8C226487265d81164396C5",
			BUSD: "0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee",
			USDC: "0x0644014472cD39f51f57ce91be871537D7A5A2Ab",
		}},
	binanceProd: rpcInfo{
		"https://bsc-dataseed.binance.org",
		"https://bscscan.com",
		contracts{
			USDT: "0x55d398326f99059fF775485246999027B3197955",
			BUSD: "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56",
			USDC: "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d",
		}},
	sherpaxTest: rpcInfo{
		"https://sherpax-testnet.chainx.org/rpc",
		"https://evm-pre.sherpax.io",
		contracts{
			USDT: "0x1635583ACf7beF762E8119887b2f3B9F9BcD1742",
			BUSD: "0x77eD6a802aB1d60A86F2e3c45B43a0Cd7Ee2572B",
			USDC: "0xa017362eB5B22302e4E5c55786f651214BD168A2",
		}},
	sherpaxProd: rpcInfo{
		"https://mainnet.sherpax.io/rpc",
		"https://evm.sherpax.io",
		contracts{
			USB:  "0xE7e312dfC08e060cda1AF38C234AEAcc7A982143",
			USDT: "0x4B53739D798EF0BEa5607c254336b40a93c75b52",
			BUSD: "0x37088186089c7D6BcD556d9A15087DFaE3Ba0C32",
			USDC: "0x935CC842f220CF3A7D10DA1c99F01B1A6894F7C5",
		}},
}

func (n *rpcInfo) Chain() *Chain {
	return NewChainWithRpc(n.url)
}
