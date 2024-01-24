package eth

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

const (
	BlockScoutURLBevm = "https://scan-canary.bevm.io/api/v2"
	BlockScoutURLEth  = "https://eth.blockscout.com/api/v2"
)

type BlockScout struct {
	BaseUrl string
}

func NewBlockScout(url string) *BlockScout {
	return &BlockScout{
		BaseUrl: url,
	}
}

// Nfts
// - params nextPage: set nil when query first page.
func (a *BlockScout) Nft(owner string, nextPage *BKSPageParams) (page *BKSNFTPage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	url := fmt.Sprintf("%v/addresses/%v/nft", a.BaseUrl, owner)
	if nextPage != nil {
		url = url + "?" + nextPage.String()
	}

	data, err := httpUtil.Get(url, nil)
	if err != nil {
		return
	}
	var rawPage bksRawItemsPage[*BKSNFT]
	err = json.Unmarshal(data, &rawPage)
	if err != nil {
		return
	}
	return &BKSNFTPage{rawPage}, nil
}

// MARK - types

type BKSNFTPage struct {
	bksRawItemsPage[*BKSNFT]
}

func (page *BKSNFTPage) ToNFTGroupedMap() *base.NFTGroupedMap {
	group := make(map[string]*base.NFTArray)
	for _, item := range page.Items {
		nft := item.ToBaseNFT()
		if arr, ok := group[nft.GroupName()]; ok {
			arr.Append(nft)
		} else {
			arr := base.NFTArray{AnyArray: []*base.NFT{nft}}
			group[nft.GroupName()] = &arr
		}
	}
	return &base.NFTGroupedMap{AnyMap: group}
}

type BKSNFT struct {
	// AnimationUrl   string `json:"animation_url"`
	// ExternalAppUrl string `json:"external_app_url"`
	// IsUnique       string `json:"is_unique"`
	// Value          string `json:"value"`
	Id        string          `json:"id"`
	ImageUrl  string          `json:"image_url"`
	Metadata  *BKSNFTMetadata `json:"metadata"`
	Owner     string          `json:"owner"`
	Token     *BKSToken       `json:"token"`
	TokenType string          `json:"token_type"`
}

func (n *BKSNFT) ToBaseNFT() *base.NFT {
	nft := &base.NFT{
		Id:       n.Id,
		Image:    n.ImageUrl,
		Standard: n.TokenType,
	}
	if n.Metadata != nil {
		nft.Name = n.Metadata.Name
		nft.Descr = n.Metadata.Descr
		nft.RelatedUrl = n.Metadata.ExternalUrl
	}
	if n.Token != nil {
		nft.Collection = n.Token.Name
		nft.ContractAddress = n.Token.Address
	}
	return nft
}

type BKSNFTMetadata struct {
	// Attributes any    `json:"attributes"`
	Descr       string `json:"description"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	ExternalUrl string `json:"external_url"`
}

type BKSToken struct {
	// CirculatingMarketCap string `json:"circulating_market_cap"`
	// Decimals             string `json:"decimals"`
	// ExchangeRate         string `json:"exchange_rate"`
	// Holders              string `json:"holders"`
	// IconUrl              string `json:"icon_url"`
	Address     string `json:"address"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	TotalSupply string `json:"total_supply"`
	Type        string `json:"type"`
}

// BlockScout Raw Items Page
type bksRawItemsPage[T any] struct {
	Items     []T            `json:"items"`
	NextPage_ *BKSPageParams `json:"next_page_params"`
}

func (p *bksRawItemsPage[T]) Count() int {
	return len(p.Items)
}

func (p *bksRawItemsPage[T]) ValueAt(index int) T {
	return p.Items[index]
}

func (p *bksRawItemsPage[T]) NextPageParams() *BKSPageParams {
	return p.NextPage_
}

func (p *bksRawItemsPage[T]) HasNextPage() bool {
	return p.NextPage_ != nil && p.NextPage_.Raw != nil
}

func (p *bksRawItemsPage[T]) JsonString() (*base.OptionalString, error) {
	return base.JsonString(p)
}

// BlockScout Next Page Params
type BKSPageParams struct {
	Raw map[string]interface{}
}

func (p BKSPageParams) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Raw)
}

func (p *BKSPageParams) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &p.Raw)
}

func (p *BKSPageParams) String() string {
	s := ""
	for k, v := range p.Raw {
		s = fmt.Sprintf("%v&%v=%v", s, k, v)
	}
	if strings.HasPrefix(s, "&") {
		return s[1:]
	} else {
		return s
	}
}
