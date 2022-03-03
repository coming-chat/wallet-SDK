package eth

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

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
	return &TransactionByHashResult{
		tx,
		msg.From(),
		isPending,
	}, nil
}

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
