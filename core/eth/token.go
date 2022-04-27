package eth

import (
	"errors"
	"strconv"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type Token struct {
	chain *Chain
}

// Warning: initial unavailable, You must create based on Chain.MainToken()
func NewToken() (*Token, error) {
	return nil, errors.New("Token initial unavailable, You must create based on Chain.MainToken()")
}

func (c *Chain) MainToken() base.Token {
	return &Token{chain: c}
}

// MARK - Implement the protocol Token

func (t *Token) Chain() base.Chain {
	return t.chain
}

// Warning: Main token does not support
func (t *Token) TokenInfo() (*base.TokenInfo, error) {
	return nil, errors.New("Main token does not support")
}

func (t *Token) BalanceOfAddress(address string) (*base.Balance, error) {
	return t.chain.BalanceOfAddress(address)
}

func (t *Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return t.BalanceOfAddress(publicKey)
}

func (t *Token) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.BalanceOfAddress(account.Address())
}

// MARK - Other

func (t *Token) EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string) (string, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return "", err
	}

	gasLimit, err := chain.EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount)
	return gasLimit, err
}

func (t *Token) BuildTransferTx(privateKey, fromAddress, receiverAddress, gasPrice, gasLimit, amount string) (string, error) {
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
	call.Value = amount

	output, err := chain.BuildTransferTx(privateKey, receiverAddress, call)
	if err != nil {
		return "", err
	}
	return output.TxHex, nil
}

func (t *Token) BuildTransferTxWithAccount(account base.Account, receiverAddress, gasPrice, gasLimit, amount string) (string, error) {
	privateKey, err := account.PrivateKeyHex()
	if err != nil {
		return "", err
	}
	return t.BuildTransferTx(privateKey, account.Address(), receiverAddress, gasPrice, gasLimit, amount)
}
