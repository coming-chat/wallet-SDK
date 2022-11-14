package sui

import (
	"context"
	"encoding/json"
	"errors"
	"sort"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (c *Chain) FetchNFTs(owner string) (map[string][]*base.NFT, error) {
	address, err := types.NewAddressFromHex(owner)
	if err != nil {
		return nil, err
	}
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	nftObjects, err := client.GetNFTsOwnedByAddress(context.Background(), *address)
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

// @param gasId gas object to be used in this transaction, the gateway will pick one from the signer's possession if not provided
func (c *Chain) MintNFT(creator, name, description, uri, gasId string, gasBudget int64) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := types.NewAddressFromHex(creator)
	if err != nil {
		return nil, errors.New("Invalid creator address")
	}
	var gas *types.ObjectId = nil
	if gasId != "" {
		gas, err = types.NewHexData(gasId)
		if err != nil {
			return nil, errors.New("Invalid gas object id")
		}
	}
	client, err := c.client()
	if err != nil {
		return
	}
	tx, err := client.MintNFT(context.Background(), *signer, name, description, uri, gas, uint64(gasBudget))
	if err != nil {
		return
	}
	return &Transaction{Txn: *tx}, nil
}

// @param gasId gas object to be used in this transaction, the gateway will pick one from the signer's possession if not provided
func (c *Chain) TransferNFT(sender, receiver, nftId, gasId string, gasBudget int64) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	senderAddress, err := types.NewAddressFromHex(sender)
	if err != nil {
		return
	}
	receiverAddress, err := types.NewAddressFromHex(receiver)
	if err != nil {
		return
	}
	nftObject, err := types.NewHexData(nftId)
	if err != nil {
		return nil, err
	}
	var gas *types.ObjectId = nil
	if gasId != "" {
		gas, err = types.NewHexData(gasId)
		if err != nil {
			return nil, errors.New("Invalid gas object id")
		}
	}
	client, err := c.client()
	if err != nil {
		return
	}
	tx, err := client.TransferObject(context.Background(), *senderAddress, *receiverAddress, *nftObject, gas, uint64(gasBudget))
	if err != nil {
		return
	}
	return &Transaction{Txn: *tx}, nil
}
