package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type UrlParam struct {
	RpcUrl string
	WsUrl  string
}

type CallMethodOpts struct {
	Nonce                int64
	Value                string
	GasPrice             string // MaxFeePerGas
	GasLimit             string
	IsPredictError       bool
	MaxPriorityFeePerGas string
}

type CallMethodOptsBigInt struct {
	Nonce                uint64
	Value                *big.Int
	GasPrice             *big.Int // MaxFeePerGas
	GasLimit             uint64
	IsPredictError       bool
	MaxPriorityFeePerGas *big.Int
}

type BuildTxResult struct {
	SignedTx *types.Transaction
	TxHex    string
}

type TransactionByHashResult struct {
	SignedTx    *types.Transaction
	From        common.Address
	IsPending   bool   // 交易是否处于Pending状态
	Status      string // 0: 交易失败, 1: 交易成功
	GasUsed     string // 实际花费gas
	BlockNumber string // 区块高度
}

type Erc20TxParams struct {
	ToAddress string `json:"toAddress"`
	Amount    string `json:"amount"`
	Method    string `json:"method"`
}

type SDKEnumInt = int8
type SDKEnumString = string
