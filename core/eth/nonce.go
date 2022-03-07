package eth

import (
	"context"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

func (e *EthChain) Nonce(spenderAddressHex string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	nonce, err := e.RemoteRpcClient.PendingNonceAt(ctx, common.HexToAddress(spenderAddressHex))
	if err != nil {
		return "0", err
	}
	return strconv.FormatUint(nonce, 10), nil
}
