package eth

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthChain struct {
	RemoteRpcClient *ethclient.Client
	RpcClient       *rpc.Client
	timeout         time.Duration
	chainId         *big.Int
	rpcUrl          string
}

func NewEthChain() *EthChain {
	timeout := 15 * time.Second
	return &EthChain{
		timeout: timeout,
	}
}

func (e *EthChain) CreateRemote(rpcUrl string) (*EthChain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	rpcClient, err := rpc.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, MapToBasicError(err)
	}

	remoteRpcClient := ethclient.NewClient(rpcClient)
	chainId, err := remoteRpcClient.ChainID(ctx)
	if err != nil {
		return nil, MapToBasicError(err)
	}
	e.chainId = chainId
	e.RpcClient = rpcClient
	e.RemoteRpcClient = remoteRpcClient
	e.rpcUrl = rpcUrl
	return e, nil
}

func (e *EthChain) ConnectRemote(rpcUrl string) error {
	_, err := e.CreateRemote(rpcUrl)
	return err
}

func (e *EthChain) Close() {
	if e.RemoteRpcClient != nil {
		e.RemoteRpcClient.Close()
	}
	if e.RpcClient != nil {
		e.RpcClient.Close()
	}
}
