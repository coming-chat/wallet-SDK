package aptos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchNFTs(t *testing.T) {
	nftFetcher := NewNFTFetcher("")
	// owner := "0x6ed6f83f1891e02c00c58bf8172e3311c982b1c4fbb1be2d85a55562d4085fb1"
	owner := "0xf5bb1482c28e3c600edf4cac9a10511b3d9a8e162d5b64d9741e2a8cb086bb50"
	creatorAddress := "0x93341710806dbac4fed2fd8251f6c2a49566c75aec00150fadc5db14c07f3d4c"
	if true { /*
			if false { /**/
		nfts, err := nftFetcher.FetchNFTs(owner)
		require.Nil(t, err)
		for name, group := range nfts {
			t.Log("=======================================")
			t.Logf("group: %s, count: %d", name, len(group))
			for idx, nft := range group {
				t.Logf("%4d: %+v", idx, nft)
			}
		}

		t.Log("\n<====================test filter============>\n")
		nftList, err := nftFetcher.FetchNFTsFilterByCreatorAddr(owner, creatorAddress)
		require.NoError(t, err)
		require.NotEmptyf(t, nftList, "fetch nft empty")
		collectName := nftList[0].Collection
		for _, v := range nftList {
			if v.Collection != collectName {
				t.Errorf("filter failed")
			}
			t.Logf("%+v", v)
		}
	}

	// if true { /*
	if false { /**/
		jsonString, err := nftFetcher.FetchNFTsJsonString(owner)
		require.Nil(t, err)
		t.Log(jsonString)
	}
}

func TestFetchNFTsTestnet(t *testing.T) {
	owner := "0xd77e3ea8aa559bc7d5a238314201f0dfd2643a0fdd7ee32d7139d8b2310b4001"

	nftFetcher := NewNFTFetcher(GraphUrlTestnet)
	nfts, err := nftFetcher.FetchNFTs(owner)
	require.Nil(t, err)
	t.Log(nfts)
}
