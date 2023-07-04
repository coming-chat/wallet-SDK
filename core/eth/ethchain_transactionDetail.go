package eth

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// 支持 对象 和 json 字符串 相互转换
type Jsonable interface {
	// 将对象转为 json 字符串
	JsonString() string

	// 接受 json 字符串，生成对象的构造方法
	// 该方法无法统一声明，每一个类需要各自提供
	// NewXxxWithJsonString(s string) *Xxx
}

func decodeSigner(txn *types.Transaction) (common.Address, error) {
	var signer types.Signer
	switch {
	case txn.Type() == types.AccessListTxType:
		signer = types.NewEIP2930Signer(txn.ChainId())
	case txn.Type() == types.DynamicFeeTxType:
		signer = types.NewLondonSigner(txn.ChainId())
	default:
		signer = types.NewEIP155Signer(txn.ChainId())
	}
	return types.Sender(signer, txn)
}

// 获取交易的详情
// @param hashString 交易的 hash
// @return 交易详情 和 交易原文信息
func (e *EthChain) FetchTransactionDetail(hashString string) (detail *base.TransactionDetail, txn *types.Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	tx, isPending, err := e.RemoteRpcClient.TransactionByHash(ctx, common.HexToHash(hashString))
	if err != nil {
		return
	}
	sender, err := decodeSigner(tx)
	if err != nil {
		return
	}

	gasFeeInt := big.NewInt(0).Mul(tx.GasPrice(), big.NewInt(0).SetUint64(tx.Gas()))
	detail = &base.TransactionDetail{
		HashString:   hashString,
		FromAddress:  sender.String(),
		ToAddress:    tx.To().String(),
		Amount:       tx.Value().String(),
		EstimateFees: gasFeeInt.String(),
	}

	if isPending {
		detail.Status = base.TransactionStatusPending
		return
	}

	// 交易receipt 状态信息，0表示失败，1表示成功
	// 当交易没有处于pending状态时，可以查询receipt信息，即交易是否成功, err为nil时，表示查询成功进入if语句赋值
	receipt, err := e.TransactionReceiptByHash(hashString)
	if err != nil {
		return
	}
	blockHeader, err := e.RemoteRpcClient.HeaderByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return
	}

	if receipt.Status == 0 {
		detail.Status = base.TransactionStatusFailure
		// get error message
		_, err := e.RemoteRpcClient.CallContract(ctx, ethereum.CallMsg{
			From:       sender,
			To:         tx.To(),
			Data:       tx.Data(),
			Gas:        tx.Gas(),
			GasPrice:   tx.GasPrice(),
			GasFeeCap:  tx.GasFeeCap(),
			GasTipCap:  tx.GasTipCap(),
			Value:      tx.Value(),
			AccessList: tx.AccessList(),
		}, receipt.BlockNumber)
		if err != nil {
			detail.FailureMessage = err.Error()
			err = nil
		}

	} else {
		detail.Status = base.TransactionStatusSuccess
	}

	effectiveGasPrice := tx.GasPrice()
	if receipt.EffectiveGasPrice != nil {
		effectiveGasPrice = receipt.EffectiveGasPrice
	}
	gasFeeInt = big.NewInt(0).Mul(effectiveGasPrice, big.NewInt(0).SetUint64(receipt.GasUsed))
	if receipt.L1Fee != nil {
		gasFeeInt = gasFeeInt.Add(gasFeeInt, receipt.L1Fee)
	}
	detail.EstimateFees = gasFeeInt.String()
	detail.FinishTimestamp = int64(blockHeader.Time)

	return detail, tx, nil
}

// 获取交易的状态
// @param hashString 交易的 hash
func (e *EthChain) FetchTransactionStatus(hashString string) base.TransactionStatus {
	if len(hashString) == 0 {
		return base.TransactionStatusNone
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	_, isPending, err := e.RemoteRpcClient.TransactionByHash(ctx, common.HexToHash(hashString))
	if err != nil {
		return base.TransactionStatusNone
	}
	if isPending {
		return base.TransactionStatusPending
	}

	// 交易receipt 状态信息，0表示失败，1表示成功
	receipt, err := e.TransactionReceiptByHash(hashString)
	if err != nil {
		return base.TransactionStatusFailure
	}
	if receipt.Status == 0 {
		return base.TransactionStatusFailure
	} else {
		return base.TransactionStatusSuccess
	}
}

// 批量获取交易的转账状态
// @param hashList 要批量查询的交易的 hash 数组
// @return 交易状态数组，它的顺序和 hashList 是保持一致的
func (e *EthChain) BatchTransactionStatus(hashList []string) []string {
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(e.FetchTransactionStatus(s)), nil
	})
	return statuses
}

