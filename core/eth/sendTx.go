package eth

import (
	"context"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/common"
)

// 对交易进行广播
func (e *EthChain) SendRawTransaction(txHex string) (string, error) {
	var hash common.Hash
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	err := e.RpcClient.CallContext(ctx, &hash, "eth_sendRawTransaction", txHex)
	if err != nil {
		return "", base.MapAnyToBasicError(err)
	}
	return hash.String(), nil
}
