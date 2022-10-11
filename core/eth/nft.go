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

	TagCollectible = "collectible"

	rss3Host = "https://pregod.rss3.dev/v1"
)

type RSS3Metadata struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Image           string `json:"image"`
	Standard        string `json:"standard"`
	Collection      string `json:"collection"`
	Description     string `json:"description"`
	ContractAddress string `json:"contract_address"`
}

type RSS3NoteAction struct {
	From string `json:"address_from"`
	To   string `json:"address_to"`
	Tag  string `json:"tag"`
	Type string `json:"type"`

	Metadata    *RSS3Metadata `json:"metadata"`
	RelatedUrls []string      `json:"related_urls"`

	Timestamp int64
	Hash      string
}

type RSS3Note struct {
	Timestamp time.Time         `json:"timestamp"`
	Hash      string            `json:"hash"`
	Success   bool              `json:"success"`
	Network   string            `json:"network"`
	Actions   []*RSS3NoteAction `json:"actions"`
}

func (a *RSS3NoteAction) IsNftAction() bool {
	return a.Tag == TagCollectible
}

// @return nft identifierKey if the action is a nft action, else return empty
func (a *RSS3NoteAction) NftIdentifierKey() string {
	if a.Tag == TagCollectible {
		return strings.ToLower(a.Metadata.Standard + a.Metadata.ContractAddress + a.Metadata.Id)
	}
	return ""
}

func (a *RSS3NoteAction) RelatedScanUrl() string {
	if len(a.RelatedUrls) == 0 {
		return ""
	}
	scanComponent := "/tx/" + a.Hash
	for _, url := range a.RelatedUrls {
		if strings.Contains(url, scanComponent) {
			return url
		}
	}
	return a.RelatedUrls[0]
}

func (a *RSS3NoteAction) Nft() *base.NFT {
	if a.Tag != TagCollectible {
		return nil
	}

	n := &base.NFT{}
	n.Id = a.Metadata.Id
	n.Name = a.Metadata.Name
	n.Image = strings.Replace(a.Metadata.Image, "ipfs://", "https://ipfs.io/ipfs/", 1)
	n.Standard = a.Metadata.Standard
	n.Collection = a.Metadata.Collection
	n.Description = a.Metadata.Description
	n.ContractAddress = a.Metadata.ContractAddress
	n.RelatedUrl = a.RelatedScanUrl()
	n.Timestamp = a.Timestamp
	n.HashString = a.Hash
	return n
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
	url := fmt.Sprintf("%v/notes/%v?network=%v&limit=%v&tag=%v", rss3Host, f.Owner, f.Network, f.Limit, TagCollectible)
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

func (f *RSS3Fetcher) FetchNFTs(owner string) (map[string][]*base.NFT, error) {
	f.Limit = 500
	f.NextCursor = ""
	f.Owner = strings.ToLower(owner)

	actions := make(map[string]*RSS3NoteAction)
	willTradeInFutureActions := []*RSS3NoteAction{}
	for true {
		notes, err := f.FetchNotesNext()
		if err != nil {
			return nil, err
		}

		for i := len(notes) - 1; i >= 0; i-- {
			note := notes[i]
			if !note.Success {
				continue
			}
			for _, action := range note.Actions {
				nftKey := action.NftIdentifierKey()
				if nftKey == "" {
					continue
				}
				if action.To == f.Owner { // receive a nft
					actions[nftKey] = action
					action.Timestamp = note.Timestamp.Unix()
					action.Hash = note.Hash
				} else if action.From == f.Owner { // send a nft
					_, exits := actions[nftKey]
					if exits {
						delete(actions, nftKey)
					} else {
						willTradeInFutureActions = append(willTradeInFutureActions, action)
					}
				} else {
					println("Invalid data, not in or out", action)
				}
			}
		}

		for i := len(willTradeInFutureActions) - 1; i >= 0; i-- {
			action := willTradeInFutureActions[i]
			_, exits := actions[action.NftIdentifierKey()]
			if exits {
				delete(actions, action.NftIdentifierKey())
				willTradeInFutureActions = append(willTradeInFutureActions[:i], willTradeInFutureActions[i+1:]...)
			}
		}

		if len(f.NextCursor) > 0 {
			continue
		} else {
			break
		}
	}

	if len(willTradeInFutureActions) != 0 {
		println("Invalid status that trade nft have not clean", willTradeInFutureActions)
	}

	nftGroupd := make(map[string][]*base.NFT)
	for _, action := range actions {
		if nft := action.Nft(); nft != nil {
			key := nft.GroupName()
			group, exist := nftGroupd[key]
			if exist {
				nftGroupd[key] = append(group, nft)
			} else {
				nftGroupd[key] = []*base.NFT{nft}
			}
		}
	}
	for _, group := range nftGroupd {
		sort.Slice(group, func(i, j int) bool {
			return group[i].Timestamp > group[j].Timestamp
		})
	}
	return nftGroupd, nil
}

// @return json string that grouped by nft's collection
func (f *RSS3Fetcher) FetchNFTsJsonString(owner string) (*base.OptionalString, error) {
	nfts, err := f.FetchNFTs(owner)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(nfts)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: string(bytes)}, nil
}
