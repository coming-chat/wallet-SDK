package eth

import (
	"context"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

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
	// 当交易没有处于pending状态时，可以查询receipt信息，即交易是否成功
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
