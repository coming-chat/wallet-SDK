package eth

import (
	"context"
	"strconv"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/common"
)

// @title    主网代币余额查询
// @description   返回主网代币余额，decimal为代币精度
// @auth      清欢
// @param     (walletAddress)     (string)  合约名称，钱包地址
// @return    (string,error)       代币余额，错误信息
func (e *EthChain) Balance(address string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	result, err := e.RemoteRpcClient.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return "0", base.MapAnyToBasicError(err)
	}
	return result.String(), nil
}

// 获取最新区块高度
func (e *EthChain) LatestBlockNumber() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	number, err := e.RemoteRpcClient.BlockNumber(ctx)
	if err != nil {
		return 0, base.MapAnyToBasicError(err)
	}

	return int64(number), nil
}

// 获取账户nonce
func (e *EthChain) Nonce(spenderAddressHex string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	nonce, err := e.RemoteRpcClient.PendingNonceAt(ctx, common.HexToAddress(spenderAddressHex))
	if err != nil {
		return "0", base.MapAnyToBasicError(err)
	}
	return strconv.FormatUint(nonce, 10), nil
}

// 获取链ID
func GetChainId(e *EthChain) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	chainId, err := e.RemoteRpcClient.ChainID(ctx)
	if err != nil {
		return "0", base.MapAnyToBasicError(err)
	}

	return chainId.String(), nil
}

// 对交易进行广播
func (e *EthChain) SendRawTransaction(txHex string) (string, error) {
	var hash common.Hash
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	err := e.RpcClient.CallContext(ctx, &hash, "eth_sendRawTransaction", txHex)
	if err != nil {
		return "", base.MapAnyToBasicError(err)
	}
	return hash.String(), nil
}
