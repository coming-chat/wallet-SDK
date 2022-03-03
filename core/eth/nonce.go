package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

func (e *EthChain) Nonce(spenderAddressHex string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	nonce, err := e.RemoteRpcClient.PendingNonceAt(ctx, common.HexToAddress(spenderAddressHex))
	if err != nil {
		return 0, err
	}

	return int64(nonce), nil
}
