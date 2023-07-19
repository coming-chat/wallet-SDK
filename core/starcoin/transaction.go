package starcoin

import (
	"encoding/hex"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/starcoinorg/starcoin-go/client"
	"github.com/starcoinorg/starcoin-go/types"
)

type Transaction struct {
	Txn *types.RawUserTransaction
}

func (t *Transaction) SignWithAccount(account base.Account) (signedTxn *base.OptionalString, err error) {
	txn, err := t.SignedTransactionWithAccount(account)
	if err != nil {
		return nil, err
	}
	return txn.HexString()
}

func (t *Transaction) SignedTransactionWithAccount(account base.Account) (signedTxn base.SignedTransaction, err error) {
	starcoinAcc := AsStarcoinAccount(account)
	if starcoinAcc == nil {
		return nil, base.ErrInvalidAccountType
	}

	privateKey, err := account.PrivateKey()
	if err != nil {
		return
	}
	txn, err := client.SignRawUserTransaction(types.Ed25519PrivateKey(privateKey), t.Txn)
	if err != nil {
		return
	}
	return &SignedTransaction{
		Txn: txn,
	}, nil
}

type SignedTransaction struct {
	Txn *types.SignedUserTransaction
}

func (txn *SignedTransaction) HexString() (res *base.OptionalString, err error) {
	txnBytes, err := txn.Txn.BcsSerialize()
	if err != nil {
		return nil, err
	}
	hexString := "0x" + hex.EncodeToString(txnBytes)

	return &base.OptionalString{Value: hexString}, nil
}

func AsSignedTransaction(txn base.SignedTransaction) *SignedTransaction {
	if res, ok := txn.(*SignedTransaction); ok {
		return res
	}
	return nil
}
