package aptos

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/coming-chat/go-aptos/graphql"
	"github.com/coming-chat/go-aptos/nft"
	txnBuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const (
	GraphUrlMainnet = graphql.GraphUrlMainnet
	GraphUrlTestnet = graphql.GraphUrlTestnet
)

type NFTFetcher struct {
	Chain    *Chain
	client   *nft.TokenClient
	GraphUrl string
}

// Deprecated: use `NewNFTFetcherGraphUrl()`
func NewNFTFetcher(chain *Chain) *NFTFetcher {
	return &NFTFetcher{
		Chain: chain,
	}
}

func NewNFTFetcherGraphUrl(url string) *NFTFetcher {
	if url == "" {
		url = GraphUrlMainnet
	}
	return &NFTFetcher{
		GraphUrl: url,
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
	if f.GraphUrl != "" {
		return f.fetchNFTsUseGraphql(owner)
	}

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

func transformGraphToken(token nft.GraphQLToken) *base.NFT {
	t, err := time.Parse(`2006-01-02T15:04:05`, token.LastTransactionTimestamp)
	if err != nil {
		t = time.Time{}
	}
	nft := base.NFT{
		HashString: strconv.FormatUint(token.LastTransactionVersion, 10),
		Timestamp:  t.Unix(),

		Id:              "",
		Name:            token.Name,
		Image:           strings.Replace(token.CurrentTokenData.MetadataUri, "ipfs://", "https://ipfs.io/ipfs/", 1),
		Standard:        "",
		Collection:      token.CollectionName,
		Description:     token.CurrentTokenData.Description,
		ContractAddress: token.CreatorAddress,
		RelatedUrl:      "",
	}
	return &nft
}

func (f *NFTFetcher) fetchNFTsUseGraphql(owner string) (map[string][]*base.NFT, error) {
	tokens, err := nft.FetchGraphqlTokensOfOwner(owner, f.GraphUrl)
	if err != nil {
		return nil, err
	}

	nftGroupd := make(map[string][]*base.NFT)
	for _, token := range tokens {
		nft := transformGraphToken(token)
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
