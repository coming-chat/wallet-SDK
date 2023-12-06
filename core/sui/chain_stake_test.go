package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	ComingChatValidatorAddress = "0x520289e77c838bae8501ae92b151b99a54407288fdd20dee6e5416bfe943eb7a"
	ComingChatValidatorMainnet = "0x11ec7353e9e08ed4ca13b935ad930a2f937112736aec5eedd08c95cc5cd4c264"
)

func TestGetValidatorState(t *testing.T) {
	chain := MainnetChain()

	state, err := chain.GetValidatorState()
	require.Nil(t, err)
	for idx, v := range state.Validators.AnyArray {
		t.Logf("%v, %-10v, APY = %v, totalStaked=%v", idx+1, v.Name, v.APY, v.TotalStaked)
	}
}

func TestStakeEarningTimems(t *testing.T) {
	state := ValidatorState{
		Epoch:                 9,
		EpochDurationMs:       86400000,
		EpochStartTimestampMs: 1681266455000,
	}

	ti := state.EarningAmountTimeAfterNowMs()
	t.Log(ti)

	delegated := DelegatedStake{
		RequestEpoch: 8,
	}
	ti2 := delegated.EarningAmountTimeAfterNowMs(&state)
	t.Log(ti2)
}

func TestGetDelegatedStakes(t *testing.T) {
	chain := TestnetChain()
	// address := M1Account(t).Address()
	address := "0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f"

	list, err := chain.GetDelegatedStakes(address)
	require.Nil(t, err)
	for _, v := range list.AnyArray {
		t.Log(v)
	}
}

func TestAddDelegation(t *testing.T) {
	chain := TestnetChain()
	acc := M1Account(t)

	amount := SUI(1).String()
	validator := ComingChatValidatorAddress
	txn, err := chain.AddDelegation(acc.Address(), amount, validator)
	require.Nil(t, err)

	gas, err := chain.EstimateTransactionFee(txn)
	require.Nil(t, err)
	t.Log(gas.Value)

	simulateTxnCheck(t, chain, txn, false)
}

func TestWithdrawDelegation(t *testing.T) {
	chain := TestnetChain()
	owner := "0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f"

	stakedArray, err := chain.GetDelegatedStakes(owner)
	require.Nil(t, err)
	require.Greater(t, stakedArray.Count(), 0)

	stake := stakedArray.ValueAt(0)
	stakeId := stake.StakeId
	txn, err := chain.WithdrawDelegation(owner, stakeId)
	require.Nil(t, err)

	simulateTxnCheck(t, chain, txn, false)
}
