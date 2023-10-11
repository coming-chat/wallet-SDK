package solana

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const aTokenAddress = "38GJtbJkJJQKkR6CTZwwok4hAXoKnByU2MFAZjgTe114"

func TestSPLToken_GetBalance(t *testing.T) {
	chain := DevnetChain()
	token := NewSPLToken(chain, "38GJtbJkJJQKkR6CTZwwok4hAXoKnByU2MFAZjgTe114")

	balance, err := token.BalanceOfAddress("9B5XszUGdMaxCZ7uSQhPzdks5ZQSmWxrmzCSvtJ6Ns6g")
	require.Nil(t, err)
	t.Log(balance)
}

func TestSPLToken_CreateTokenAccount(t *testing.T) {
	signer := M1Account(t)
	owner := "abc"

	chain := DevnetChain()
	token := NewSPLToken(chain, aTokenAddress)

	txn, err := token.CreateTokenAccount(owner, signer.Address())
	require.Nil(t, err)

	fee, err := chain.EstimateTransactionFee(txn)
	require.Nil(t, err)
	t.Log(fee.Value)

	signedTxn, err := txn.SignedTransactionWithAccount(signer)
	require.Nil(t, err)

	if false {
		txhash, err := chain.SendSignedTransaction(signedTxn)
		require.Nil(t, err)
		t.Log("create token account for other user success, hash = ", txhash)
	}
}
