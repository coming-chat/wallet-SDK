package solana

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToken_BuildTransfer_SignedTransaction(t *testing.T) {
	account := M1Account(t)
	chain := TestnetChain()
	token := chain.MainToken()

	balance, err := token.BalanceOfAddress(account.Address())
	require.Nil(t, err)
	t.Log(balance.Total)

	txn, err := token.BuildTransfer(account.Address(), account.Address(), "100000")
	require.Nil(t, err)
	signedTxn, err := txn.SignedTransactionWithAccount(account)
	require.Nil(t, err)

	if false {
		hash, err := chain.SendSignedTransaction(signedTxn)
		require.Nil(t, err)
		t.Log("Transaction hash = ", hash.Value)
	}
}

func TestToken_BuildTransfer_GasFee(t *testing.T) {
	account := M1Account(t)
	chain := TestnetChain()
	token := chain.MainToken()

	txn, err := token.BuildTransfer(account.Address(), account.Address(), "10000")
	require.Nil(t, err)
	gas, err := chain.EstimateTransactionFee(txn)
	require.Nil(t, err)
	t.Log("Estimate gas fee = ", gas.Value)
}
