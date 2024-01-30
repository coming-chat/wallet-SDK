package graphql

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type donutResp struct {
	Data struct {
		Type  string          `json:"@type"`
		Value json.RawMessage `json:"value"`
	} `json:"data,omitempty"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func donutParser(resp donutResp, out any) error {
	if resp.Code != 0 {
		return fmt.Errorf("error code %v: %v", resp.Code, resp.Msg)
	}
	return json.Unmarshal(resp.Data.Value, out)
}

func TestDonutGraphql(t *testing.T) {
	graphUrl := "https://bc.dnt.social/v1/common/search"
	holder := "0xa2cCF83EA437565a37E1F2d49940e0C4C7D7591e"
	query := fmt.Sprintf(`{
		src20Balances(holder: "%v", first: 100) {
			edges{
				node{
					tick
					amount
				}
			}
		}
	}`, holder)

	var resp interface{}
	err := QueryString(query, graphUrl, donutParser, &resp)
	require.NoError(t, err)
	t.Log(resp)

	errQuery := "aaaa"
	err = QueryString(errQuery, graphUrl, donutParser, &resp)
	require.Error(t, err)
}
