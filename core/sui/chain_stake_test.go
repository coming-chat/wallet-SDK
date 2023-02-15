package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetValidatorState(t *testing.T) {
	chain := TestnetChain()

	state, err := chain.GetValidatorState()
	require.Nil(t, err)
	for _, v := range state.Validators {
		t.Logf("%-10v, APY = %v", v.Name, v.APY)
	}
}

func TestGetDelegatedStakes(t *testing.T) {
	chain := TestnetChain()
	acc := M1Account(t)

	list, err := chain.GetDelegatedStakes(acc.Address())
	require.Nil(t, err)
	for _, v := range list {
		t.Log(v)
	}
}

func TestAddDelegation(t *testing.T) {
	chain := TestnetChain()
	acc := M1Account(t)

	amount := "500000000" // 0.5
	validator := "0x0018bb48352b63c246bdb154b15a3b0d17dff193"
	txn, err := chain.AddDelegation(acc.Address(), amount, validator)
	require.Nil(t, err)

	signedTxn, err := txn.SignWithAccount(acc)
	require.Nil(t, err)

	if false {
		hash, err := chain.SendRawTransaction(signedTxn.Value)
		require.Nil(t, err)

		t.Log("add stake delegation succeed:", hash)
	}
}

func TestWithdrawDelegation(t *testing.T) {
	chain := TestnetChain()
	acc := M1Account(t)

	if false {
		delegationId := "0xd1e5f57aa2eb1ef7481e715c46c72fdfb46ec048"
		stakeId := "0x5cdb23dacf54329660467b900a2598bb796353fa"
		txn, err := chain.WithdrawDelegation(acc.Address(), delegationId, stakeId)
		require.Nil(t, err)

		signedTxn, err := txn.SignWithAccount(acc)
		require.Nil(t, err)

		hash, err := chain.SendRawTransaction(signedTxn.Value)
		require.Nil(t, err)

		t.Log("withdraw stake delegation succeed:", hash)
	}
}
