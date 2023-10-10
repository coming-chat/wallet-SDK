package solana

import (
	"context"
	"errors"
	"strconv"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type Token struct {
	chain *Chain
}

func NewToken(chain *Chain) *Token {
	return &Token{chain}
}

// MARK - Implement the protocol Token

func (t *Token) Chain() base.Chain {
	return t.chain
}

func (t *Token) TokenInfo() (*base.TokenInfo, error) {
	return &base.TokenInfo{
		Name:    "SOL",
		Symbol:  "SOL",
		Decimal: 9,
	}, nil
}

func (t *Token) BalanceOfAddress(address string) (*base.Balance, error) {
	return t.chain.BalanceOfAddress(address)
}

func (t *Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return t.chain.BalanceOfPublicKey(publicKey)
}

func (t *Token) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.BalanceOfAddress(account.Address())
}

// MARK - Solana token

func (t *Token) BuildTransferTx(privateKey, receiverAddress, amount string) (*base.OptionalString, error) {
	account, err := AccountWithPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return t.BuildTransferTxWithAccount(account, receiverAddress, amount)
}

func (t *Token) BuildTransferTxWithAccount(account *Account, receiverAddress, amount string) (*base.OptionalString, error) {
	txn, err := t.BuildTransfer(account.Address(), receiverAddress, amount)
	if err != nil {
		return nil, err
	}
	signedTxn, err := txn.SignedTransactionWithAccount(account)
	if err != nil {
		return nil, err
	}
	return signedTxn.HexString()
}

func (t *Token) EstimateFees(receiverAddress, amount string) (*base.OptionalString, error) {
	txn, err := t.BuildTransfer(receiverAddress, receiverAddress, amount)
	if err != nil {
		return nil, err
	}
	return t.chain.EstimateTransactionFee(txn)
}

func transactionMessage(client *client.Client, fromAddress, toAddress, amount string) (*types.Message, error) {
	if !IsValidAddress(toAddress) {
		return nil, errors.New("Invalid receiver address")
	}
	amountUint, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return nil, errors.New("Invalid amount")
	}
	pubTo := common.PublicKeyFromString(toAddress)
	pubFrom := common.PublicKeyFromString(fromAddress) // from is same as to, or it's must valid

	// to fetch recent blockhash
	res, err := client.GetLatestBlockhash(context.Background())
	if err != nil {
		return nil, err
	}

	// create a message
	message := types.NewMessage(types.NewMessageParam{
		FeePayer:        pubFrom,
		RecentBlockhash: res.Blockhash, // recent blockhash
		Instructions: []types.Instruction{
			system.Transfer(system.TransferParam{
				From:   pubFrom, // from
				To:     pubTo,   // to
				Amount: amountUint,
			}),
		},
	})

	return &message, nil
}

func (t *Token) BuildTransfer(sender, receiver, amount string) (txn base.Transaction, err error) {
	client := t.chain.client()
	message, err := transactionMessage(client, sender, receiver, amount)
	if err != nil {
		return nil, err
	}
	return &Transaction{
		Message: *message,
	}, nil
}
func (t *Token) CanTransferAll() bool {
	return false
}
func (t *Token) BuildTransferAll(sender, receiver string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}
