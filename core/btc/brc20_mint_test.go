package btc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMintTransaction(t *testing.T) {
	user := "tb1pdq423fm5dv00sl2uckmcve8y3w7guev8ka6qfweljlu23mmsw63qk6w2v3"
	tick := "ovvo"
	amount := "1000"

	chain, err := NewChainWithChainnet(ChainTestnet)
	require.Nil(t, err)
	txn, err := chain.BuildBrc20MintTransaction(user, user, tick, amount, 2)
	require.Nil(t, err)

	t.Log(txn)
}
