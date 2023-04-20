package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildMergeCoinRequest(t *testing.T) {
	chain := TestnetChain()

	owner := "0x9bab5b2fa325fe2b103fd6a56a93bf91925b269a2dd31ee146b693e5cb9d2901"
	amountStr := SUI(400).String()

	req, err := chain.BuildMergeCoinRequest(owner, "", amountStr)
	require.Nil(t, err)

	preview, err := chain.BuildMergeCoinPreview(req)
	require.Nil(t, err)

	t.Logf("preview success %v", preview.SimulateSuccess)
	t.Logf("preview achieved %v", preview.WillBeAchieved)
}

func TestBuildSplitCoinTransaction(t *testing.T) {
	chain := TestnetChain()

	owner := "0x9bab5b2fa325fe2b103fd6a56a93bf91925b269a2dd31ee146b693e5cb9d2901"
	amountStr := SUI(1).String()

	txn, err := chain.BuildSplitCoinTransaction(owner, "", amountStr)
	require.Nil(t, err)

	simulateCheck(t, chain, &txn.Txn, true)
}
