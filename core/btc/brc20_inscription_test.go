package btc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchBrc20Inscription(t *testing.T) {
	owner := "bc1p6wtl8gp8s6k94j3ryuqdeq4dgdcm4yyyc265g3rlh2x9m4cqn32scpts08"

	page, err := FetchBrc20Inscription(owner, "1", 20)
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
