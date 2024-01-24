package base

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base/inter"
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
	Descr           string `json:"descr"`
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

type NFTArray struct {
	inter.AnyArray[*NFT]
}

type NFTGroupedMap struct {
	inter.AnyMap[string, *NFTArray]
}

func (g *NFTGroupedMap) Keys() *StringArray {
	keys := inter.KeysOf(g.AnyMap)
	return &StringArray{keys}
}

func (g *NFTGroupedMap) ToNFTGroupArray() *NFTGroupArray {
	arr := make([]*NFTGroup, 0, len(g.AnyMap))
	for k, v := range g.AnyMap {
		collection := NFTGroup{
			Collection: k,
			Items:      v,
		}
		arr = append(arr, &collection)
	}
	return &NFTGroupArray{AnyArray: arr}
}

type NFTGroup struct {
	Collection string    `json:"collection"`
	Items      *NFTArray `json:"items"`
}

type NFTGroupArray struct {
	inter.AnyArray[*NFTGroup]
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
	url = strings.Replace(url, "ipfs://", "https://ipfs.io/ipfs/", 1)
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
