package btc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchBrc20Inscription(t *testing.T) {
	owner := "bc1pdq423fm5dv00sl2uckmcve8y3w7guev8ka6qfweljlu23mmsw63qpjc9k7"

	chain, err := NewChainWithChainnet(ChainMainnet)
	require.Nil(t, err)
	page, err := chain.FetchBrc20Inscription(owner, "0", 20)
	require.Nil(t, err)
	require.True(t, page.TotalCount() >= 1)
	t.Log(page.ItemArray().Values...)
	t.Log(page.ItemAt(0))

	jsonstring, err := page.JsonString()
	require.Nil(t, err)
	rePage, err := NewBrc20InscriptionPageWithJsonString(jsonstring.Value)
	require.Nil(t, err)
	require.Equal(t, page.TotalCount_, rePage.TotalCount_)
	require.Equal(t, page.Items[0], rePage.Items[0])
}

func TestFetchBrc20TransferableInscription(t *testing.T) {
	owner := "tb1p2hsjm57fsxrqcq5p42get87ttrw069kqa2ar444ma4ussquuaklqfsrknz"

	chain, err := NewChainWithChainnet(ChainTestnet)
	require.Nil(t, err)
	page, err := chain.FetchBrc20TransferableInscription(owner, "txtx")
	require.Nil(t, err)
	t.Log(page.ItemArray())
}
