package sui

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

// Only support devnet now.
func (c *Chain) FetchNFTs(owner string) (map[string][]*base.NFT, error) {
	address, err := types.NewAddressFromHex(owner)
	if err != nil {
		return nil, err
	}
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	nftObjects, err := client.GetDevnetNFTOwnedByAddress(context.Background(), *address)
	if err != nil {
		return nil, err
	}
	nfts := []*base.NFT{}
	for _, obj := range nftObjects {
		nft := transformNFT(&obj)
		if nft != nil {
			nfts = append(nfts, nft)
		}
	}

	sort.Slice(nfts, func(i, j int) bool {
		return nfts[i].Timestamp > nfts[j].Timestamp
	})
	group := make(map[string][]*base.NFT)
	group["Other"] = nfts
	return group, nil
}

// Only support devnet now.
func (c *Chain) FetchNFTsJsonString(owner string) (*base.OptionalString, error) {
	nfts, err := c.FetchNFTs(owner)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(nfts)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: string(bytes)}, nil
}

func transformNFT(nft *types.ObjectRead) *base.NFT {
	if nft.Status != types.ObjectStatusExists {
		return nil
	}

	meta := struct {
		Fields struct {
			Id struct {
				Id string `json:"id"`
			} `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Url         string `json:"url"`
		} `json:"fields"`
	}{}
	metaBytes, err := json.Marshal(nft.Details.Data)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(metaBytes, &meta)
	if err != nil {
		return nil
	}

	return &base.NFT{
		HashString: nft.Details.PreviousTransaction.String(),

		Id:          meta.Fields.Id.Id,
		Name:        meta.Fields.Name,
		Description: meta.Fields.Description,
		Image:       meta.Fields.Url,
	}
}
