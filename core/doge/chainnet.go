package doge

import (
	"errors"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

const (
	ChainTestnet = "testnet"
	ChainMainnet = "mainnet"
)

var (
	ErrUnsupportedChain = errors.New("Unsupported Doge chainnet")

	// https://pkg.go.dev/github.com/renproject/multichain@v0.2.9/chain/dogecoin
	mainnetCfg = chaincfg.Params{
		Name: "mainnet",
		Net:  0xc0c0c0c0,

		PubKeyHashAddrID: 30,
		ScriptHashAddrID: 22,
		PrivateKeyID:     158,

		HDPrivateKeyID: [4]byte{0x02, 0xfa, 0xc3, 0x98},
		HDPublicKeyID:  [4]byte{0x02, 0xfa, 0xca, 0xfd},

		Bech32HRPSegwit: "doge",
	}
	testnetCfg = chaincfg.Params{
		Name: "testnet",
		Net:  0xfcc1b7dc,

		PubKeyHashAddrID: 113,
		ScriptHashAddrID: 196,
		PrivateKeyID:     241,

		HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94},
		HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf},

		Bech32HRPSegwit: "doget",
	}
)

func isValidChain(chainnet string) bool {
	switch chainnet {
	case ChainTestnet, ChainMainnet:
		return true
	default:
		return false
	}
}

func netParamsOf(chainnet string) (*chaincfg.Params, error) {
	switch chainnet {
	case ChainTestnet:
		return &testnetCfg, nil
	case ChainMainnet:
		return &mainnetCfg, nil
	}
	return nil, ErrUnsupportedChain
}

func restUrlOf(chainnet string) (string, error) {
	switch chainnet {
	case ChainMainnet:
		return "https://api.blockcypher.com/v1/doge/main", nil
	case ChainTestnet:
		return "", ErrUnsupportedChain
	default:
		return "", ErrUnsupportedChain
	}
}

func scanHostOf(chainnet string) (string, error) {
	switch chainnet {
	case ChainTestnet:
		return "https://electrs-pre.coming.chat", nil
	case ChainMainnet:
		return "https://electrs-mainnet.coming.chat", nil
	default:
		return "", ErrUnsupportedChain
	}
}

func rpcClientOf(chainnet string) (*rpcclient.Client, error) {
	switch chainnet {
	case ChainTestnet:
		return rpcclient.New(&rpcclient.ConnConfig{
			Host:         "115.29.163.193:38332",
			User:         "auth",
			Pass:         "bitcoin-b2dd077",
			HTTPPostMode: true,
			DisableTLS:   true,
		}, nil)

	case ChainMainnet:
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
	case ChainTestnet:
		return "Doge", nil
	case ChainMainnet:
		return "Doget", nil
	}
	return "", ErrUnsupportedChain
}
