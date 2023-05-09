package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransaction_SignWithAccount(t *testing.T) {
	account := M1Account(t)
	chain := TestnetChain()
	token := NewTokenMain(chain)

	toAddress := M1Account(t).Address()
	amount := SUI(1).String()
	// toAddress := account.Address()
	// amount := SUI(4).String() // test big amount transfer

	txn, err := token.BuildTransferTransaction(account, toAddress, amount)
	require.Nil(t, err)

	signedTx, err := txn.SignWithAccount(account)
	require.Nil(t, err)

	t.Log(signedTx.Value)
}
