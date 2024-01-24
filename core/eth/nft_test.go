package eth

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
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

func TestBEVMjosjo(t *testing.T) {
	url := "https://eth.blockscout.com/api/v2/addresses/0x30b31174e5FEd4598aA0dF3191EbAA5AAb48d43E/tokens"
	datas, err := httpUtil.Get(url, map[string]string{
		"type":       "ERC-721",
		"fiat_value": "",
		"id":         "4674005429",
		"value":      "3",
		// "items_count": "550",
	})
	require.Nil(t, err)

	var resp struct {
		Items []struct {
			Token struct {
				Address string `json:"address"`
			} `json:"token"`
			Value string `json:"value"`
		} `json:"items"`
		Next_page_params struct {
			FiatValue  any    `json:"fiat_value"`
			Id         int64  `json:"id"`
			ItemsCount int    `json:"items_count"`
			Value      string `json:"value"`
		} `json:"next_page_params"`
	}

	err = json.Unmarshal(datas, &resp)
	require.Nil(t, err)

	total := int64(0)
	for _, item := range resp.Items {
		v, _ := strconv.ParseInt(item.Value, 10, 64)
		total += v
		println(item.Token.Address)
	}
	println("total:", total)
	println("fiat:", resp.Next_page_params.FiatValue)
	println("id:", resp.Next_page_params.Id)
	println("item_count:", resp.Next_page_params.ItemsCount)
	println("value:", resp.Next_page_params.Value)
}

func TestBKSNFTFetcher_FetchNextPage(t *testing.T) {
	owner := "0x30b31174e5FEd4598aA0dF3191EbAA5AAb48d43E" // nft very much
	// owner := "0x9dE416AB881b322eE0b4c189C2dE624090280cF2" // nft few
	fetcher := NewBKSNFTFetcher(BlockScoutURLEth, owner)

	page1, err := fetcher.FetchNextPage()
	require.Nil(t, err)
	showPage(page1)

	page2, err := fetcher.FetchNextPage()
	require.Nil(t, err)
	showPage(page2)

	fetcher.ResetPage()
	page1re, err := fetcher.FetchNextPage()
	require.Nil(t, err)
	showPage(page1re)
	require.True(t, reflect.DeepEqual(page1, page1re))
}
