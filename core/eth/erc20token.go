package eth

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/common"
)

var (
	// 合约 ABI json文件，查询ERC20 相关代币信息需要使用 ABI 文件
	Erc20AbiStr = `[{"inputs":[{"internalType":"address","name":"operator","type":"address"},{"internalType":"address","name":"pauser","type":"address"},{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"symbol","type":"string"},{"internalType":"uint8","name":"decimal","type":"uint8"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Unpaused","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"burn","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"new_operator","type":"address"},{"internalType":"address","name":"new_pauser","type":"address"}],"name":"changeUser","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"pause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unpause","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

type Erc20TokenInfo struct {
	*base.TokenInfo

	ContractAddress string
	ChainId         string
	TokenIcon       string

	// Deprecated: Balance is not a token's info.
	Balance string
}

type Erc20Token struct {
	*Token

	ContractAddress string
}

// Warning: initial unavailable, You must create based on Chain.Erc20Token()
func NewErc20Token() (*Erc20Token, error) {
	return nil, errors.New("Token initial unavailable, You must create based on Chain.MainToken()")
}

func (c *Chain) Erc20Token(contractAddress string) base.Token {
	return &Erc20Token{
		Token:           &Token{chain: c},
		ContractAddress: contractAddress,
	}
}

// MARK - Implement the protocol Token, Override

// cannot get balance
func (t *Erc20Token) Erc20TokenInfo() (*Erc20TokenInfo, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return nil, err
	}
	baseInfo, err := t.TokenInfo()
	if err != nil {
		return nil, err
	}

	return &Erc20TokenInfo{
		TokenInfo:       baseInfo,
		ContractAddress: t.ContractAddress,
		ChainId:         chain.chainId.String(),
	}, nil
}

func (t *Erc20Token) TokenInfo() (*base.TokenInfo, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return nil, err
	}

	info := &base.TokenInfo{}
	info.Name, err = chain.TokenName(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	info.Symbol, err = chain.TokenSymbol(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	info.Decimal, err = chain.TokenDecimal(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (t *Erc20Token) BalanceOfAddress(address string) (*base.Balance, error) {
	b := base.EmptyBalance()
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return b, err
	}

	balance, err := chain.TokenBalance(t.ContractAddress, address)
	if err != nil {
		return b, err
	}
	return &base.Balance{
		Total:  balance,
		Usable: balance,
	}, nil
}

func (t *Erc20Token) EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string) (string, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return "", err
	}

	erc20JsonParams := fmt.Sprintf(
		"{\"toAddress\":\"%s\", \"amount\":\"%s\", \"method\":\"%s\"}",
		receiverAddress,
		amount,
		ERC20_METHOD_TRANSFER)
	gasLimit, err := chain.EstimateContractGasLimit(fromAddress, t.ContractAddress, Erc20AbiStr, ERC20_METHOD_TRANSFER, erc20JsonParams)

	return gasLimit, err
}

func (t *Erc20Token) BuildTransferTx(privateKey, fromAddress, receiverAddress, gasPrice, gasLimit, amount string) (string, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return "", err
	}

	nonce, err := t.chain.NonceOfAddress(fromAddress)
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

	erc20JsonParams := fmt.Sprintf(
		"{\"toAddress\":\"%s\", \"amount\":\"%s\", \"method\":\"%s\"}",
		receiverAddress,
		amount,
		ERC20_METHOD_TRANSFER)
	output, err := chain.BuildCallMethodTx(privateKey, t.ContractAddress, Erc20AbiStr, ERC20_METHOD_TRANSFER, call, erc20JsonParams)
	if err != nil {
		return "", err
	}

	return output.TxHex, nil
}

// Deprecated: Erc20TokenInfo is deprecated. Please Use Chain.Erc20Token().Erc20TokenInfo()
// @title    Erc20代币基础信息
// @description   返回代币基础信息
// @auth      清欢
// @param     (contractAddress, walletAddress)     (string,string)  合约名称，钱包地址
// @return    (*Erc20Token,error)       Erc20Token，错误信息
func (e *EthChain) Erc20TokenInfo(contractAddress string, walletAddress string) (*Erc20TokenInfo, error) {
	var token Erc20TokenInfo
	token.ContractAddress = contractAddress
	token.ChainId = e.chainId.String()
	var err error
	token.Decimal, err = e.TokenDecimal(contractAddress)
	if err != nil {
		return nil, err
	}
	token.Symbol, err = e.TokenSymbol(contractAddress)
	if err != nil {
		return nil, err
	}
	token.Name, err = e.TokenName(contractAddress)
	if err != nil {
		return nil, err
	}
	token.Balance, err = e.TokenBalance(contractAddress, walletAddress)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// @title    Erc20代币余额
// @description   返回erc20代币余额
// @auth      清欢
// @param     (contractAddress，walletAddress)     合约地址,钱包地址
// @return    (string,error)       余额，错误信息
func (e *EthChain) TokenBalance(contractAddress, address string) (string, error) {
	if len(contractAddress) == 0 || len(address) == 0 {
		return "0", errors.New("The address of the contract or wallet is empty.")
	}
	result := new(big.Int)
	err := e.CallContractConstant(
		&result,
		contractAddress,
		Erc20AbiStr,
		"balanceOf",
		nil,
		common.HexToAddress(address),
	)
	if err != nil {
		return "0", err
	}

	return result.String(), err
}

// @title    Erc20代币精度
// @description   返回代币精度
// @auth      清欢
// @param     (contractAddress)     合约地址
// @return    (string,error)       代币精度，错误信息
func (e *EthChain) TokenDecimal(contractAddress string) (int16, error) {
	decimal := int16(0)
	err := e.CallContractConstant(
		&decimal,
		contractAddress,
		Erc20AbiStr,
		"decimals",
		nil,
	)
	if err != nil {
		return 0, err
	}
	return decimal, err
}

// @title    Erc20代币符号
// @description   返回代币符号
// @auth      清欢
// @param     (contractAddress)     合约地址
// @return    (string,error)       符号，错误信息
func (e *EthChain) TokenSymbol(contractAddress string) (string, error) {
	tokenSymbol := ""
	err := e.CallContractConstant(
		&tokenSymbol,
		contractAddress,
		Erc20AbiStr,
		"symbol",
		nil,
	)
	if err != nil {
		return "", err
	}

	return tokenSymbol, err
}

// @title    Erc20代币名称
// @description   返回代币名称
// @auth      清欢
// @param     (contractAddress)     合约地址
// @return    (string,error)       名称，错误信息
func (e *EthChain) TokenName(contractAddress string) (string, error) {
	tokenName := ""
	err := e.CallContractConstant(
		&tokenName,
		contractAddress,
		Erc20AbiStr,
		"name",
		nil,
	)
	if err != nil {
		return "", err
	}

	return tokenName, err
}
