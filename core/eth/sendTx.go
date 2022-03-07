package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

func (e *EthChain) SendRawTransaction(txHex string) (string, error) {
	var hash common.Hash
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	err := e.RpcClient.CallContext(ctx, &hash, "eth_sendRawTransaction", txHex)
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}
