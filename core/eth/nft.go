package eth

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

const (
	NetworkNameEthereum  = "ethereum"
	NetworkNameBSC       = "binance_smart_chain"
	NetworkNamePolygon   = "polygon"
	NetworkNameArbitrum  = "arbitrum"
	NetworkNameOptimism  = "optimism"
	NetworkNameAvalanche = "avalanche"

	rss3Host = "https://pregod.rss3.dev/v1.0.0"
)

type NftMetadata struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

type Nft struct {
	Timestamp  int64  `json:"timestamp"`
	HashString string `json:"hashString"`

	TokenId       string `json:"tokenId"`
	TokenAddress  string `json:"tokenAddress"`
	TokenStandard string `json:"tokenStandard"`

	// Metadata

	Name   string `json:"name"`
	Image  string `json:"image"`
	Detail string `json:"detail"`
}

type nft_temp struct {
	TokenId       string       `json:"token_id"`
	TokenAddress  string       `json:"token_address"`
	TokenStandard string       `json:"token_standard"`
	Metadata      *NftMetadata `json:"nft_metadata"`
}

type RSS3NoteAction struct {
	From string `json:"address_from"`
	To   string `json:"address_to"`
	Tag  string `json:"tag"`
	Type string `json:"type"`

	Metadata map[string]interface{}
}

type RSS3Note struct {
	Timestamp time.Time         `json:"timestamp"`
	Hash      string            `json:"hash"`
	Success   bool              `json:"success"`
	Actions   []*RSS3NoteAction `json:"actions"`
}

type RSS3Fetcher struct {
	// Filter notes by networks.
	// Possible values: [ethereum, ethereum_classic, binance_smart_chain, polygon, zksync, xdai, arweave, arbitrum, optimism, fantom, avalanche, crossbell]
	Network string
	// Limit the page number of Notes returned by the server. max 500.
	Limit int
	// owner eth address
	Owner string

	NextCursor string
	// PreCursor  string
}

func (n *Nft) IdentifierKey() string {
	return n.TokenStandard + n.TokenAddress + n.TokenId
}

func (a *RSS3NoteAction) Nft() *Nft {
	if a.Tag != "collectible" {
		return nil
	}

	bytes, err := json.Marshal(a.Metadata["token"])
	if err != nil {
		return nil
	}
	var nftTemp = nft_temp{}
	err = json.Unmarshal(bytes, &nftTemp)
	if err != nil {
		println("decode nft error", err)
		return nil
	}

	nft := Nft{}
	nft.TokenStandard = nftTemp.TokenStandard
	nft.TokenAddress = nftTemp.TokenAddress
	nft.TokenId = nftTemp.TokenId
	if m := nftTemp.Metadata; m != nil {
		nft.Name = m.Name
		nft.Image = strings.Replace(m.Image, "ipfs://", "https://ipfs.io/ipfs/", 1)
		nft.Detail = m.Description
	}

	return &nft
}

func NewRSS3FetcherWithNetwork(network string, owner string) *RSS3Fetcher {
	return &RSS3Fetcher{
		Network: network,
		Owner:   owner,
		Limit:   100,
	}
}

func (f *RSS3Fetcher) FetchNotes(cursor string) ([]RSS3Note, error) {
	if !IsValidAddress(f.Owner) {
		return nil, fmt.Errorf("Invalid owner address %v", f.Owner)
	}
	// https://pregod.rss3.dev/v1.0.0/notes/0x8c951f58F63C0018BFBb47A29e55e84507eD63Bd?tag=collectible
	url := fmt.Sprintf("%v/notes/%v?network=%v&limit=%v&tag=%v", rss3Host, f.Owner, f.Network, f.Limit, "collectible")
	if len(cursor) > 0 {
		url = url + "&cursor=" + cursor
	}
	body, err := httpUtil.Get(url, nil)
	if err != nil {
		return nil, err
	}

	var res = struct {
		Total  int        `json:"total"`
		Cursor string     `json:"cursor"`
		Result []RSS3Note `json:"result"`
	}{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	if len(res.Cursor) > 0 {
		f.NextCursor = res.Cursor
	}

	return res.Result, nil
}

func (f *RSS3Fetcher) FetchNotesNext() ([]RSS3Note, error) {
	return f.FetchNotes(f.NextCursor)
}

// func (f *RSS3Fetcher) FetchNotesPre() ([]RSS3Note, error) {
// 	return f.FetchNotes(f.PreCursor)
// }

func (f *RSS3Fetcher) FetchNtfs() ([]*Nft, error) {
	f.Limit = 500
	f.NextCursor = ""
	f.Owner = strings.ToLower(f.Owner)

	nfts := make(map[string]*Nft)
	willTradeInFutureNfts := []*Nft{}
	for true {
		notes, err := f.FetchNotesNext()
		if err != nil {
			return nil, err
		}

		for i := len(willTradeInFutureNfts) - 1; i >= 0; i-- {
			nft := willTradeInFutureNfts[i]
			_, exits := nfts[nft.IdentifierKey()]
			if exits {
				delete(nfts, nft.IdentifierKey())
				willTradeInFutureNfts = append(willTradeInFutureNfts[:i], willTradeInFutureNfts[i+1:]...)
			}
		}

		for i := len(notes) - 1; i >= 0; i-- {
			note := notes[i]
			for _, action := range note.Actions {
				nft := action.Nft()
				if nft == nil {
					continue
				}
				if action.To == f.Owner {
					nfts[nft.IdentifierKey()] = nft
					nft.Timestamp = note.Timestamp.Unix()
					nft.HashString = note.Hash
				} else if action.From == f.Owner {
					_, exits := nfts[nft.IdentifierKey()]
					if !exits {
						willTradeInFutureNfts = append(willTradeInFutureNfts, nft)
					} else {
						delete(nfts, nft.IdentifierKey())
					}
				} else {
					println("Invalid data, not in or out", action)
				}
			}
		}

		if len(f.NextCursor) > 0 {
			continue
		} else {
			break
		}
	}

	for i := len(willTradeInFutureNfts) - 1; i >= 0; i-- {
		nft := willTradeInFutureNfts[i]
		_, exits := nfts[nft.IdentifierKey()]
		if exits {
			delete(nfts, nft.IdentifierKey())
			willTradeInFutureNfts = append(willTradeInFutureNfts[:i], willTradeInFutureNfts[i+1:]...)
		}
	}

	if len(willTradeInFutureNfts) != 0 {
		println("Invalid status that trade nft have not clean", willTradeInFutureNfts)
	}

	nftList := []*Nft{}
	for _, nft := range nfts {
		nftList = append(nftList, nft)
	}
	sort.Slice(nftList, func(i, j int) bool {
		return nftList[i].Timestamp > nftList[j].Timestamp
	})
	return nftList, nil
}

func (f *RSS3Fetcher) FetchNftsJsonString() (*base.OptionalString, error) {
	nfts, err := f.FetchNtfs()
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(nfts)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: string(bytes)}, nil
}
