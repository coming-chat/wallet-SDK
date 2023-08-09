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
	println("address: ", user.Value)

	chain, err := NewChainWithChainnet(ChainTestnet)
	require.Nil(t, err)
	balance, err := chain.BalanceOfAddress(user.Value)
	require.Nil(t, err)
	println("balance: ", balance.Total)

	txn, err := chain.BuildBrc20MintTransaction(user.Value, user.Value, "mint", tick, amount, 1)
	require.Nil(t, err)

	signedTxn, err := txn.SignedTransactionWithAccount(account)
	require.Nil(t, err)

	if false {
		hash, err := chain.SendSignedTransaction(signedTxn)
		require.Nil(t, err)
		println("brc20 mint success: ", hash.Value)
	}
}

func TestFetchBrc20InscriptionList(t *testing.T) {
	account, err := NewAccountWithMnemonic(testcase.M1, ChainTestnet)
	require.Nil(t, err)
	chain, err := NewChainWithChainnet(ChainTestnet)
	require.Nil(t, err)
	addr, err := account.TaprootAddress()
	require.Nil(t, err)

	list, err := chain.FetchBrc20Inscription(addr.Value, "", 10)
	require.Nil(t, err)
	t.Log(list.ItemArray().Values...)
}

func TestNewBrc20MintTransactionWithJsonString(t *testing.T) {
	jsonStr := `{
		"commit": "abc",
		"reveal": ["11", "22"]
	}`
	txn, err := NewBrc20MintTransactionWithJsonString(jsonStr)
	require.Nil(t, err)
	require.Equal(t, txn.Commit, "abc")
	require.Equal(t, txn.Reveal[1], "22")
}
