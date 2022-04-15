package eth

import (
	"context"
)

// 获取链ID
func GetChainId(e *EthChain) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	chainId, err := e.RemoteRpcClient.ChainID(ctx)
	if err != nil {
		return "0", MapToBasicError(err)
	}

	return chainId.String(), nil
}
