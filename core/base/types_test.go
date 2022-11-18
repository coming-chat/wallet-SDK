package base

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringArray(t *testing.T) {
	arr := StringArray{}

	arr.Append("AAA")
	require.Equal(t, arr.Count(), 1)

	arr.Append("bbb")
	require.Equal(t, arr.ValueOf(0), "AAA")
	require.Equal(t, arr.ValueOf(1), "bbb")

	arr.SetValue("ccc", 1)
	require.Equal(t, arr.ValueOf(1), "ccc")
	require.Equal(t, arr.String(), `["AAA","ccc"]`)
	require.Equal(t, arr.Count(), 2)

	arr.Append("ddd")
	arr.Remove(0)
	require.Equal(t, arr.Count(), 2)
	require.Equal(t, arr.String(), `["ccc","ddd"]`)
}

func TestExtractNFTImageUrl(t *testing.T) {
	url := "https://cdn-2.galxe.com/galaxy/images/alienswap/1667153514800858058.gif"

	r, err := ExtractNFTImageUrl(url)
	require.Nil(t, err)
	t.Log(r)
}
