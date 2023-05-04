package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FetchSuiCatGlobalData(t *testing.T) {
	chain := MainnetChain()
	data, err := chain.FetchSuiCatGlobalData()
	require.Nil(t, err)
	t.Log(data.JsonString())
}

func Test_QueryIsInSuiCatWhiteList(t *testing.T) {
	chain := TestnetChain()

	res, err := chain.QueryIsInSuiCatWhiteList("0xf9ed7d8de1a6c44d703b64318a1cc687c324fdec35454281035a53ea3ba1a95a")
	require.Nil(t, err)
	t.Log(res.Value)

	res, err = chain.QueryIsInSuiCatWhiteList("0x09ed7d8de1a6c44d703b64318a1cc687c324fdec35454281035a53ea3ba1a95a")
	require.Nil(t, err)
	t.Log(res.Value)
}

func Test_MintSuiCatNFT(t *testing.T) {
	acc := M3Account(t)
	chain := TestnetChain()

	txn, err := chain.MintSuiCatNFT(acc.Address(), SUI(0.9).String())
	require.Nil(t, err)

	// simulateCheck(t, chain, &txn.Txn, true)
	executeTransaction(t, chain, &txn.Txn, acc.account)
}
