package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

func (e *EthChain) Balance(address string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	result, err := e.RemoteRpcClient.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return "0", err
	}
	return result.String(), nil
}
