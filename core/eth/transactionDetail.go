package eth

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
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

type TransactionStatus = SDKEnumInt

const (
	TransactionStatusNone    TransactionStatus = 0
	TransactionStatusPending TransactionStatus = 1
	TransactionStatusSuccess TransactionStatus = 2
	TransactionStatusFailure TransactionStatus = 3
)

// 可以从链上获取的转账详情信息
// 客户端的详情展示还需要 FromCID, ToCID, CreateTimestamp, Transfer(转出/收入), CoinType, Decimal
// 这些信息需要客户端自己维护
type TransactionDetail struct {
	// 交易在链上的哈希
	HashString string
	// 交易额
	Amount string
	// 交易手续费, Pending 时为预估手续费，交易结束时为真实手续费
	EstimateFees string
	// 转账人的地址
	FromAddress string
	// 收款人的地址
	ToAddress string
	// 交易状态 枚举常量
	// 0: TransactionStatusNone;
	// 1: TransactionStatusPending;
	// 2: TransactionStatusSuccess;
	// 3: TransactionStatusFailure;
	Status TransactionStatus
	// 交易完成时间, 如果在 Pending 中，为 0
	FinishTimestamp int64
	// 失败描述
	FailureMessage string
}

func (i *TransactionDetail) JsonString() string {
	json, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(json)
}

func NewTransactionDetailWithJsonString(s string) *TransactionDetail {
	var i TransactionDetail
	json.Unmarshal([]byte(s), &i)
	return &i
}

// 获取交易的详情
// @param hashString 交易的 hash
// @return 详情对象，该对象无法提供 CID 信息
func (e *EthChain) FetchTransactionDetail(hashString string) (*TransactionDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	tx, isPending, err := e.RemoteRpcClient.TransactionByHash(ctx, common.HexToHash(hashString))
	if err != nil {
		return nil, err
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(e.chainId), nil)
	if err != nil {
		return nil, err
	}

	gasPrice := msg.GasPrice().Uint64()
	estimateGasLimit := msg.Gas()
	detail := &TransactionDetail{
		HashString:   hashString,
		FromAddress:  msg.From().String(),
		ToAddress:    msg.To().String(),
		Amount:       strconv.FormatUint(msg.Value().Uint64(), 10),
		EstimateFees: strconv.FormatUint(gasPrice*estimateGasLimit, 10),
	}

	if isPending {
		detail.Status = TransactionStatusPending
		return detail, nil
	}

	// 交易receipt 状态信息，0表示失败，1表示成功
	// 当交易没有处于pending状态时，可以查询receipt信息，即交易是否成功, err为nil时，表示查询成功进入if语句赋值
	receipt, err := e.TransactionReceiptByHash(hashString)
	if err != nil {
		return nil, err
	}
	blockHeader, err := e.RemoteRpcClient.HeaderByNumber(ctx, receipt.BlockNumber)
	if err != nil {
		return nil, err
	}

	if receipt.Status == 0 {
		detail.Status = TransactionStatusFailure
		// get error message
		_, err := e.RemoteRpcClient.CallContract(ctx, ethereum.CallMsg{
			From: msg.From(),
			To:   msg.To(),
			Data: msg.Data(),
		}, receipt.BlockNumber)
		detail.FailureMessage = err.Error()

	} else {
		detail.Status = TransactionStatusSuccess
	}
	gasUsed := receipt.GasUsed
	detail.EstimateFees = strconv.FormatUint(gasPrice*gasUsed, 10)
	detail.FinishTimestamp = int64(blockHeader.Time)

	return detail, nil
}

// 获取交易的状态
// @param hashString 交易的 hash
func (e *EthChain) FetchTransactionStatus(hashString string) TransactionStatus {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	_, isPending, err := e.RemoteRpcClient.TransactionByHash(ctx, common.HexToHash(hashString))
	if err != nil {
		return TransactionStatusNone
	}
	if isPending {
		return TransactionStatusPending
	}

	// 交易receipt 状态信息，0表示失败，1表示成功
	receipt, err := e.TransactionReceiptByHash(hashString)
	if receipt.Status == 0 {
		return TransactionStatusFailure
	} else {
		return TransactionStatusSuccess
	}
}
