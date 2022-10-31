package aptos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchNFTs(t *testing.T) {
	nftFetcher := NewNFTFetcherGraphUrl("")
	// owner := "0x6ed6f83f1891e02c00c58bf8172e3311c982b1c4fbb1be2d85a55562d4085fb1"
	owner := "0xf5bb1482c28e3c600edf4cac9a10511b3d9a8e162d5b64d9741e2a8cb086bb50"

	if true { /*
			if false { /**/
		nfts, err := nftFetcher.FetchNFTs(owner)
		require.Nil(t, err)
		for name, group := range nfts {
			t.Log("=======================================")
			t.Logf("group: %v, count: %v", name, len(group))
			for idx, nft := range group {
				t.Logf("%4v: %v", idx, nft)
			}
		}
	}

	// if true { /*
	if false { /**/
		jsonString, err := nftFetcher.FetchNFTsJsonString(owner)
		require.Nil(t, err)
		t.Log(jsonString)
	}
}
