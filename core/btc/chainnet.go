package btc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

const (
	ChainSignet  = "signet"
	ChainMainnet = "mainnet"
	// ComingChat used, similar mainnet's alias.
	ChainBitcoin = "bitcoin"
)

func isValidChain(chainnet string) bool {
	switch chainnet {
	case ChainSignet, ChainBitcoin, ChainMainnet:
		return true
	default:
		return false
	}
}

func netParamsOf(chainnet string) (*chaincfg.Params, error) {
	switch chainnet {
	case ChainSignet:
		return &chaincfg.SigNetParams, nil
	case ChainMainnet, ChainBitcoin:
		return &chaincfg.MainNetParams, nil
	}
	return nil, ErrUnsupportedChain
}

func scanHostOf(chainnet string) (string, error) {
	switch chainnet {
	case ChainSignet:
		return "https://electrs.coming.chat/signet", nil
	case ChainMainnet, ChainBitcoin:
		return "https://electrs.coming.chat/mainnet", nil
	default:
		return "", ErrUnsupportedChain
	}
}

func rpcClientOf(chainnet string) (*rpcclient.Client, error) {
	switch chainnet {
	case ChainSignet:
		return rpcclient.New(&rpcclient.ConnConfig{
			Host:         "115.29.163.193:38332",
			User:         "auth",
			Pass:         "bitcoin-b2dd077",
			HTTPPostMode: true,
			DisableTLS:   true,
		}, nil)

	case ChainMainnet, ChainBitcoin:
		return rpcclient.New(&rpcclient.ConnConfig{
			Host:         "115.29.163.193:8332",
			User:         "auth",
			Pass:         "bitcoin-b2dd077",
			HTTPPostMode: true,
			DisableTLS:   true,
		}, nil)
	}

	return nil, ErrUnsupportedChain
}

func nameOf(chainnet string) (string, error) {
	switch chainnet {
	case ChainSignet:
		return "sBTC", nil
	case ChainMainnet, ChainBitcoin:
		return "BTC", nil
	}
	return "", ErrUnsupportedChain
}
