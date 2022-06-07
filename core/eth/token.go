package eth

import (
	"crypto/ecdsa"
	"errors"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/crypto"
)

type TokenProtocol interface {
	base.Token

	EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string) (string, error)
	BuildTransferTx(privateKey string, transaction *Transaction) (*base.OptionalString, error)
	BuildTransferTxWithAccount(account *Account, transaction *Transaction) (*base.OptionalString, error)
}

type Token struct {
	chain *Chain
}

// Warning: initial unavailable, You must create based on Chain.MainToken()
func NewToken() (*Token, error) {
	return nil, errors.New("Token initial unavailable, You must create based on Chain.MainToken()")
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

// MARK - Eth TokenProtocol

func (t *Token) EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string) (string, error) {
	msg := NewCallMsg()
	msg.SetFrom(fromAddress)
	msg.SetTo(receiverAddress)
	msg.SetGasPrice(gasPrice)
	msg.SetValue(amount)

	res, err := t.chain.EstimateGasLimit(msg)
	if err != nil {
		return "", err
	}
	return res.Value, nil
}

func (t *Token) BuildTransferTx(privateKey string, transaction *Transaction) (*base.OptionalString, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	return t.buildTransfer(privateKeyECDSA, transaction)
}

func (t *Token) BuildTransferTxWithAccount(account *Account, transaction *Transaction) (*base.OptionalString, error) {
	return t.buildTransfer(account.privateKeyECDSA, transaction)
}

func (t *Token) buildTransfer(privateKey *ecdsa.PrivateKey, transaction *Transaction) (*base.OptionalString, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return nil, err
	}

	if transaction.Nonce == "" || transaction.Nonce == "0" {
		address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
		nonce, err := chain.Nonce(address)
		if err != nil {
			nonce = "0"
			err = nil
		}
		transaction.Nonce = nonce
	}

	rawTx, err := transaction.GetRawTx()
	if err != nil {
		return nil, err
	}

	txResult, err := chain.buildTxWithTransaction(rawTx, privateKey)
	if err != nil {
		return nil, err
	}

	return &base.OptionalString{Value: txResult.TxHex}, nil
}
