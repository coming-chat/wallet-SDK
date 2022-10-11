package aptos

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/coming-chat/go-aptos/nft"
	txnBuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type NFTFetcher struct {
	Chain  *Chain
	client *nft.TokenClient
}

func NewNFTFetcher(chain *Chain) *NFTFetcher {
	return &NFTFetcher{
		Chain: chain,
	}
}

func (f *NFTFetcher) tokenClient() (*nft.TokenClient, error) {
	if f.client == nil {
		restClient, err := f.Chain.client()
		if err != nil {
			return nil, err
		}
		f.client = nft.NewTokenClient(restClient)
	}
	return f.client, nil
}

func (f *NFTFetcher) FetchNFTs(owner string) (map[string][]*base.NFT, error) {
	account, err := txnBuilder.NewAccountAddressFromHex(owner)
	if err != nil {
		return nil, err
	}
	client, err := f.tokenClient()
	if err != nil {
		return nil, err
	}
	allTokens, err := client.GetAllTokenForAccount(*account)
	if err != nil {
		return nil, err
	}

	nftGroupd := make(map[string][]*base.NFT)
	for _, token := range allTokens {
		nft := transformNFT(token)
		key := nft.GroupName()
		group, exists := nftGroupd[key]
		if exists {
			nftGroupd[key] = append(group, nft)
		} else {
			nftGroupd[key] = []*base.NFT{nft}
		}
	}
	for _, group := range nftGroupd {
		sort.Slice(group, func(i, j int) bool {
			return group[i].Timestamp > group[j].Timestamp
		})
	}
	return nftGroupd, nil
}

func (f *NFTFetcher) FetchNFTsJsonString(owner string) (*base.OptionalString, error) {
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

func transformNFT(token *nft.NFTInfo) *base.NFT {
	nft := base.NFT{
		HashString: token.RelatedHash,
		Timestamp:  int64(token.RelatedTimestamp / 1000),

		Id:              "",
		Name:            token.TokenData.Name,
		Image:           strings.Replace(token.TokenData.Uri, "ipfs://", "https://ipfs.io/ipfs/", 1),
		Standard:        "",
		Collection:      token.TokenData.Collection,
		Description:     token.TokenData.Description,
		ContractAddress: token.TokenId.Creator,
		RelatedUrl:      "",
	}
	return &nft
}
