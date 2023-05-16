package sui

import (
	"testing"

	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/stretchr/testify/require"
)

func TestBuildMergeCoinRequest(t *testing.T) {
	chain := TestnetChain()

	owner := M1Account(t).Address()
	amountStr := SUI(1).String()

	req, err := chain.BuildMergeCoinRequest(owner, "", amountStr)
	require.Nil(t, err)

	preview, err := chain.BuildMergeCoinPreview(req)
	require.Nil(t, err)

	t.Logf("preview success %v", preview.SimulateSuccess)
	t.Logf("preview achieved %v", preview.WillBeAchieved)
	simulateCheck(t, chain, &preview.Transaction.Txn, true)
}

func TestBuildSplitCoinTransaction(t *testing.T) {
	chain := TestnetChain()

	owner := M1Account(t).Address()
	amountStr := SUI(1).String()

	txn, err := chain.BuildSplitCoinTransaction(owner, "", amountStr)
	require.Nil(t, err)

	simulateCheck(t, chain, &txn.Txn, true)
}

func TestRunableSplitCoin(t *testing.T) {
	chain := TestnetChain()
	acc := M3Account(t)

	owner, err := sui_types.NewAddressFromHex(acc.Address())
	require.Nil(t, err)
	amount := SUI(1).String()

	txn, err := chain.BuildSplitCoinTransaction(owner.String(), "", amount)
	require.Nil(t, err)

	simulateCheck(t, chain, &txn.Txn, true)
	// executeTransaction(t, chain, &txn.Txn, acc.account)
}
