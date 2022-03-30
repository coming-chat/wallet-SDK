package btc

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

// 检查地址是否有效
// @param address 比特币地址
// @param chainnet 链名称
func IsValidAddress(address, chainnet string) bool {
	var netParams *chaincfg.Params
	switch chainnet {
	case "signet":
		netParams = &chaincfg.SigNetParams
	case "mainnet", "bitcoin": // bitcoin is to fit ComingChat
		netParams = &chaincfg.MainNetParams
	default:
		return false
	}

	_, err := btcutil.DecodeAddress(address, netParams)
	return err == nil
}
