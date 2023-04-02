package sui

import (
	"strconv"
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
	address := M1Account(t).Address()
	// address := "0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f"

	list, err := chain.GetDelegatedStakes(address)
	require.Nil(t, err)
	for _, v := range list.Values {
		vv := v.(*DelegatedStake)
		t.Log(vv)
	}
}

func TestAddDelegation(t *testing.T) {
	chain := DevnetChain()
	acc := M1Account(t)

	amount := strconv.FormatInt(1e9, 10) // 1 SUI
	validator := "0x8ce890590fed55c37d44a043e781ad94254b413ee079a53fb5c037f7a6311304"
	txn, err := chain.AddDelegation(acc.Address(), amount, validator)
	require.Nil(t, err)

	signedTxn, err := txn.SignWithAccount(acc)
	require.Nil(t, err)

	if true {
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
