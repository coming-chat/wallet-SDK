package aptos

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	txbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type Transaction struct {
	RawTxn txbuilder.RawTransaction
}

// type SignedTransaction struct {
// 	SignedBytes []byte
// }

func (t *Transaction) SignWithAccount(account base.Account) (signedTx *base.OptionalString, err error) {
	acc, ok := account.(*Account)
	if !ok {
		return nil, base.ErrInvalidAccountType
	}
	signedBytes, err := txbuilder.GenerateBCSTransaction(acc.account, &t.RawTxn)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: types.HexEncodeToString(signedBytes)}, nil
}
