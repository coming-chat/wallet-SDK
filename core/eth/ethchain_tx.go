package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// buildTx 创建交易
func (e *EthChain) buildTx(
	privateKeyECDSA *ecdsa.PrivateKey,
	nonce uint64,
	toAddressObj common.Address,
	value *big.Int,
	gasLimit uint64,
	data []byte,
	opts *CallMethodOpts) (*BuildTxResult, error) {
	var rawTx *types.Transaction

	optsBigInt := OptsTobigInt(opts)

	if optsBigInt.MaxPriorityFeePerGas == nil {
		var gasPrice *big.Int = nil
		if opts != nil && optsBigInt.GasPrice != nil {
			gasPrice = optsBigInt.GasPrice
		}
		if gasPrice == nil {
			ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
			defer cancel()
			_gasPrice, err := e.RemoteRpcClient.SuggestGasPrice(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to suggest gas price: %v", err)
			}
			gasPrice = _gasPrice
		}
		rawTx = types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &toAddressObj,
			Value:    value,
			Gas:      gasLimit,
			GasPrice: gasPrice,
			Data:     data,
		})
	} else {
		rawTx = types.NewTx(&types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &toAddressObj,
			Value:     value,
			Gas:       gasLimit,
			GasFeeCap: optsBigInt.GasPrice,             // maxFeePerGas 最大的 gasPrice（包含 baseFee），减去 baseFee 就是小费。gasPrice = min(maxFeePerGas, baseFee + maxPriorityFeePerGas)
			GasTipCap: optsBigInt.MaxPriorityFeePerGas, // maxPriorityFeePerGas，也就是最大的小费。GasTipCap 和 gasFeeCap - baseFee 的更小值才是真正的给矿工的，baseFee 是销毁的。
			Data:      data,
		})
	}
	return e.buildTxWithTransaction(rawTx, privateKeyECDSA)
}

func (e *EthChain) buildTxWithTransaction(transaction *types.Transaction, privateKeyCDSA *ecdsa.PrivateKey) (*BuildTxResult, error) {
	signedTx, err := types.SignTx(transaction, types.LatestSignerForChainID(e.chainId), privateKeyCDSA)
	if err != nil {
		return nil, err
	}
	txBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &BuildTxResult{
		SignedTx: signedTx,
		TxHex:    hexutil.Encode(txBytes),
	}, nil
}

// 创建ETH转账交易
func (e *EthChain) BuildTransferTx(privateKey, toAddress string, opts *CallMethodOpts) (*BuildTxResult, error) {

	privateKey = strings.TrimPrefix(privateKey, "0x")

	toAddressObj := common.HexToAddress(toAddress)
	privateKeyBuf, err := hex.DecodeString(privateKey)

	optsBigInt := OptsTobigInt(opts)
	if err != nil {
		return nil, err
	}

	var value = big.NewInt(0)

	var gasLimit = uint64(0)
	var nonce uint64 = 0
	if opts != nil {
		value = optsBigInt.Value
		gasLimit = optsBigInt.GasLimit
		nonce = optsBigInt.Nonce
	}
	if gasLimit == 0 {
		gasLimit = uint64(21000)
	}

	privateKeyECDSA, err := crypto.ToECDSA(privateKeyBuf)
	if err != nil {
		return nil, err
	}
	publicKeyECDSA := privateKeyECDSA.PublicKey
	fromAddress := crypto.PubkeyToAddress(publicKeyECDSA)
	if nonce == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
		defer cancel()
		nonce, err = e.RemoteRpcClient.PendingNonceAt(ctx, fromAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
		}
	}

	return e.buildTx(privateKeyECDSA, nonce, toAddressObj, value, gasLimit, nil, opts)
}

