package sui

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchNfts(t *testing.T) {
	owner := "0xbb8f7e72ae99d371020a1ccfe703bfb64a8a430f"
	chain := DevnetChain()

	nfts, err := chain.FetchNFTs(owner)
	require.Nil(t, err)
	for name, group := range nfts {
		t.Log("=======================================")
		t.Logf("group: %v, count: %v", name, len(group))
		for idx, nft := range group {
			t.Logf("%4v: %v", idx, nft)
		}
	}
}
