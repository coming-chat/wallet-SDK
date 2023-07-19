package starcoin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBalance(t *testing.T) {
	chain := BarnardChain(t)
	token := chain.MainToken()
	account := M1Account(t)

	balance, err := token.BalanceOfAddress(account.Address())
	require.Nil(t, err)

	t.Log(balance)
}

func TestTransfer(t *testing.T) {
	chain := BarnardChain(t)
	token := NewMainToken(chain)

	account1 := M1Account(t)
	receiver := account1
	amount := "10000"

	rawTx, err := token.BuildTransferTxWithAccount(account1, receiver.Address(), amount)
	require.Nil(t, err)

	hash, err := chain.SendRawTransaction(rawTx.Value)
	require.Nil(t, err)

	t.Log(hash)
}

func TestEstimateFee(t *testing.T) {
	chain := BarnardChain(t)
	token := NewMainToken(chain)
	account := M1Account(t)
	gasFee, err := token.EstimateFees(account, account.Address(), "100")
	require.Nil(t, err)
	t.Log("gasfee = ", gasFee.Value)
}

func TestToken_BuildTransfer_SignedTransaction(t *testing.T) {
	account := M1Account(t)
	chain := BarnardChain(t)
	token := chain.MainToken()

	balance, err := token.BalanceOfAddress(account.Address())
	require.Nil(t, err)
	t.Log("sender address = ", account.Address())
	t.Log("balance = ", balance.Usable)

	txn, err := token.BuildTransfer(account.Address(), account.Address(), "100")
	require.Nil(t, err)

	gasfee, err := chain.EstimateTransactionFeeUsePublicKey(txn, account.PublicKeyHex())
	require.Nil(t, err)
	t.Log("Estimate fee = ", gasfee.Value)

	signedTxn, err := txn.SignedTransactionWithAccount(account)
	require.Nil(t, err)

	if false {
		hash, err := chain.SendSignedTransaction(signedTxn)
		require.Nil(t, err)
		t.Log("Transaction hash = ", hash.Value)
	}
}
