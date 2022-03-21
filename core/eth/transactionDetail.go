package eth

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

	address := msg.To().String()
	amount := strconv.FormatUint(msg.Value().Uint64(), 10)

	if len(tx.Data()) != 0 {
		address, amount, err = decodeErc20TransferInput(tx.Data())
		if err != nil {
			return nil, err
		}
	}

	gasPrice := msg.GasPrice().Uint64()
	estimateGasLimit := msg.Gas()
	detail := &TransactionDetail{
		HashString:   hashString,
		FromAddress:  msg.From().String(),
		ToAddress:    address,
		Amount:       amount,
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
			From:       msg.From(),
			To:         msg.To(),
			Data:       msg.Data(),
			Gas:        msg.Gas(),
			GasPrice:   msg.GasPrice(),
			GasFeeCap:  msg.GasFeeCap(),
			GasTipCap:  msg.GasTipCap(),
			Value:      msg.Value(),
			AccessList: msg.AccessList(),
		}, receipt.BlockNumber)
		if err != nil {
			detail.FailureMessage = err.Error()
		}

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
	if len(hashString) == 0 {
		return TransactionStatusNone
	}

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

// 批量获取交易的转账状态
// @param hashList 要批量查询的交易的 hash 数组
// @return 交易状态数组，它的顺序和 hashList 是保持一致的
func (e *EthChain) BatchTransactionStatus(hashList []string) []string {
	thread := 0
	max := 10
	wg := sync.WaitGroup{}

	dict := &safeMap{Map: make(map[string]string)}
	for _, hashString := range hashList {
		if thread == max {
			wg.Wait()
			thread = 0
		}
		if thread < max {
			wg.Add(1)
		}

		go func(w *sync.WaitGroup, hashString string, dict *safeMap) {
			status := e.FetchTransactionStatus(hashString)
			dict.writeMap(hashString, strconv.Itoa(status))
			wg.Done()
		}(&wg, hashString, dict)
		thread++
	}
	wg.Wait()

	result := []string{}
	for _, hashString := range hashList {
		result = append(result, dict.Map[hashString])
	}
	return result
}

// SDK 批量获取交易的转账状态，hash 列表和返回值，都只能用字符串，逗号隔开传递
// @param hashListString 要批量查询的交易的 hash，用逗号拼接的字符串："hash1,hash2,hash3"
// @return 批量的交易状态，它的顺序和 hashListString 是保持一致的: "status1,status2,status3"
func (e *EthChain) SdkBatchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	array := e.BatchTransactionStatus(hashList)
	return strings.Join(array, ",")
}

// 解析 erc20 转账的 input data
// @return 返回转账地址和金额
func decodeErc20TransferInput(data []byte) (string, string, error) {
	if len(data) == 0 {
		// 主币
		return "", "", nil
	}

	parsedAbi, err := abi.JSON(strings.NewReader(Erc20AbiStr))
	if err != nil {
		return "", "", err
	}

	method, err := parsedAbi.MethodById(data[:4])
	if err != nil {
		return "", "", err
	}
	if method.RawName != ERC20_METHOD_TRANSFER {
		return "", "", nil
	}

	params, err := method.Inputs.Unpack(data[4:])
	address := params[0].(common.Address).String()
	amount := params[1].(*big.Int).String()
	return address, amount, nil
}
