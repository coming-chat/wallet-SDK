package eth

import (
	"context"

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
		return "0", MapToBasicError(err)
	}
	return result.String(), nil
}
