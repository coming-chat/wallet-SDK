package sui

import (
	"encoding/json"
	"errors"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type Transaction struct {
	Txn types.TransactionBytes

	MaxGasBudget   int64
	EstimateGasFee int64
}

func (txn *Transaction) SignWithAccount(account *Account) (signedTx *base.OptionalString, err error) {
	if account == nil {
		return nil, errors.New("Invalid account.")
	}
	signedTxn := txn.Txn.SignSerializedSigWith(account.account.PrivateKey)
	bytes, err := json.Marshal(signedTxn)
	if err != nil {
		return
	}
	txnString := types.Bytes(bytes).GetBase64Data().String()

	return &base.OptionalString{Value: txnString}, nil
}
