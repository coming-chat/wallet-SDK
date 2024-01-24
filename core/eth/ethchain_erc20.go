package eth

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// 合约 ABI json文件，查询ERC20 相关代币信息需要使用 ABI 文件
	Erc20AbiStr = `[{"inputs":[{"internalType":"address","name":"operator","type":"address"},{"internalType":"address","name":"pauser","type":"address"},{"internalType":"string","name":"name","type":"string"},{"internalType":"string","name":"symbol","type":"string"},{"internalType":"uint8","name":"decimal","type":"uint8"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Unpaused","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"burn","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"new_operator","type":"address"},{"internalType":"address","name":"new_pauser","type":"address"}],"name":"changeUser","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"subtractedValue","type":"uint256"}],"name":"decreaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"addedValue","type":"uint256"}],"name":"increaseAllowance","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"pause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unpause","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	// erc721 的 ABI 文件, 只支持 transferFrom 方法
	Erc721Abi_TransferOnly = `[{"inputs":[{"internalType":"address","name":"from","type":"address"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

// Deprecated: SdkBatchTokenBalance is deprecated. Please Use Chain.BatchFetchErc20TokenBalance() instead.
func (e *EthChain) SdkBatchTokenBalance(contractListString, address string) (string, error) {
	c := NewChainWithRpc(e.rpcUrl)
	return c.BatchFetchErc20TokenBalance(contractListString, address)
}

// Deprecated: Erc20TokenInfo is deprecated. Please Use Chain.Erc20Token().Erc20TokenInfo()
// @title    Erc20代币基础信息
// @description   返回代币基础信息
// @auth      清欢
// @param     (contractAddress, walletAddress)     (string,string)  合约名称，钱包地址
// @return    (*Erc20Token,error)       Erc20Token，错误信息
func (e *EthChain) Erc20TokenInfo(contractAddress string, walletAddress string) (*Erc20TokenInfo, error) {
	var token = Erc20TokenInfo{TokenInfo: &base.TokenInfo{}}
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
	decimal := uint8(0)
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
	return int16(decimal), err
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

// Call Contract

func (e *EthChain) CallContractConstant(out interface{}, contractAddress, abiStr, methodName string, opts *bind.CallOpts, params ...interface{}) (err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	parsedAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return err
	}
	inputParams, err := parsedAbi.Pack(methodName, params...)
	if err != nil {
		return err
	}

	method, ok := parsedAbi.Methods[methodName]
	if !ok {
		return errors.New("method not found")
	}
	err = e.CallContractConstantWithPayload(out, contractAddress, hex.EncodeToString(inputParams), method.Outputs, opts)
	return err
}

func (e *EthChain) CallContractConstantWithPayload(out interface{}, contractAddress, payload string, outputTypes abi.Arguments, opts *bind.CallOpts) (err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if opts == nil {
		opts = new(bind.CallOpts)
	}

	contractAddressObj := common.HexToAddress(contractAddress)

	payload = strings.TrimPrefix(payload, "0x")
	payloadBuf, err := hex.DecodeString(payload)
	if err != nil {
		return err
	}
	var (
		msg    = ethereum.CallMsg{From: opts.From, To: &contractAddressObj, Data: payloadBuf}
		ctx    = opts.Context
		code   []byte
		output []byte
	)
	if ctx == nil {
		ctxTemp, cancel := context.WithTimeout(context.Background(), e.timeout)
		defer cancel()
		ctx = ctxTemp
	}
	if opts.Pending {
		pb := bind.PendingContractCaller(e.RemoteRpcClient)
		output, err = pb.PendingCallContract(ctx, msg)
		if err != nil {
			return err
		}
		if len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = pb.PendingCodeAt(ctx, contractAddressObj); err != nil {
				return err
			} else if len(code) == 0 {
				return errors.New(bind.ErrNoCode.Error())
			}
		}
	} else {
		output, err = bind.ContractCaller(e.RemoteRpcClient).CallContract(ctx, msg, opts.BlockNumber)
		if err != nil {
			return err
		}
		if len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = bind.ContractCaller(e.RemoteRpcClient).CodeAt(ctx, contractAddressObj, opts.BlockNumber); err != nil {
				return err
			} else if len(code) == 0 {
				return errors.New(bind.ErrNoCode.Error())
			}
		}
	}
	err = e.UnpackParams(out, outputTypes, hex.EncodeToString(output))
	if err != nil {
		return err
	}
	return nil
}
