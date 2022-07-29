package eth

import (
	"testing"
)

func TestRSS3Notes(t *testing.T) {
	fetcher := RSS3Fetcher{
		Owner: "0x8c951f58F63C0018BFBb47A29e55e84507eD63Bd",
		// Owner:   "0xFCC3299Eb11790d36836F1A9aBDbE3D2435794C1",
		Limit:   100,
		Network: NetworkNameEthereum,
	}

	// res, err := fetcher.FetchNotes("")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// println(res)

	// nfts, err := fetcher.FetchNtfs()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// for _, nft := range nfts {
	// 	fmt.Printf("%v %v %v: %v %v\n", nft.Timestamp, nft.TokenAddress, nft.TokenId, nft.Name, nft.Image)
	// }
	// t.Log("total ", len(nfts))

	json, err := fetcher.FetchNftsJsonString()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(json)
}
