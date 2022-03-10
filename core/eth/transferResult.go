package eth

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

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

type coin struct {
	Type    string // 与后台交互的类型：ETH_USDT、BSC_USDT
	Symbol  string // 货币符号：BTC、ETH、USDT ...
	Decimal int16  // 货币精度
}

type TransferDetail struct {
	HashString string // 交易在链上的哈希
	CoinType   string // 货币符号：ETH、BSC ...
	Decimal    int16  // 货币精度

	Transfer     string // 收支 income: 收入; expense: 支出;
	Amount       string // 交易额
	EstimateFees string // 交易手续费

	From        string // 转账人的 CID, 客户端维护的
	FromAddress string // 转账人的地址
	To          string // 收款人的 CID, 客户端维护的
	ToAddress   string // 收款人的地址

	Status          int    // 交易状态 0: None; 1: Pending; 2: Success; 3: Failure;
	CreateTimestamp int64  // 交易创建时间, 客户端维护的
	FinishTimestamp int64  // 交易完成时间, 如果在 Pending 中，为 0
	FailureMessage  string // 失败描述
}

func (i *TransferDetail) JsonString() string {
	json, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(json)
}

func NewTransferDetailWithJsonString(s string) *TransferDetail {
	var i TransferDetail
	json.Unmarshal([]byte(s), &i)
	return &i
}

// 创建一个交易信息, 需要传入客户端维护的信息
func NewTransferDetailWithCoinType(coinType string, decimal int16, fromCID string, toCID string) *TransferDetail {
	return &TransferDetail{
		CoinType:        coinType,
		Decimal:         decimal,
		From:            fromCID,
		To:              toCID,
		CreateTimestamp: time.Now().Unix(),
	}
}

// 从本地存储的转账信息里面，更新无法从链上获取到的交易信息
func (d *TransferDetail) UpdateLocalInfoFrom(old *TransferDetail) {
	d.CoinType = old.CoinType
	d.Decimal = old.Decimal
	d.From = old.From
	d.To = old.To
	d.CreateTimestamp = old.CreateTimestamp
}

// 获取交易的详情
// @param hashString 交易的 hash
// @return 详情对象，该对象无法提供 CID 信息
func (e *EthChain) FetchTransferDetail(hashString string) (*TransferDetail, error) {
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

	statusInt := 1 // Pending
	gasPrice := msg.GasPrice().Uint64()
	estimateGasLimit := msg.Gas()
	estimateFees := strconv.FormatUint(gasPrice*estimateGasLimit, 10)
	amount := strconv.FormatUint(msg.Value().Uint64(), 10)
	fromAddress := msg.From().String()
	toAddress := msg.To().String()
	if isPending {
		return &TransferDetail{
			HashString: hashString,
			Status:     statusInt,

			From:         fromAddress,
			To:           toAddress,
			Amount:       amount,
			EstimateFees: estimateFees,
		}, nil
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
		statusInt = 3
	} else {
		statusInt = 2
	}
	gasUsed := receipt.GasUsed
	estimateFees = strconv.FormatUint(gasPrice*gasUsed, 10)

	// sss, _ := receipt.MarshalJSON()
	return &TransferDetail{
		HashString: hashString,
		Status:     statusInt,

		From:         fromAddress,
		To:           toAddress,
		Amount:       amount,
		EstimateFees: estimateFees,

		FinishTimestamp: int64(blockHeader.Time),
	}, nil
}

// 获取交易的状态
// @param hashString 交易的 hash
func (e *EthChain) FetchTransferStatus(hashString string) int {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	_, isPending, err := e.RemoteRpcClient.TransactionByHash(ctx, common.HexToHash(hashString))
	if err != nil {
		return 0
	}
	if isPending {
		return 1
	}

	// 交易receipt 状态信息，0表示失败，1表示成功
	receipt, err := e.TransactionReceiptByHash(hashString)
	if receipt.Status == 0 {
		return 3
	} else {
		return 2
	}
}
