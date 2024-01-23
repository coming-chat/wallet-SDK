package eth

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/base/inter"
)

const (
	BlockScoutBevmUrl = "https://scan-canary.bevm.io/api/v2"
	BlockScoutEthUrl  = "https://eth.blockscout.com/api/v2"
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
// - params nextPage: set nil or empty when query first page.
func (a *BlockScout) Nfts(owner string, nextPage *BKSPageParams) (*BKSNFTPage, error) {
	return nil, nil
}

// MARK - types

type BKSNFTPage struct {
	raw bksRawItemsPage[*BKSNFT]
}

func (p *BKSNFTPage) NFTArray() *BKSNFTArray {
	return &BKSNFTArray{
		inter.AnyArray[*BKSNFT](p.raw.Items),
	}
}

func (p *BKSNFTPage) NextPageParams() *BKSPageParams {
	return p.raw.NextPageParams
}

func (p *BKSNFTPage) PageEnd() bool {
	return p.raw.PageEnd()
}

// AppendNextPage
// page.items = page.items + next.items
// page.nextPageItems = next.nextPageItems
func (page *BKSNFTPage) AppendNextPage(next *BKSNFTPage) {
	page.raw.Items = append(page.raw.Items, next.raw.Items...)
	page.raw.NextPageParams = next.raw.NextPageParams
}

type BKSNFTArray struct {
	inter.AnyArray[*BKSNFT]
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

func (n *BKSNFT) BaseNFT() *base.NFT {
	return &base.NFT{
		Id:              n.Id,
		Name:            n.Metadata.Name,
		Image:           n.ImageUrl,
		Standard:        n.TokenType,
		Collection:      n.Token.Name,
		Descr:           n.Metadata.Descr,
		ContractAddress: n.Token.Address,
		RelatedUrl:      n.Metadata.ExternalUrl,
	}
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
	Items          []T            `json:"items"`
	NextPageParams *BKSPageParams `json:"next_page_params"`
}

func (p *bksRawItemsPage[T]) PageEnd() bool {
	return p.NextPageParams == nil
}

// BlockScout Next Page Params
type BKSPageParams struct {
	Raw map[string]interface{}
}

func (p *BKSPageParams) Unmarshal(data []byte) error {
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
