package btc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

const (
	ChainSignet  = "signet"
	ChainMainnet = "mainnet"
	ChainTestnet = "testnet"
	// ComingChat used, similar mainnet's alias.
	ChainBitcoin = "bitcoin"
)

func isValidChain(chainnet string) bool {
	switch chainnet {
	case ChainSignet, ChainBitcoin, ChainMainnet, ChainTestnet:
		return true
	default:
		return false
	}
}

func netParamsOf(chainnet string) (*chaincfg.Params, error) {
	switch chainnet {
	case ChainSignet:
		return &chaincfg.SigNetParams, nil
	case ChainTestnet, "testnet3":
		return &chaincfg.TestNet3Params, nil
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
	case ChainTestnet:
		return "https://electrs.coming.chat/testnet", nil
	default:
		return "", ErrUnsupportedChain
	}
}

func rpcClientOf(chainnet string) (*rpcclient.Client, error) {
	switch chainnet {
	case ChainSignet:
		return rpcclient.New(&rpcclient.ConnConfig{
			Host:         "bitcoin.coming.chat/signet/",
			User:         "auth",
			Pass:         "bitcoin-b2dd077",
			HTTPPostMode: true,
			DisableTLS:   true,
		}, nil)

	case ChainMainnet, ChainBitcoin:
		return rpcclient.New(&rpcclient.ConnConfig{
			Host:         "bitcoin.coming.chat/mainnet/",
			User:         "auth",
			Pass:         "bitcoin-b2dd077",
			HTTPPostMode: true,
			DisableTLS:   true,
		}, nil)
	case ChainTestnet:
		return rpcclient.New(&rpcclient.ConnConfig{
			Host:         "bitcoin.coming.chat/testnet/",
			User:         "auth",
			Pass:         "bitcoin-b2dd077",
			HTTPPostMode: true,
			DisableTLS:   true,
		}, nil)
	}

	return nil, ErrUnsupportedChain
}

func zeroWalletHost(chainnet string) (string, error) {
	switch chainnet {
	case ChainMainnet, ChainBitcoin:
		return "https://coming-zero-wallet.coming.chat", nil
	case ChainTestnet:
		return "https://coming-zero-wallet-pre.coming.chat", nil
	}
	return "", ErrUnsupportedChain
}

func comingOrdHost(chainnet string) (string, error) {
	switch chainnet {
	case ChainMainnet, ChainBitcoin:
		return "https://ord.bevm.io/mainnet", nil
	case ChainTestnet:
		return "https://ord.bevm.io/testnet", nil
	}
	return "", ErrUnsupportedChain
}

func unisatHost(chainnet string) (string, error) {
	switch chainnet {
	case ChainMainnet, ChainBitcoin:
		return "https://api.unisat.io", nil
	case ChainTestnet:
		return "https://api-testnet.unisat.io", nil
	}
	return "", ErrUnsupportedChain
}

func nameOf(chainnet string) (string, error) {
	switch chainnet {
	case ChainSignet:
		return "sBTC", nil
	case ChainMainnet, ChainBitcoin:
		return "BTC", nil
	case ChainTestnet:
		return "BTC", nil
	}
	return "", ErrUnsupportedChain
}
