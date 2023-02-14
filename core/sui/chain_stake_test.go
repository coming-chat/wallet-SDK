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
