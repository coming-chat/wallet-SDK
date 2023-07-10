package sui

import (
	"encoding/json"

	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type Transaction struct {
	Txn types.TransactionBytes

	TxnBytes lib.Base64Data

	EstimateGasFee int64
}

type SignedTransaction struct {
	// transaction data bytes
	TxBytes *lib.Base64Data `json:"tx_bytes"`

	// transaction signature
	Signature *sui_types.Signature `json:"signature"`
}

func (t *Transaction) TransactionBytes() []byte {
	if t.TxnBytes != nil {
		return t.TxnBytes
	}
	return t.Txn.TxBytes
}

func (t *Transaction) SignWithAccount(account base.Account) (signedTx *base.OptionalString, err error) {
	signedTxn, err := t.SignedTransactionWithAccount(account)
	if err != nil {
		return
	}
	bytes, err := json.Marshal(signedTxn)
	if err != nil {
		return
	}
	txnString := lib.Base64Data(bytes).String()

	return &base.OptionalString{Value: txnString}, nil
}

func (t *Transaction) SignedTransactionWithAccount(account base.Account) (signedTx base.SignedTransaction, err error) {
	acc, ok := account.(*Account)
	if !ok {
		return nil, base.ErrInvalidAccountType
	}
	txnBytes := t.TransactionBytes()
	signature, err := acc.account.SignSecureWithoutEncode(txnBytes, sui_types.DefaultIntent())
	if err != nil {
		return nil, err
	}
	base64data := lib.Base64Data(txnBytes)
	signedTxn := SignedTransaction{
		TxBytes:   &base64data,
		Signature: &signature,
	}
	return &signedTxn, nil
}

func AsSignedTransaction(txn base.SignedTransaction) *SignedTransaction {
	if res, ok := txn.(*SignedTransaction); ok {
		return res
	}
	if res, ok := txn.(SignedTransaction); ok {
		return &res
	}
	return nil
}
