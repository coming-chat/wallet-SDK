package solana

import (
	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/portto/solana-go-sdk/types"
)

type Transaction struct {
	Message types.Message
}

func (t *Transaction) SignWithAccount(account base.Account) (signedTxn *base.OptionalString, err error) {
	txn, err := t.SignedTransactionWithAccount(account)
	if err != nil {
		return nil, err
	}
	return txn.HexString()
}

func (t *Transaction) SignedTransactionWithAccount(account base.Account) (signedTxn base.SignedTransaction, err error) {
	solanaAcc := AsSolanaAccount(account)
	if solanaAcc == nil {
		return nil, base.ErrInvalidAccountType
	}

	// create tx by message + signer
	txn, err := types.NewTransaction(types.NewTransactionParam{
		Message: t.Message,
		Signers: []types.Account{*solanaAcc.account, *solanaAcc.account},
	})
	if err != nil {
		return nil, err
	}
	return &SignedTransaction{
		Transaction: txn,
	}, nil
}

type SignedTransaction struct {
	Transaction types.Transaction
}

func (txn *SignedTransaction) HexString() (res *base.OptionalString, err error) {
	bytes, err := txn.Transaction.Serialize()
	if err != nil {
		return nil, err
	}
	hexString := hexTypes.HexEncodeToString(bytes)

	return &base.OptionalString{Value: hexString}, nil
}

func AsSignedTransaction(txn base.SignedTransaction) *SignedTransaction {
	if res, ok := txn.(*SignedTransaction); ok {
		return res
	}
	return nil
}
