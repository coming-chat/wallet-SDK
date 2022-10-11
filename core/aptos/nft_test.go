package aptos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchNFTs(t *testing.T) {
	chain := NewChainWithRestUrl(testnetRestUrl)
	nftFetcher := NewNFTFetcher(chain)
	owner := "0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5"
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
