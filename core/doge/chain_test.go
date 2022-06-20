package doge

type chainInfo struct {
	scan string
}

type chainCfg struct {
	mainnet chainInfo
	testnet chainInfo
}

var chains = &chainCfg{
	mainnet: chainInfo{
		scan: "https://dogechain.info/",
	},
	testnet: chainInfo{
		scan: "https://sochain.com/testnet/doge",
	},
}
