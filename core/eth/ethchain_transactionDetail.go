package eth

import (
	"context"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
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

// 获取交易的详情
// @param hashString 交易的 hash
// @return 详情对象，该对象无法提供 CID 信息
func (e *EthChain) FetchTransactionDetail(hashString string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	tx, isPending, err := e.RemoteRpcClient.TransactionByHash(ctx, common.HexToHash(hashString))
	if err != nil {
		return
	}
	msg, err := tx.AsMessage(types.NewLondonSigner(e.chainId), nil)
	if err != nil {
		return
	}

	address := msg.To().String()
	amount := strconv.FormatUint(msg.Value().Uint64(), 10)

	if len(tx.Data()) != 0 {
		address, amount, err = decodeErc20TransferInput(tx.Data())
		if err != nil {
			return
		}
	}

	gasPrice := msg.GasPrice().Uint64()
	estimateGasLimit := msg.Gas()
	detail = &base.TransactionDetail{
		HashString:   hashString,
		FromAddress:  msg.From().String(),
		ToAddress:    address,
		Amount:       amount,
		EstimateFees: strconv.FormatUint(gasPrice*estimateGasLimit, 10),
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
			err = nil
		}

	} else {
		detail.Status = base.TransactionStatusSuccess
	}
	gasUsed := receipt.GasUsed
	detail.EstimateFees = strconv.FormatUint(gasPrice*gasUsed, 10)
	detail.FinishTimestamp = int64(blockHeader.Time)

	return detail, nil
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

// 根据交易hash查询交易状态
// TO-DO  返回更详细的信息，解析交易余额，交易动作
func (e *EthChain) TransactionByHash(txHash string) (*TransactionByHashResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	tx, isPending, err := e.RemoteRpcClient.TransactionByHash(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(e.chainId), nil)
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
		msg.From(),
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
func (e *EthChain) TransactionReceiptByHash(txHash string) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	receipt, err := e.RemoteRpcClient.TransactionReceipt(ctx, common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (e *EthChain) WaitConfirm(txHash string, interval time.Duration) *types.Receipt {
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
