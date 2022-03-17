package eth

import (
	"fmt"
	"strconv"
)

type CoinUtil struct {
	// 链的 RPC 地址
	RpcUrl string
	// 币种的合约地址，如果为 nil，表示是主网的币
	ContractAddress string
	// 用户的钱包地址
	WalletAddress string
}

// 创建 CoinUtil 对象
// @param contractAddress 币种的合约地址，如果是主网的币，可传 nil
// @param walletAddress 用户的钱包地址
func NewCoinUtilWithRpc(rpcUrl, contractAddress, walletAddress string) *CoinUtil {
	return &CoinUtil{
		RpcUrl:          rpcUrl,
		ContractAddress: contractAddress,
		WalletAddress:   walletAddress,
	}
}

// 是否是主币
func (u *CoinUtil) IsMainCoin() bool {
	return u.ContractAddress == ""
}

// 查询币种信息
// 主币只能获取到余额
// 合约币能获取到所有信息
func (u *CoinUtil) CoinInfo() (*Erc20Token, error) {
	chain, err := GetConnection(u.RpcUrl)
	if err != nil {
		return nil, err
	}

	if u.IsMainCoin() {
		balance, err := u.QueryBalance()
		if err != nil {
			return nil, err
		}
		return &Erc20Token{
			Balance: balance,
		}, nil
	} else {
		return chain.Erc20TokenInfo(u.ContractAddress, u.WalletAddress)
	}
}

// 查询用户的钱包余额
func (u *CoinUtil) QueryBalance() (string, error) {
	chain, err := GetConnection(u.RpcUrl)
	if err != nil {
		return "", err
	}

	var balance string
	if u.IsMainCoin() {
		balance, err = chain.Balance(u.WalletAddress)
	} else {
		balance, err = chain.TokenBalance(u.ContractAddress, u.WalletAddress)
	}
	if err != nil {
		return "", err
	}

	return balance, nil
}

// 获取 gasPrice
func (u *CoinUtil) SuggestGasPrice() (string, error) {
	chain, err := GetConnection(u.RpcUrl)
	if err != nil {
		return "", err
	}
	return chain.SuggestGasPrice()
}

// 获取交易的 nonce
func (u *CoinUtil) Nonce() (string, error) {
	chain, err := GetConnection(u.RpcUrl)
	if err != nil {
		return "", err
	}
	return chain.Nonce(u.WalletAddress)
}

// 获取转账的 预估 gasLimit
func (u *CoinUtil) EstimateGasLimit(receiverAddress, gasPrice, amount string) (string, error) {
	chain, err := GetConnection(u.RpcUrl)
	if err != nil {
		return "", err
	}

	var gasLimit string
	if u.IsMainCoin() {
		gasLimit, err = chain.EstimateGasLimit(u.WalletAddress, receiverAddress, gasPrice, amount)
	} else {
		erc20JsonParams := fmt.Sprintf(
			"{\"toAddress\":\"%s\", \"amount\":\"%s\", \"method\":\"%s\"}",
			receiverAddress,
			amount,
			ERC20_METHOD_TRANSFER)
		gasLimit, err = chain.EstimateContractGasLimit(u.WalletAddress, u.ContractAddress, Erc20AbiStr, ERC20_METHOD_TRANSFER, erc20JsonParams)
	}
	if err != nil {
		return "", err
	}

	return gasLimit, nil
}

// 创建 ETH 转账交易
func (u *CoinUtil) BuildTransferTx(privateKey, receiverAddress, nonce, gasPrice, gasLimit, amount string) (string, error) {
	chain, err := GetConnection(u.RpcUrl)
	if err != nil {
		return "", err
	}

	nonceInt, err := strconv.ParseInt(nonce, 10, 64)
	if err != nil {
		return "", err
	}
	call := &CallMethodOpts{
		Nonce:    nonceInt,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
	}

	var output *BuildTxResult
	if u.IsMainCoin() {
		call.Value = amount
		output, err = chain.BuildTransferTx(privateKey, receiverAddress, call)

	} else {
		erc20JsonParams := fmt.Sprintf(
			"{\"toAddress\":\"%s\", \"amount\":\"%s\", \"method\":\"%s\"}",
			receiverAddress,
			amount,
			ERC20_METHOD_TRANSFER)
		output, err = chain.BuildCallMethodTx(privateKey, u.ContractAddress, Erc20AbiStr, ERC20_METHOD_TRANSFER, call, erc20JsonParams)
	}
	if err != nil {
		return "", err
	}

	return output.TxHex, nil
}

// 对交易进行广播
func (u *CoinUtil) SendRawTransaction(txHex string) (string, error) {
	chain, err := GetConnection(u.RpcUrl)
	if err != nil {
		return "", err
	}
	return chain.SendRawTransaction(txHex)
}

// 获取交易的状态
// @param hashString 交易的 hash
func (u *CoinUtil) FetchTransactionStatus(hashString string) TransactionStatus {
	chain, err := GetConnection(u.RpcUrl)
	if err != nil {
		return TransactionStatusNone
	}
	return chain.FetchTransactionStatus(hashString)
}

// 获取交易的详情
// @param hashString 交易的 hash
// @return 详情对象，该对象无法提供 CID 信息
func (u *CoinUtil) FetchTransactionDetail(hashString string) (*TransactionDetail, error) {
	chain, err := GetConnection(u.RpcUrl)
	if err != nil {
		return nil, err
	}
	return chain.FetchTransactionDetail(hashString)
}