// 对合约进行调用
func (e *EthChain) BuildCallMethodTx(
	privateKey, contractAddress, abiStr, methodName string,
	opts *CallMethodOpts,
	erc20JsonParams string) (*BuildTxResult, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")

	parsedAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, err
	}
	contractAddressObj := common.HexToAddress(contractAddress)
	privateKeyBuf, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	optsBigInt := OptsTobigInt(opts)

	var value = big.NewInt(0)
	var gasLimit uint64 = 0
	var nonce uint64 = 0
	var isPredictError = true
	if opts != nil {
		value = optsBigInt.Value
		gasLimit = optsBigInt.GasLimit
		nonce = optsBigInt.Nonce
		isPredictError = opts.IsPredictError
	}

	privateKeyECDSA, err := crypto.ToECDSA(privateKeyBuf)
	if err != nil {
		return nil, err
	}
	publicKeyECDSA := privateKeyECDSA.PublicKey
	fromAddress := crypto.PubkeyToAddress(publicKeyECDSA)
	if nonce == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
		defer cancel()
		nonce, err = e.RemoteRpcClient.PendingNonceAt(ctx, fromAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
		}
	}

	var erc20TxParams Erc20TxParams
	var input []byte
	// 对交易参数进行格式化
	if err := json.Unmarshal([]byte(erc20JsonParams), &erc20TxParams); err != nil {
		return nil, err
	}

	amountBigInt, _ := new(big.Int).SetString(erc20TxParams.Amount, 10)

	if methodName == ERC20_METHOD_TRANSFER || methodName == ERC20_METHOD_APPROVE {
		// 将string地址类型转化为hex类型
		input, err = parsedAbi.Pack(methodName,
			common.HexToAddress(erc20TxParams.ToAddress),
			amountBigInt)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unsupported method name: %s", methodName)
	}
	if gasLimit == 0 || isPredictError {
		msg := ethereum.CallMsg{From: fromAddress, To: &contractAddressObj, GasPrice: new(big.Int).SetInt64(10), Value: value, Data: input}
		tempGasLimit, err := e.estimateGasLimit(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas: %v", err)
		}
		tempGasLimitUint, err := strconv.ParseUint(tempGasLimit, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed parse estimate gas: %v", err)
		}
		if gasLimit == 0 {
			gasLimit = uint64(float64(tempGasLimitUint) * 1.3)
		}
	}

	return e.buildTx(privateKeyECDSA, nonce, contractAddressObj, value, gasLimit, input, opts)
}

func (e *EthChain) BuildCallMethodTxWithPayload(
	privateKey, contractAddress, payload string, opts *CallMethodOpts) (*BuildTxResult, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyBuf, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	payload = strings.TrimPrefix(payload, "0x")

	payloadBuf, err := hex.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	contractAddressObj := common.HexToAddress(contractAddress)

	optsBigInt := OptsTobigInt(opts)

	var value = big.NewInt(0)
	var gasLimit uint64 = 0
	var nonce uint64 = 0
	var isPredictError = true
	if opts != nil {
		value = optsBigInt.Value
		gasLimit = optsBigInt.GasLimit
		nonce = optsBigInt.Nonce
		isPredictError = opts.IsPredictError
	}

	privateKeyECDSA, err := crypto.ToECDSA(privateKeyBuf)
	if err != nil {
		return nil, err
	}
	publicKeyECDSA := privateKeyECDSA.PublicKey
	fromAddress := crypto.PubkeyToAddress(publicKeyECDSA)
	if nonce == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
		defer cancel()
		nonce, err = e.RemoteRpcClient.PendingNonceAt(ctx, fromAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
		}
	}

	if gasLimit == 0 || isPredictError {
		msg := ethereum.CallMsg{From: fromAddress, To: &contractAddressObj, GasPrice: new(big.Int).SetInt64(10), Value: value, Data: payloadBuf}
		tempGasLimit, err := e.estimateGasLimit(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas: %v", err)
		}
		tempGasLimitUint, err := strconv.ParseUint(tempGasLimit, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed parse estimate gas: %v", err)
		}
		if gasLimit == 0 {
			gasLimit = uint64(float64(tempGasLimitUint) * 1.3)
		}
	}

	return e.buildTx(privateKeyECDSA, nonce, contractAddressObj, value, gasLimit, payloadBuf, opts)
}
