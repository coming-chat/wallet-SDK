package aptos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchNFTs(t *testing.T) {
	chain := NewChainWithRestUrl(mainnetRestUrl)
	nftFetcher := NewNFTFetcher(chain)
	owner := "0x6ed6f83f1891e02c00c58bf8172e3311c982b1c4fbb1be2d85a55562d4085fb1"
	run := false

	// run = true
	if run {
		nfts, err := nftFetcher.FetchNFTs(owner)
		require.Nil(t, err)
		for name, group := range nfts {
			t.Log("=======================================")
			t.Logf("group: %v, count: %v", name, len(group))
			for idx, nft := range group {
				t.Logf("%4v: %v", idx, nft)
			}
		}
		run = false
	}

	run = true
	if run {
		jsonString, err := nftFetcher.FetchNFTsJsonString(owner)
		require.Nil(t, err)
		t.Log(jsonString)
		run = false
	}
}
