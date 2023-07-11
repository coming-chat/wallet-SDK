package aptos

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	txbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type Transaction struct {
	RawTxn txbuilder.RawTransaction
}

func (t *Transaction) SignWithAccount(account base.Account) (signedTx *base.OptionalString, err error) {
	txn, err := t.SignedTransactionWithAccount(account)
	if err != nil {
		return nil, err
	}
	return txn.HexString()
}

func (t *Transaction) SignedTransactionWithAccount(account base.Account) (signedTx base.SignedTransaction, err error) {
	acc, ok := account.(*Account)
	if !ok {
		return nil, base.ErrInvalidAccountType
	}
	signedBytes, err := txbuilder.GenerateBCSTransaction(acc.account, &t.RawTxn)
	if err != nil {
		return nil, err
	}
	return &SignedTransaction{
		RawTxn:      &t.RawTxn,
		SignedBytes: signedBytes,
	}, nil
}

type SignedTransaction struct {
	RawTxn *txbuilder.RawTransaction

	SignedBytes []byte
}

func (txn *SignedTransaction) HexString() (res *base.OptionalString, err error) {
	return &base.OptionalString{Value: types.HexEncodeToString(txn.SignedBytes)}, nil
}

func AsSignedTransaction(txn base.SignedTransaction) *SignedTransaction {
	if res, ok := txn.(*SignedTransaction); ok {
		return res
	}
	return nil
}
