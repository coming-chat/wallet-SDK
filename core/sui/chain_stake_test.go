package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStakeList(t *testing.T) {
	chain := TestnetChain()

	list, err := chain.GetStakeList()
	require.Nil(t, err)
	for _, v := range list {
		t.Logf("%v, APY = %v", v.Name, v.APY)
	}
}
