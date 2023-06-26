package starknet

import (
	"errors"
	"math/big"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/dontpanicdao/caigo/gateway"
)

var ErrUnknownNetwork = errors.New("unknown network (known: mainnet, goerli)")

type Network = base.SDKEnumInt

const (
	NetworkMainnet = 0
	NetworkGoerli  = 1
	// NetworkGoerli2 = 2
)

func NetworkString(n Network) (string, error) {
	switch n {
	case NetworkMainnet:
		return "mainnet", nil
	case NetworkGoerli:
		return "goerli", nil
	// case NetworkGoerli2:
	// return "goerli2", nil
	default:
		return "", ErrUnknownNetwork
	}
}

func NetworkChainID(n Network) (*big.Int, error) {
	switch n {
	case NetworkMainnet:
		return big.NewInt(0).SetBytes([]byte(gateway.MAINNET_ID)), nil
	case NetworkGoerli:
		return big.NewInt(0).SetBytes([]byte(gateway.GOERLI_ID)), nil
	// case NetworkGoerli2:
	// return big.NewInt(0).SetBytes([]byte("SN_GOERLI2")), nil
	default:
		return nil, ErrUnknownNetwork
	}
}
