package graphql

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type Params struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName,omitempty"`
	Variables     map[string]any `json:"variables,omitempty"`
}

type Parser[RawResponse any] func(resp RawResponse, out any) error

func Query[RawResponse any](params Params, graphqlUrl string, parser Parser[RawResponse], out any) (err error) {
	if parser == nil || out == nil {
		return errors.New("the parser of the query cannot be nil")
	}
	body, err := json.Marshal(params)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, graphqlUrl, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header["Content-Type"] = []string{"application/json"}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var rawResp RawResponse
	err = json.Unmarshal(respBody, &rawResp)
	if err != nil {
		return
	}
	return parser(rawResp, out)
}

func QueryString[RawResponse any](query string, graphqlUrl string, parser Parser[RawResponse], out any) (err error) {
	return Query[RawResponse](Params{Query: query}, graphqlUrl, parser, out)
}
