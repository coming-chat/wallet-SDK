package wallet

import (
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/btc"
	"github.com/coming-chat/wallet-SDK/core/cosmos"
	"github.com/coming-chat/wallet-SDK/core/eth"
	"github.com/coming-chat/wallet-SDK/core/solana"
)

const (
	// contains ethereum, bsc, chainx_eth, polygon...
	ChainTypeEthereum = "ethereum"
	ChainTypeBitcoin  = "bitcoin"
	ChainTypeCosmos   = "cosmos"
	ChainTypeSolana   = "solana"

	// contains chainx, minix, sherpax, polkadot...
	// ChainTypePolka    = "polka"
	// ChainTypeSignet   = "signet"
	// ChainTypeDoge     = "dogecoin"
	// ChainTypeTerra    = "terra"
	// ChainTypeAptos    = "aptos"
	// ChainTypeSui      = "sui"
	// ChainTypeStarcoin = "starcoin"
)

func ChainTypeFrom(address string) *base.StringArray {
	res := &base.StringArray{}
	for _, ch := range []byte(address) {
		valid := (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
		if !valid {
			return res
		}
	}
	if strings.HasPrefix(address, "0x") || strings.HasPrefix(address, "0X") {
		if eth.IsValidAddress(address) {
			res.Append(ChainTypeEthereum)
		}
		// if aptos.IsValidAddress(address) {
		// 	res.Append(ChainTypeAptos)
		// }
		// if sui.IsValidAddress(address) {
		// 	res.Append(ChainTypeSui)
		// }
		// if starcoin.IsValidAddress(address) {
		// 	res.Append(ChainTypeStarcoin)
		// }
	} else {
		// if polka.IsValidAddress(address) {
		// 	res.Append(ChainTypePolka)
		// }
		if btc.IsValidAddress(address, btc.ChainMainnet) {
			res.Append(ChainTypeBitcoin)
		}
		// if btc.IsValidAddress(address, btc.ChainSignet) {
		// 	res.Append(ChainTypeSignet)
		// }
		// if doge.IsValidAddress(address, doge.ChainMainnet) {
		// 	res.Append(ChainTypeDoge)
		// }
		if cosmos.IsValidAddress(address, cosmos.CosmosPrefix) {
			res.Append(ChainTypeCosmos)
		}
		// if cosmos.IsValidAddress(address, cosmos.TerraPrefix) {
		// 	res.Append(ChainTypeTerra)
		// }
		if solana.IsValidAddress(address) {
			res.Append(ChainTypeSolana)
		}
	}
	return res
}
