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
}

func NewEthChain() *EthChain {
	timeout := 60 * time.Second
	return &EthChain{
		timeout: timeout,
	}
}

func (e *EthChain) CreateRemote(rpcUrl string) (*EthChain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	rpcClient, err := rpc.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, err
	}

	remoteRpcClient := ethclient.NewClient(rpcClient)
	ctx, cancel = context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	chainId, err := remoteRpcClient.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	e.chainId = chainId
	e.RpcClient = rpcClient
	e.RemoteRpcClient = remoteRpcClient
	return e, nil
}

func (e *EthChain) Close() {
	if e.RemoteRpcClient != nil {
		e.RemoteRpcClient.Close()
	}
	if e.RpcClient != nil {
		e.RpcClient.Close()
	}
}
