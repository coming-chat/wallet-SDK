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
	receiver := M2Account(t)
	amount := "10000"

	rawTx, err := token.BuildTransferTxWithAccount(account1, receiver.Address(), amount)
	require.Nil(t, err)

	hash, err := chain.SendRawTransaction(rawTx.Value)
	require.Nil(t, err)

	t.Log(hash)

	detail, err := chain.FetchTransactionDetail(hash)
	t.Log(err)
	t.Log(detail)
}
