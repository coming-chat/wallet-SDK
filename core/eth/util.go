package eth

import (
	"math/big"
	"strconv"
)

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
