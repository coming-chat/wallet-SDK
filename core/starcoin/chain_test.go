package starcoin

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	BarnardRpcUrl = "https://barnard-seed.starcoin.org"
	MainRpcUrl    = "https://main-seed.starcoin.org"
)

func BarnardChain(t *testing.T) *Chain {
	return NewChainWithRpc(BarnardRpcUrl)
}

func TestTransactionDetail(t *testing.T) {
	chain := BarnardChain(t)

	hash := "0x5a16ec93f46d137be51d93df374426cb212b2bb04dbe476cbb8a6b72385291fc"
	// hash := "0x098045bd1f817f37f7759c7f0e1dc4a6d6e0a5f33a8e0748cd8928653ed6a31a"

	detail, err := chain.FetchTransactionDetail(hash)
	require.Nil(t, err)

	t.Log(detail)
}
