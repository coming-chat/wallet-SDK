package eth

import (
	"math/big"
	"strconv"
)

// 将MethodOpts 进行转化，由于端的限制，只能传入string字符
func OptsTobigInt(opts *CallMethodOpts) *CallMethodOptsBigInt {

	GasPrice, _ := new(big.Int).SetString(opts.GasPrice, 10)
	GasLimit, _ := strconv.Atoi(opts.GasLimit)
	MaxPriorityFeePerGas, _ := new(big.Int).SetString(opts.MaxPriorityFeePerGas, 10)
	Value, _ := new(big.Int).SetString(opts.Value, 10)

	return &CallMethodOptsBigInt{
		Nonce:                uint64(opts.Nonce),
		Value:                Value,
		GasPrice:             GasPrice,
		MaxPriorityFeePerGas: MaxPriorityFeePerGas,
		GasLimit:             uint64(GasLimit),
	}

}
