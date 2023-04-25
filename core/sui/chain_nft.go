package sui

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (c *Chain) FetchNFTs(owner string) (res map[string][]*base.NFT, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	address, err := types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	client, err := c.Client()
	if err != nil {
		return
	}
	nftObjects, err := client.BatchGetFilteredObjectsOwnedByAddress(context.Background(), *address, types.SuiObjectDataOptions{
		ShowType:                true,
		ShowContent:             true,
		ShowPreviousTransaction: true,
	}, func(sod *types.SuiObjectData) bool {
		if strings.HasPrefix(*sod.Type, "0x2::coin::Coin<") {
			return false
		}
		return true
	})
	if err != nil {
		return
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

func transformNFT(nft *types.SuiObjectResponse) *base.NFT {
	if nft == nil || nft.Data == nil || nft.Data.Content == nil || nft.Data.Content.Data.MoveObject == nil {
		return nil
	}
	fields := struct {
		Id struct {
			Id string `json:"id"`
		} `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Url         string `json:"url"`
	}{}
	metaBytes, err := json.Marshal(nft.Data.Content.Data.MoveObject.Fields)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(metaBytes, &fields)
	if err != nil {
		return nil
	}
	if fields.Name == "" && fields.Url == "" {
		return nil
	}

	return &base.NFT{
		HashString: *nft.Data.PreviousTransaction,

		Id:          fields.Id.Id,
		Name:        fields.Name,
		Description: fields.Description,
		Image:       strings.Replace(fields.Url, "ipfs://", "https://ipfs.io/ipfs/", 1),
	}
}

func (c *Chain) MintNFT(creator, name, description, uri string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := types.NewAddressFromHex(creator)
	if err != nil {
		return nil, errors.New("Invalid creator address")
	}
	client, err := c.Client()
	if err != nil {
		return
	}
	return c.EstimateTransactionFeeAndRebuildTransaction(MaxGasBudget, func(gasBudget uint64) (*Transaction, error) {
		txBytes, err := client.MintNFT(context.Background(), *signer, name, description, uri, nil, gasBudget)
		if err != nil {
			return nil, err
		}
		return &Transaction{
			Txn: *txBytes,
		}, nil
	})
}

// Just encapsulation and callbacks to method `TransferObject`.
// @param gasId gas object to be used in this transaction, the gateway will pick one from the signer's possession if not provided
func (c *Chain) TransferNFT(sender, receiver, nftId string) (txn *Transaction, err error) {
	return c.TransferObject(sender, receiver, nftId, MaxGasBudget)
}
