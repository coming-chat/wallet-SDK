package base

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
}

func (n *NFT) GroupName() string {
	if n.Collection == "" {
		return "Others"
	} else {
		return n.Collection
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
