package btc

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestBrc20MintTransaction_SignAndSend(t *testing.T) {
	account, err := NewAccountWithMnemonic(testcase.M1, ChainTestnet)
	require.Nil(t, err)
	tick := "ovvo"
	amount := "1000"

	user, err := account.TaprootAddress()
	require.Nil(t, err)
	chain, err := NewChainWithChainnet(ChainTestnet)
	require.Nil(t, err)
	txn, err := chain.BuildBrc20MintTransaction(user.Value, user.Value, "mint", tick, amount, 1)
	require.Nil(t, err)

	signedTxn, err := txn.SignedTransactionWithAccount(account)
	require.Nil(t, err)

	if false {
		hash, err := chain.SendSignedTransaction(signedTxn)
		require.Nil(t, err)
		t.Log("brc20 mint success:", hash.Value)
	}
}
