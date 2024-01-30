package eth

import (
	"encoding/json"
	"fmt"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/base/inter"
	"github.com/coming-chat/wallet-SDK/graphql"
)

type donutGraphResp struct {
	Data struct {
		Type  string `json:"@type"`
		Value string `json:"value"`
	} `json:"data,omitempty"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func donutParser(resp donutGraphResp, out any) error {
	if resp.Code != 0 {
		return fmt.Errorf("error code %v: %v", resp.Code, resp.Msg)
	}
	return json.Unmarshal([]byte(resp.Data.Value), out)
}

type DonutInscription struct {
	Tick   string `json:"tick"`
	Amount string `json:"amount"`
}

type DonutInscriptionArray struct {
	inter.AnyArray[*DonutInscription]
}

// FetchDonutInscriptions
// - param graphURL Default "https://bc.dnt.social/v1/common/search"
func FetchDonutInscriptions(owner string, graphURL string) (arr *DonutInscriptionArray, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if graphURL == "" {
		graphURL = "https://bc.dnt.social/v1/common/search"
	}
	query := fmt.Sprintf(`{
		src20Balances(holder: "%v", first: 100) {
			edges{
				node{
					tick
					amount
				}
			}
		}
	}`, owner)

	var out struct {
		Src20Balances struct {
			Edges []struct {
				Node *DonutInscription `json:"node"`
			} `json:"edges"`
		} `json:"src20Balances"`
	}
	err = graphql.QueryString(query, graphURL, donutParser, &out)
	if err != nil {
		return
	}
	inscriptions := make([]*DonutInscription, len(out.Src20Balances.Edges))
	for idx, node := range out.Src20Balances.Edges {
		inscriptions[idx] = node.Node
	}
	return &DonutInscriptionArray{AnyArray: inscriptions}, nil
}
