package eth

import (
	"context"
	"math/big"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
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
	timeout := 60 * time.Second
	return &EthChain{
		timeout: timeout,
	}
}

func (e *EthChain) CreateRemote(rpcUrl string) (chain *EthChain, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	rpcClient, err := rpc.DialContext(ctx, rpcUrl)
	if err != nil {
		return
	}

	remoteRpcClient := ethclient.NewClient(rpcClient)
	chainId, err := remoteRpcClient.ChainID(ctx)
	if err != nil {
		return
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
