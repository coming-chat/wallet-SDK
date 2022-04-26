package eth

import (
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// 合约 ABI json文件，查询ERC20 相关代币信息需要使用 ABI 文件
	Erc20AbiStr = `[{"inputs":[{"internalType":"address","name":"operator","type":"address"},{"internalType":"address","name":"pauser","type":"address"},{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"symbol","type":"string"},{"internalType":"uint8","name":"decimal","type":"uint8"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Unpaused","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"burn","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"new_operator","type":"address"},{"internalType":"address","name":"new_pauser","type":"address"}],"name":"changeUser","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"pause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unpause","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

// Erc20Token  erc20 代币对象，定义了代币的基础信息：合约地址、代币名、代币符号、代币精度、代币余额
type Erc20Token struct {
	ContractAddress string // 代币合约地址
	Name            string // 代币名称
	Symbol          string // 代币符号
	Decimal         string // 代币精度
	ChainId         string // 链ID
	Balance         string // 代币余额
	TokenIcon       string // 代币图标
}

type Erc20Balance struct {
	ContractAddress string
	Balance         string
}

// @title    Erc20代币基础信息
// @description   返回代币基础信息
// @auth      清欢
// @param     (contractAddress, walletAddress)     (string,string)  合约名称，钱包地址
// @return    (*Erc20Token,error)       Erc20Token，错误信息
func (e *EthChain) Erc20TokenInfo(contractAddress string, walletAddress string) (*Erc20Token, error) {
	var token Erc20Token
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

// SDK 批量请求代币余额，因为 golang 导出的方法不支持数组，因此传参和返回都只能用字符串
// @param contractListString 批量查询的代币的合约地址字符串，用逗号隔开，例如: "add1,add2,add3"
// @param address 用户的钱包地址
// @return 余额列表，顺序与传入合约顺序是保持一致的，例如: "b1,b2,b3"
// @throw 如果任意一个代币请求余额出错时，会抛出错误
func (e *EthChain) SdkBatchTokenBalance(contractListString, address string) (string, error) {
	contractList := strings.Split(contractListString, ",")
	balances, err := e.BatchTokenBalance(contractList, address)
	if err != nil {
		return "", err
	}
	return strings.Join(balances, ","), nil
}

// 批量请求代币余额
// @param contractList 批量查询的代币的合约地址数组
// @param address 用户的钱包地址
// @return 余额数组，顺序与传入的 contractList 是保持一致的
// @throw 如果任意一个代币请求余额出错时，会抛出错误
func (e *EthChain) BatchTokenBalance(contractList []string, address string) ([]string, error) {
	return MapListConcurrentStringToString(contractList, func(s string) (string, error) {
		return e.TokenBalance(s, address)
	})
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
func (e *EthChain) TokenDecimal(contractAddress string) (string, error) {
	result := uint8(0)
	err := e.CallContractConstant(
		&result,
		contractAddress,
		Erc20AbiStr,
		"decimals",
		nil,
	)
	if err != nil {
		return "0", err
	}
	tokenDecimal := strconv.Itoa(int(result))
	return tokenDecimal, err
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
