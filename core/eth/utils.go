package eth

import (
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
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

// 私钥转地址
func PrivateKeyToAddress(privateKey string) (string, error) {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}
	return crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex(), nil
}
