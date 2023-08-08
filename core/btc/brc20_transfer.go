package btc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type Brc20TransferTransaction struct {
	Transaction string `json:"transaction"`
	NetworkFee  int64  `json:"network_fee"`
	// CommitCustom []string `json:"commit_custom"`
}

func (t *Brc20TransferTransaction) ToPsbtTransaction() (*PsbtTransaction, error) {
	packet, err := DecodePsbtTxToPacket(t.Transaction)
	if err != nil {
		return nil, err
	}
	return &PsbtTransaction{Packet: *packet}, nil
}

func (c *Chain) BuildBrc20TransferTransaction(
	sender, receiver string,
	inscriptionIds *base.StringArray,
	feeRate int64, opReturn string) (txn *Brc20TransferTransaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if inscriptionIds == nil || inscriptionIds.Count() <= 0 {
		return nil, errors.New("no inscriptions")
	}
	insIds := inscriptionIds.Values

	host, err := comingOrdHost(c.Chainnet)
	if err != nil {
		return
	}
	url := fmt.Sprintf("%v/transfer", host)
	header := map[string]string{"Content-Type": "application/json"}
	requestBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "transfer",
		"params": map[string]any{
			"source":            sender,
			"destination":       receiver,
			"fee_rate":          feeRate,
			"op_return":         opReturn,
			"outgoing":          insIds[0],
			"addition_outgoing": insIds[1:],
			"brc20_transfer":    true,
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
