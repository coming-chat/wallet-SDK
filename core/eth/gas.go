package eth

import (
	"context"

	"github.com/ethereum/go-ethereum"
)

func (e *EthChain) SuggestGasPrice() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasPrice, err := e.RemoteRpcClient.SuggestGasPrice(ctx)

	if err != nil {
		return 0, nil
	}
	return gasPrice.Int64(), err
}

func (e *EthChain) EstimateGas(msg ethereum.CallMsg) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasCount, err := e.RemoteRpcClient.EstimateGas(ctx, msg)
	if err != nil {
		return 0, err
	}
	return int64(gasCount), nil
}
