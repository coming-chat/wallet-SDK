package solana

import (
	"context"
	"errors"
	"strconv"

	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/types"
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
	account, err := NewAccountWithPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return t.BuildTransferTxWithAccount(account, receiverAddress, amount)
}

func (t *Token) BuildTransferTxWithAccount(account *Account, receiverAddress, amount string) (*base.OptionalString, error) {
	client := t.chain.client()
	message, err := transactionMessage(client, account.Address(), receiverAddress, amount)
	if err != nil {
		return nil, err
	}

	// create tx by message + signer
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: *message,
		Signers: []types.Account{*account.account, *account.account},
	})
	if err != nil {
		return nil, err
	}

	bytes, err := tx.Serialize()
	if err != nil {
		return nil, err
	}
	hash := hexTypes.HexEncodeToString(bytes)

	return &base.OptionalString{Value: hash}, nil
}

func (t *Token) EstimateFees(receiverAddress, amount string) (*base.OptionalString, error) {
	client := t.chain.client()
	message, err := transactionMessage(client, receiverAddress, receiverAddress, amount)
	if err != nil {
		return nil, err
	}

	fee, err := client.GetFeeForMessage(context.Background(), *message)
	if err != nil {
		return nil, err
	}
	feeString := strconv.FormatUint(*fee, 10)

	return &base.OptionalString{Value: feeString}, nil
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
	res, err := client.GetRecentBlockhash(context.Background())
	if err != nil {
		return nil, err
	}

	// create a message
	message := types.NewMessage(types.NewMessageParam{
		FeePayer:        pubFrom,
		RecentBlockhash: res.Blockhash, // recent blockhash
		Instructions: []types.Instruction{
			sysprog.Transfer(sysprog.TransferParam{
				From:   pubFrom, // from
				To:     pubTo,   // to
				Amount: amountUint,
			}),
		},
	})

	return &message, nil
}
