package eth

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/coming-chat/wallet-SDK/core/base"
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
