package eth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

func (e *EthChain) gasFactor() float32 {
	if strings.HasPrefix(e.rpcUrl, "https://mainnet.infura.io/v3") {
		return 1.8
	}
	return 1.3
}

// 获取标准gas价格
func (e *EthChain) SuggestGasPrice() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasPrice, err := e.RemoteRpcClient.SuggestGasPrice(ctx)

	if err != nil {
		return "0", nil
	}
	return gasPrice.String(), err
}

// erc20 代币 Transfer，Approve GasLimit 估计
// var erc20TxParams Erc20TxParams
func (e *EthChain) EstimateContractGasLimit(
	// 用户钱包地址，由私钥可以转地址， util包 PrivateKeyToAddress
	fromAddress,
	contractAddress,
	abiStr,
	methodName string,
	erc20JsonParams string) (string, error) {
	var err error
	defer func() {
		err = MapToBasicError(err)
	}()

	parsedAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return DEFAULT_CONTRACT_GAS_LIMIT, err
	}
	contractAddressObj := common.HexToAddress(contractAddress)

	var erc20TxParams Erc20TxParams
	var input []byte
	// 对交易参数进行格式化
	if err := json.Unmarshal([]byte(erc20JsonParams), &erc20TxParams); err != nil {
		return DEFAULT_CONTRACT_GAS_LIMIT, err
	}

	amountBigInt, _ := new(big.Int).SetString(erc20TxParams.Amount, 10)

	if methodName == ERC20_METHOD_TRANSFER || methodName == ERC20_METHOD_APPROVE {
		// 将string地址类型转化为hex类型
		input, err = parsedAbi.Pack(methodName,
			common.HexToAddress(erc20TxParams.ToAddress),
			amountBigInt)
		if err != nil {
			return DEFAULT_CONTRACT_GAS_LIMIT, err
		}
	} else {
		return DEFAULT_CONTRACT_GAS_LIMIT, fmt.Errorf("unsupported method name: %s", methodName)
	}
	value := big.NewInt(0)

	// 获取标准 gasprice, 如果失败则使用默认值 20000000000
	gasPrice, err := e.SuggestGasPrice()
	if err != nil {
		gasPrice = DEFAULT_ETH_GAS_PRICE
	}
	gasPriceBigInt, _ := new(big.Int).SetString(gasPrice, 10)

	// 如果method为transfer，合约余额不足会导致估算手续费失败掉
	msg := ethereum.CallMsg{
		From:     common.HexToAddress(fromAddress),
		To:       &contractAddressObj,
		GasPrice: gasPriceBigInt,
		Value:    value,
		Data:     input,
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	tempGasLimitUint, err := e.RemoteRpcClient.EstimateGas(ctx, msg)
	if err != nil {
		return DEFAULT_CONTRACT_GAS_LIMIT, err
	}
	gasLimit := uint64(float64(tempGasLimitUint) * float64(e.gasFactor()))
	if gasLimit < 60000 {
		gasLimit = 60000
	}
	gasLimitStr := strconv.FormatUint(gasLimit, 10)
	return gasLimitStr, nil
}

func (e *EthChain) estimateGasLimit(msg ethereum.CallMsg) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasCount, err := e.RemoteRpcClient.EstimateGas(ctx, msg)
	if err != nil {
		return DEFAULT_ETH_GAS_LIMIT, MapToBasicError(err)
	}
	gasLimitStr := strconv.FormatUint(gasCount, 10)
	return gasLimitStr, nil
}

// Estimated gasLimit
// @param fromAddress The address where the transfer originated
// @param receiverAddress The address where the transfer will received
// @param gasPrice Previously acquired or entered by the user
// @param amount The amount transferred
// @return Estimate gasLimit, is a `String` converted from `Uint64`
func (e *EthChain) EstimateGasLimit(fromAddress string, receiverAddress string, gasPrice string, amount string) (string, error) {
	from := common.HexToAddress(fromAddress)
	contractAddressObj := common.HexToAddress(receiverAddress)
	value := big.NewInt(0)

	amountBigInt, _ := new(big.Int).SetString(amount, 10)
	parsedAbi, err := abi.JSON(strings.NewReader(Erc20AbiStr))
	if err != nil {
		return DEFAULT_ETH_GAS_LIMIT, err
	}
	input, err := parsedAbi.Pack(ERC20_METHOD_TRANSFER, common.HexToAddress(receiverAddress), amountBigInt)
	if err != nil {
		return DEFAULT_ETH_GAS_LIMIT, err
	}

	price, isNumber := new(big.Int).SetString(gasPrice, 10)
	if !isNumber {
		return DEFAULT_ETH_GAS_LIMIT, errors.New("gasPrice is invalid")
	}

	msg := ethereum.CallMsg{From: from, To: &contractAddressObj, GasPrice: price, Value: value, Data: input}
	gasLimit, err := e.estimateGasLimit(msg)
	if err != nil {
		return DEFAULT_ETH_GAS_LIMIT, MapToBasicError(err)
	}
	gasLimitDecimal, err := decimal.NewFromString(gasLimit)
	if err != nil {
		return DEFAULT_ETH_GAS_LIMIT, err
	}
	return gasLimitDecimal.Mul(decimal.NewFromFloat32(e.gasFactor())).Round(0).String(), nil
}
