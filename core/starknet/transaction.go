package starknet

import (
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/dontpanicdao/caigo/types"
)

type Transaction struct {
	calls   []types.FunctionCall
	details types.ExecuteDetails
}

type SignedTransaction struct {
	Account *Account

	// depoly Txn
	depolyTxn *DeployAccountTransaction

	// invoke Txn
	invokeTxn *Transaction
}

func (t *Transaction) SignWithAccount(account base.Account) (signedTxn *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

func (t *Transaction) SignedTransactionWithAccount(account base.Account) (signedTx base.SignedTransaction, err error) {
	starknetAccount := AsStarknetAccount(account)
	if starknetAccount == nil {
		return nil, base.ErrInvalidAccountType
	}
	return &SignedTransaction{
		Account:   starknetAccount,
		invokeTxn: t,
	}, nil
}
