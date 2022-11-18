package base

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type NFT struct {
	Timestamp  int64  `json:"timestamp"`
	HashString string `json:"hashString"`

	Id              string `json:"id"`
	Name            string `json:"name"`
	Image           string `json:"image"`
	Standard        string `json:"standard"`
	Collection      string `json:"collection"`
	Description     string `json:"description"`
	ContractAddress string `json:"contract_address"`

	RelatedUrl string `json:"related_url"`

	// Aptos token's largest_property_version
	AptTokenVersion int64 `json:"apt_token_version"`
	// Aptos token's amount
	AptAmount int64 `json:"apt_amount"`
}

func (n *NFT) GroupName() string {
	if n.Collection == "" {
		return "Others"
	} else {
		return n.Collection
	}
}

func (n *NFT) ExtractedImageUrl() string {
	url, err := ExtractNFTImageUrl(n.Image)
	if err != nil {
		return n.Image
	} else {
		return url.Value
	}
}

type NFTFetcher interface {
	/** Gets all NFTs for the specified account and groups them by Collection name
	 * @owner The specified account address
	 * @return Grouped NFTs in below format
	 * ```
	 *  {
	 *    "Collection1": [ NFT1, NFT2 ],
	 *    "Collection2": [ NFT3 ],
	 *    "Collection3": [ NFT4, NFT5, NFT6, ... ],
	 *  }
	 * ```
	 */
	FetchNFTs(owner string) (map[string][]*NFT, error)

	/** Gets all NFT JSON Strings for the specified account
	 * @owner The specified account address
	 * @return This method directly calls `FetchNFTs()` and jsonifies the result and returns
	 */
	FetchNFTsJsonString(owner string) (*OptionalString, error)
}

// ExtractNFTImageUrl
// Extract the nft's real image url.
// If the content type of the given url is JSON, it's will return the `image` field specified url.
func ExtractNFTImageUrl(url string) (u *OptionalString, err error) {
	u = &OptionalString{Value: url}
	resp, err := httpUtil.Request(http.MethodHead, url, nil, nil) // HEAD request
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return u, nil
	}

	contentType := resp.Header["Content-Type"]
	if contentType == nil || len(contentType) == 0 {
		// no content type, return directly
		return u, nil
	}
	if !strings.Contains(contentType[0], "application/json") {
		// not json url, return directly
		return u, nil
	}

	// It should be return the original url if have anything error
	resp, err = httpUtil.Request(http.MethodGet, url, nil, nil) // GET request
	if resp.Code != 200 {
		return u, nil
	}
	jsonValue := struct {
		Image string `json:"image"`
	}{}
	err = json.Unmarshal(resp.Body, &jsonValue)
	if err != nil {
		return u, nil
	}
	return &OptionalString{Value: jsonValue.Image}, nil
}
