package eth

import (
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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

// Encode erc20 transfer data
func EncodeErc20Transfer(toAddress, amount string) ([]byte, error) {
	if !common.IsHexAddress(toAddress) {
		return nil, errors.New("Invalid receiver address")
	}
	amountInt, valid := big.NewInt(0).SetString(amount, 10)
	if !valid {
		return nil, errors.New("Invalid transfer amount")
	}
	return EncodeContractData(Erc20AbiStr, ERC20_METHOD_TRANSFER, common.HexToAddress(toAddress), amountInt)
}

func EncodeErc20Approve(spender string, amount *big.Int) ([]byte, error) {
	if !common.IsHexAddress(spender) {
		return nil, errors.New("Invalid receiver address")
	}
	return EncodeContractData(Erc20AbiStr, ERC20_METHOD_APPROVE, common.HexToAddress(spender), amount)
}

func EncodeContractData(abiString, method string, params ...interface{}) ([]byte, error) {
	parsedAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}
	return parsedAbi.Pack(method, params...)
}

func DecodeContractParams(abiString string, data []byte) (string, []interface{}, error) {
	if len(data) <= 4 {
		return "", nil, nil
	}
	parsedAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return "", nil, err
	}

	method, err := parsedAbi.MethodById(data[:4])
	if err != nil {
		return "", nil, err
	}

	params, err := method.Inputs.Unpack(data[4:])
	return method.RawName, params, err
}
