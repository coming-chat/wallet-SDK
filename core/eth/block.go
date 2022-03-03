package eth

import (
	"context"
)

func (e *EthChain) LatestBlockNumber() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	number, err := e.RemoteRpcClient.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return int64(number), nil
}
