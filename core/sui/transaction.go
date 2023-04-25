package sui

import (
	"encoding/json"
	"errors"

	"github.com/coming-chat/go-sui/sui_types"
	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type Transaction struct {
	Txn types.TransactionBytes

	EstimateGasFee int64
}

type SignedTransaction struct {
	// transaction data bytes
	TxBytes *types.Base64Data `json:"tx_bytes"`

	// transaction signature
	Signature *sui_types.Signature `json:"signature"`
}

func (txn *Transaction) SignWithAccount(account *Account) (signedTx *base.OptionalString, err error) {
	if account == nil {
		return nil, errors.New("Invalid account.")
	}
	signature, err := account.account.SignSecureWithoutEncode(txn.Txn.TxBytes, sui_types.DefaultIntent())
	if err != nil {
		return nil, err
	}
	signedTxn := SignedTransaction{
		TxBytes:   &txn.Txn.TxBytes,
		Signature: &signature,
	}
	bytes, err := json.Marshal(signedTxn)
	if err != nil {
		return
	}
	txnString := types.Bytes(bytes).GetBase64Data().String()

	return &base.OptionalString{Value: txnString}, nil
}
