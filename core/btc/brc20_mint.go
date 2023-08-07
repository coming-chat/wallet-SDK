package btc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type Brc20MintTransaction struct {
	Commit      string   `json:"commit"`
	Reveal      []string `json:"reveal"`
	Inscription []string `json:"inscription"`

	NetworkFee  int64 `json:"network_fee"`
	SatpointFee int64 `json:"satpoint_fee"`
	ServiceFee  int64 `json:"service_fee"`
	CommitFee   int64 `json:"commit_fee"`
	CommitVsize int64 `json:"commit_vsize"`
}

func (c *Chain) BuildBrc20MintTransaction(sender, receiver string, ticker, amount string, feeRate int64) (txn *Brc20MintTransaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	url := ""
	switch c.Chainnet {
	case ChainMainnet, ChainBitcoin:
		url = "https://bitcoin.coming.chat/ord/mint"
	case ChainTestnet:
		url = "https://bitcoin.coming.chat/ord_testnet/mint"
	case ChainSignet:
		return nil, ErrUnsupportedChain
	}

	header := map[string]string{"Content-Type": "application/json"}
	requestBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "mint",
		"params": map[string]any{
			"source":      sender,
			"fee_rate":    feeRate,
			"content":     fmt.Sprintf(`{"p":"brc-20","op":"mint","tick":"%v","amt":"%v"}`, ticker, amount),
			"destination": receiver,
			"extension":   ".txt",
			"repeat":      1,
		},
	}
	requestBytes, _ := json.Marshal(requestBody)
	resp, err := httpUtil.Request(http.MethodPost, url, header, requestBytes)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", resp.Code, string(resp.Body))
	}

	err = json.Unmarshal(resp.Body, &txn)
	if err != nil {
		return nil, err
	}
	return txn, nil
}
