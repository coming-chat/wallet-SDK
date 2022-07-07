package testcase

import "os"

type AccountCase struct {
	Mnemonic string
	Address  string
}

type AccountGroup struct {
	BtcMainnet  AccountCase
	BtcSignet   AccountCase
	Cosmos      AccountCase
	Terra       AccountCase
	DogeMainnet AccountCase
	DogeTestnet AccountCase
	Ethereum    AccountCase
	Polka0      AccountCase
	Polka2      AccountCase
	Polka44     AccountCase
	Solana      AccountCase
}

// 你需要在你的电脑环境中配置助记词的环境变量，并更新对应的币种地址
// You need to configure the environment variable of the mnemonic phrase in your computer environment and update the corresponding currency address
var M1 = os.Getenv("WalletSdkTestM1")
var M2 = os.Getenv("WalletSdkTestM2")
var Mterra = os.Getenv("WalletSdkTestMterra")

var Accounts = AccountGroup{
	BtcMainnet:  AccountCase{M1, "bc1pe7mlvfszt45zvffl3lzctrn6t39c2v6sksu946gfeanp8zm3a3msaw8sys"},
	BtcSignet:   AccountCase{M1, "tb1pe7mlvfszt45zvffl3lzctrn6t39c2v6sksu946gfeanp8zm3a3ms2x3l7l"},
	Cosmos:      AccountCase{M1, "cosmos1r64ug62cytmg28leeu42hp9mzq7p442myf0s3m"},
	Terra:       AccountCase{Mterra, "terra14tt4mzwrfxlgv7ly4xgym79jrxugjtdvgsalpj"},
	DogeMainnet: AccountCase{M1, "DHkgR9b8TGtFe6ZKqhAcmYmBKUeYmJ2uqg"},
	DogeTestnet: AccountCase{M1, "ngok9AL3PFLyX58WsWp51xMUZM2qk1RGQ6"},
	Ethereum:    AccountCase{M1, "0x6334d64D5167F726d8A44f3fbCA66613708E59E7"},
	Polka0:      AccountCase{M1, "12eV7FtPbXBgDG6mX4zwPaJdQKgigVtnYofSpS8mgEQbX623"},
	Polka2:      AccountCase{M1, "EDodEyCN6w8XNuhL8kz9NqUhHyJns9pvgmi3oRNbwba5hxN"},
	Polka44:     AccountCase{M1, "5RHWUui8WKff5quBNVhz1E1Kqfyf6ZbgrC3DtWS23ra3u4vV"},
	Solana:      AccountCase{M1, "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqPY"},
}

var Accounts2 = AccountGroup{
	DogeMainnet: AccountCase{M2, "DLGwnBwHB9FpFJw9apYato1kPXyeNYhq6H"},
	DogeTestnet: AccountCase{M2, "njL1WCgC77iY8HWLceC39Cc3dQMwSbXugR"},
	Solana:      AccountCase{M2, "GDqGCNxkZK3QWcWTnXK3TDuMBV168oUPUqd5spdRN8QW"},
}

var EmptyMnemonic = AccountCase{Mnemonic: ""}
var ErrorMnemonic = AccountCase{Mnemonic: "unaware oxygen allow method allow property predict various slice travel please error"}
