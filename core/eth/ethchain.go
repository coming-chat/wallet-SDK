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
	timeout         time.Duration
	RemoteRpcClient *ethclient.Client
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
	return e.CreateRemoteWithTimeout(rpcUrl, 0)
}

// @param timeout time unit millisecond. 0 means use chain's default: 60000ms.
func (e *EthChain) CreateRemoteWithTimeout(rpcUrl string, timeout int64) (chain *EthChain, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	var t time.Duration
	if timeout == 0 {
		t = e.timeout
	} else {
		t = time.Duration(timeout * int64(time.Millisecond))
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
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
}
