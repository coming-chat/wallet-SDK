package btc

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func IsValidAddress(address, chainnet string) bool {
	var netParams *chaincfg.Params
	switch chainnet {
	case "signet":
		netParams = &chaincfg.SigNetParams
	case "mainnet", "bitcoin":
		netParams = &chaincfg.MainNetParams
	default:
		return false
	}

	_, err := btcutil.DecodeAddress(address, netParams)
	return err == nil
}