// SDK 批量获取交易的转账状态，hash 列表和返回值，都只能用字符串，逗号隔开传递
// @param hashListString 要批量查询的交易的 hash，用逗号拼接的字符串："hash1,hash2,hash3"
// @return 批量的交易状态，它的顺序和 hashListString 是保持一致的: "status1,status2,status3"
func (e *EthChain) SdkBatchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses := e.BatchTransactionStatus(hashList)
	return strings.Join(statuses, ",")
}

// 根据交易hash查询交易状态
// TO-DO  返回更详细的信息，解析交易余额，交易动作
func (e *EthChain) TransactionByHash(txHash string) (*TransactionByHashResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	tx, isPending, err := e.RemoteRpcClient.TransactionByHash(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}
	sender, err := decodeSigner(tx)
	if err != nil {
		return nil, err
	}

	// 交易receipt 状态信息，0表示失败，1表示成功
	receipt, err := e.TransactionReceiptByHash(txHash)
	var status, gasUsed, blockNumber string
	// 当交易没有处于pending状态时，可以查询receipt信息，即交易是否成功, err为nil时，表示查询成功进入if语句赋值
	if err == nil {
		gasUsed = strconv.FormatUint(receipt.GasUsed, 10)
		status = strconv.FormatUint(receipt.Status, 10)
		blockNumber = receipt.BlockHash.String()
	}

	return &TransactionByHashResult{
		tx,
		sender,
		isPending,
		status,
		gasUsed,
		blockNumber,
	}, nil
}

// TransactionReceipt 是指交易的收据，每笔交易执行完
// 会产生一个收据，收据中包含交易的状态，交易的gas使用情况，交易执行是否成功的状态码等信息
// 交易收据属性列表：
// gasUsed: 交易执行时使用的gas数量
// bloomFilter：交易信息日志检索
// logInfoList: 交易日志集合
// postTxState: 交易执行后的状态，1 表示成功，0表示失败
func (e *EthChain) TransactionReceiptByHash(txHash string) (*Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	var r *Receipt
	err := e.RpcClient.CallContext(ctx, &r, "eth_getTransactionReceipt", txHash)
	if err == nil {
		if r == nil {
			return nil, ethereum.NotFound
		}
	}
	return r, err
}

func (e *EthChain) WaitConfirm(txHash string, interval time.Duration) *Receipt {
	timer := time.NewTimer(0)
	for range timer.C {
		transRes, err := e.TransactionByHash(txHash)
		if err != nil {
			timer.Reset(interval)
			continue
		}
		if transRes.IsPending {
			timer.Reset(interval)
			continue
		}
		receipt, err := e.TransactionReceiptByHash(txHash)
		if err != nil {
			timer.Reset(interval)
			continue
		}
		timer.Stop()
		return receipt
	}
	return nil
}

// customReceipt is inherit from eth/core types.Receipt, and added some necessary properties
type Receipt struct {
	types.Receipt

	// Optimism Layer2 gas info
	L1Fee *big.Int `json:"l1Fee"`

	// Arbitrum Layer2 gas info
	EffectiveGasPrice *big.Int `json:"effectiveGasPrice"`
}

func (r *Receipt) UnmarshalJSON(data []byte) error {
	err := r.Receipt.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	// If the basic eth receipt unmarshal is successed, The latter should not return an error

	type Layer2Addtional struct {
		L1Fee             *hexutil.Big `json:"l1Fee"`
		EffectiveGasPrice *hexutil.Big `json:"effectiveGasPrice"`
	}
	var layer2Gas Layer2Addtional
	err = json.Unmarshal(data, &layer2Gas)
	if err != nil {
		return nil // should not return an error
	}
	r.L1Fee = (*big.Int)(layer2Gas.L1Fee)
	r.EffectiveGasPrice = (*big.Int)(layer2Gas.EffectiveGasPrice)

	return nil
}
