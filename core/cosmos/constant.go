package cosmos

const (
	// All comments are made for sdk

	// 0.01
	GasPriceLow = "0.01"
	// 0.025
	GasPriceAverage = "0.025"
	// 0.04
	GasPriceHigh = "0.04"

	// 100000
	GasLimitDefault = "100000"

	// 118
	CosmosCointype = 118
	// cosmos
	CosmosPrefix = "cosmos"
	// uatom
	CosmosAtomDenom = "uatom"

	// 330
	TerraCointype = 330
	// terra
	TerraPrefix = "terra"
	// uluna
	TerraLunaDenom = "uluna"
	// uusd
	TerraUSTDenom = "uusd"

	// 10
	TerraGasPrice = "10"
	// 0.25
	TerraGasPriceUST = "0.25"
	// 80000
	TerraGasLimitDefault = "80000"
)

type GradedGasPrice struct {
	Low     string
	Average string
	High    string
}

type KnownTokenInfo struct {
	Cointype int64
	Prefix   string
	Denom    string
	GasPrice *GradedGasPrice
	GasLimit string
}

var (
	CosmosAtom = &KnownTokenInfo{
		Cointype: 118,
		Prefix:   "cosmos",
		Denom:    "uatom",
		GasPrice: &GradedGasPrice{
			Low:     "0.01",
			Average: "0.025",
			High:    "0.04",
		},
		GasLimit: "100000",
	}
	TerraLunc = &KnownTokenInfo{
		Cointype: 330,
		Prefix:   "terra",
		Denom:    "uluna",
		GasPrice: &GradedGasPrice{"10", "10", "10"},
		GasLimit: "80000",
	}
	TerraUst = &KnownTokenInfo{
		Cointype: 330,
		Prefix:   "terra",
		Denom:    "uusd",
		GasPrice: &GradedGasPrice{"0.25", "0.25", "0.25"},
		GasLimit: "80000",
	}
	// TerraLuna = &KnownTokenInfo{
	// 	Cointype: 330,
	// 	Prefix:   "terra",
	// 	Denom:    "uluna",
	// 	GasPrice: &GradedGasPrice{"10", "10", "10"},
	// 	GasLimit: "80000",
	// }
)
