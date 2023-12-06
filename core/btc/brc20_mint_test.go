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
	t.Log(list.Items)
}

func TestNewBrc20MintTransactionWithJsonString(t *testing.T) {
	jsonStr := `{
		"commit": "abc",
		"reveal": ["11", "22"]
	}`
	txn, err := NewBrc20MintTransactionWithJsonString(jsonStr)
	require.Nil(t, err)
	require.Equal(t, txn.Commit, "abc")
	require.Equal(t, txn.Reveal.ValueAt(1), "22")
}

func Test_Marshal_Brc20MintTransaction(t *testing.T) {
	jsonStr := `
	{
		"inscription": [
			"9ea4ca939fb364944a1a51ca42aa5a8e99c06b2d0e7dce45bf317c601ebd4d21i0"
		],
		"commit": "70736274ff010",
		"commit_custom": [
			"0118b014000000",
			"ed4fa9e028aa59e5384b219562034b8254347d591c3f667855616e8f33ab561e",
			"2"
		],
		"reveal": [
			"0100000000010"
		],
		"service_fee": 0,
		"satpoint_fee": 546,
		"network_fee": 70716,
		"commit_vsize": 154,
		"commit_fee": 37422
	}
	`

	txn, err := NewBrc20MintTransactionWithJsonString(jsonStr)
	require.Nil(t, err)
	require.Equal(t, txn.CommitCustom.BaseTx, "0118b014000000")
	require.Equal(t, txn.CommitCustom.Utxos.Count(), 1)
	require.Equal(t, txn.CommitCustom.Utxos.ValueAt(0).Index, int64(2))
}
