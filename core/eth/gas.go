package eth

import (
	"context"
	"strconv"

	"github.com/ethereum/go-ethereum"
)

// 获取标准gas价格
func (e *EthChain) SuggestGasPrice() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasPrice, err := e.RemoteRpcClient.SuggestGasPrice(ctx)

	if err != nil {
		return "0", nil
	}
	return gasPrice.String(), err
}

func (e *EthChain) EstimateGasLimit(msg ethereum.CallMsg) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasCount, err := e.RemoteRpcClient.EstimateGas(ctx, msg)
	if err != nil {
		return "0", err
	}
	gasLimitStr := strconv.FormatUint(gasCount, 10)
	return gasLimitStr, nil
}
