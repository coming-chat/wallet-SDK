package sui

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (c *Chain) FetchNFTs(owner string) (res map[string][]*base.NFT, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	address, err := sui_types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	client, err := c.Client()
	if err != nil {
		return
	}
	nftObjects, err := client.BatchGetFilteredObjectsOwnedByAddress(context.Background(), *address, types.SuiObjectDataOptions{
		ShowType:                true,
		ShowDisplay:             true,
		ShowPreviousTransaction: true,
	}, func(sod *types.SuiObjectData) bool {
		isCoin := strings.HasPrefix(*sod.Type, "0x2::coin::Coin<")
		return !isCoin
	})
	if err != nil {
		return
	}
	nfts := []*base.NFT{}
	for _, obj := range nftObjects {
		nft := TransformNFT(&obj)
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

func TransformNFT(nft *types.SuiObjectResponse) *base.NFT {
	if nft == nil || nft.Data == nil || nft.Data.Display == nil {
		return nil
	}
	var nftStruct struct {
		Data struct {
			ImageUrl    string `json:"image_url"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Creator     string `json:"creator"`
			// ProjectUrl  string `json:"project_url"`
			Link string `json:"link"`
		} `json:"data"`
	}
	metaBytes, err := json.Marshal(nft.Data.Display)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(metaBytes, &nftStruct)
	if err != nil {
		return nil
	}
	if nftStruct.Data.ImageUrl == "" {
		return nil
	}

	contractAddress := ""
	if nft.Data.Type != nil {
		typ, err := types.NewResourceType(*nft.Data.Type)
		if err == nil {
			contractAddress = typ.Address.String()
		}
		err = nil
	}
	hash := ""
	if nft.Data.PreviousTransaction != nil {
		hash = nft.Data.PreviousTransaction.String()
	}
	name := nftStruct.Data.Name
	if name == "" {
		name = nft.Data.ObjectId.String()
	}
	return &base.NFT{
		HashString:      hash,
		ContractAddress: contractAddress,

		Name:        name,
		Id:          nft.Data.ObjectId.String(),
		Description: nftStruct.Data.Description,
		RelatedUrl:  nftStruct.Data.Link,
		Image:       strings.Replace(nftStruct.Data.ImageUrl, "ipfs://", "https://ipfs.io/ipfs/", 1),
	}
}

func (c *Chain) MintNFT(creator, name, description, uri string) (txn *Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}

// Just encapsulation and callbacks to method `TransferObject`.
// @param gasId gas object to be used in this transaction, the gateway will pick one from the signer's possession if not provided
func (c *Chain) TransferNFT(sender, receiver, nftId string) (txn *Transaction, err error) {
	return c.TransferObject(sender, receiver, nftId, MaxGasBudget)
}
