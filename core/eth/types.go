package eth

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type UrlParam struct {
	RpcUrl string
	WsUrl  string
}

type CallMethodOpts struct {
	Nonce                int64
	Value                int64
	GasPrice             int64 // MaxFeePerGas
	GasLimit             int64
	IsPredictError       bool
	MaxPriorityFeePerGas int64
}

type BuildTxResult struct {
	SignedTx *types.Transaction
	TxHex    string
}

type TransactionByHashResult struct {
	SignedTx  *types.Transaction
	From      common.Address
	IsPending bool
}
