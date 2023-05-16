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

	EstimateGasFee int64
}

type SignedTransaction struct {
	// transaction data bytes
	TxBytes *lib.Base64Data `json:"tx_bytes"`

	// transaction signature
	Signature *sui_types.Signature `json:"signature"`
}

func (t *Transaction) SignWithAccount(account base.Account) (signedTx *base.OptionalString, err error) {
	acc, ok := account.(*Account)
	if !ok {
		return nil, base.ErrInvalidAccountType
	}
	signature, err := acc.account.SignSecureWithoutEncode(t.Txn.TxBytes, sui_types.DefaultIntent())
	if err != nil {
		return nil, err
	}
	signedTxn := SignedTransaction{
		TxBytes:   &t.Txn.TxBytes,
		Signature: &signature,
	}
	bytes, err := json.Marshal(signedTxn)
	if err != nil {
		return
	}
	txnString := lib.Base64Data(bytes).String()

	return &base.OptionalString{Value: txnString}, nil
}
