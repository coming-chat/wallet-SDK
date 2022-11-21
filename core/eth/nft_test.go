package eth

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/stretchr/testify/require"
)

func TestRSS3Notes(t *testing.T) {
	owner := "0x8c951f58F63C0018BFBb47A29e55e84507eD63Bd"
	// owner := "0xFCC3299Eb11790d36836F1A9aBDbE3D2435794C1"
	rss3Fetcher := RSS3Fetcher{
		Owner:   owner,
		Limit:   100,
		Network: NetworkNameEthereum,
	}
	run := false

	run = true
	if run {
		res, err := rss3Fetcher.FetchNotes("")
		require.Nil(t, err)
		t.Log(res)
		run = false
	}

	run = true
	if run {
		var nftFetcher base.NFTFetcher = &rss3Fetcher
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
		var nftFetcher base.NFTFetcher = &rss3Fetcher
		jsonString, err := nftFetcher.FetchNFTsJsonString(owner)
		require.Nil(t, err)
		t.Log(jsonString)
		run = false
	}
}
