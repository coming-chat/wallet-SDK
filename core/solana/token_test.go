package solana

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func SOL(amount float64) testcase.Amount {
	return testcase.Amount{Amount: amount, Multiple: 1e9}
}

func TestToken_BuildTransfer_SignedTransaction(t *testing.T) {
	account := M1Account(t)
	chain := TestnetChain()
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
