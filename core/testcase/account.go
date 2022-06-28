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
}

// 你需要在你的电脑环境中配置助记词的环境变量，并更新对应的币种地址
// You need to configure the environment variable of the mnemonic phrase in your computer environment and update the corresponding currency address
var m1 = os.Getenv("WalletSdkTestM1")
var m2 = os.Getenv("WalletSdkTestM2")
var mterra = os.Getenv("WalletSdkTestMterra")

var Accounts = AccountGroup{
	BtcMainnet:  AccountCase{m1, "bc1pe7mlvfszt45zvffl3lzctrn6t39c2v6sksu946gfeanp8zm3a3msaw8sys"},
	BtcSignet:   AccountCase{m1, "tb1pe7mlvfszt45zvffl3lzctrn6t39c2v6sksu946gfeanp8zm3a3ms2x3l7l"},
	Cosmos:      AccountCase{m1, "cosmos1r64ug62cytmg28leeu42hp9mzq7p442myf0s3m"},
	Terra:       AccountCase{mterra, "terra14tt4mzwrfxlgv7ly4xgym79jrxugjtdvgsalpj"},
	DogeMainnet: AccountCase{m1, "DHkgR9b8TGtFe6ZKqhAcmYmBKUeYmJ2uqg"},
	DogeTestnet: AccountCase{m1, "ngok9AL3PFLyX58WsWp51xMUZM2qk1RGQ6"},
	Ethereum:    AccountCase{m1, "0x6334d64D5167F726d8A44f3fbCA66613708E59E7"},
	Polka0:      AccountCase{m1, "12eV7FtPbXBgDG6mX4zwPaJdQKgigVtnYofSpS8mgEQbX623"},
	Polka2:      AccountCase{m1, "EDodEyCN6w8XNuhL8kz9NqUhHyJns9pvgmi3oRNbwba5hxN"},
	Polka44:     AccountCase{m1, "5RHWUui8WKff5quBNVhz1E1Kqfyf6ZbgrC3DtWS23ra3u4vV"},
}

var Accounts2 = AccountGroup{
	DogeMainnet: AccountCase{m2, "DLGwnBwHB9FpFJw9apYato1kPXyeNYhq6H"},
	DogeTestnet: AccountCase{m2, "njL1WCgC77iY8HWLceC39Cc3dQMwSbXugR"},
}
