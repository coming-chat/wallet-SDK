package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetValidatorState(t *testing.T) {
	chain := DevnetChain()

	state, err := chain.GetValidatorState()
	require.Nil(t, err)
	for _, v := range state.Validators.Values {
		vv := v.(*Validator)
		t.Logf("%-10v, APY = %v", vv.Name, vv.APY)
	}
}

func TestGetDelegatedStakes(t *testing.T) {
	chain := DevnetChain()
	acc := M1Account(t)

	list, err := chain.GetDelegatedStakes(acc.Address())
	require.Nil(t, err)
	for _, v := range list.Values {
		vv := v.(*DelegatedStake)
		t.Log(vv)
	}
}

func TestAddDelegation(t *testing.T) {
	chain := DevnetChain()
	acc := M1Account(t)

	amount := "10000000" // 0.01
	validator := "0x0399e8864553720dac9258c7708ca821221bb246"
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
		stakeId := "0x5cdb23dacf54329660467b900a2598bb796353fa"
		txn, err := chain.WithdrawDelegation(acc.Address(), stakeId)
		require.Nil(t, err)

		signedTxn, err := txn.SignWithAccount(acc)
		require.Nil(t, err)

		hash, err := chain.SendRawTransaction(signedTxn.Value)
		require.Nil(t, err)

		t.Log("withdraw stake delegation succeed:", hash)
	}
}
