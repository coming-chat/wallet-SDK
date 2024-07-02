package btc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type Brc20TransferTransaction struct {
	Transaction string `json:"transaction"`
	NetworkFee  int64  `json:"network_fee"`
	CommitId    string `json:"commit_id"`
	CommitFee   int64  `json:"commit_fee"`
	CommitVsize int64  `json:"commit_vsize"`

	CommitCustom *Brc20CommitCustom `json:"commit_custom"`
}

func NewBrc20TransferTransaction() *Brc20TransferTransaction {
	return &Brc20TransferTransaction{}
}

func NewBrc20TransferTransactionWithJsonString(jsonStr string) (*Brc20TransferTransaction, error) {
	data := []byte(jsonStr)
	var out Brc20TransferTransaction
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
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
	insIds := inscriptionIds.AnyArray[:]
	for idx, id := range insIds {
		insIds[idx] = strings.Replace(id, "i", ":", 1)
	}

	host, err := comingOrdHost(c.Chainnet)
	if err != nil {
		return
	}
	url := fmt.Sprintf("%v/TransferOutpoints", host)
	header := map[string]string{"Content-Type": "application/json"}
	requestBody := map[string]any{
		"commit_fee_rate": strconv.FormatInt(feeRate, 10),
		"source":          sender,
		"destination":     receiver,
		"outpoint":        insIds,
		"postage":         546,
		"op_return":       opReturn,
	}
	requestBytes, _ := json.Marshal(requestBody)
	resp, err := httpUtil.Request(http.MethodPost, url, header, requestBytes)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", resp.Code, string(resp.Body))
	}

	var r struct {
		Commit_id     string             `json:"commit_id"`
		Commit_psbt   string             `json:"commit_psbt"`
		Commit_fee    int64              `json:"commit_fee"`
		Commit_vsize  int64              `json:"commit_vsize"`
		Commit_custom *Brc20CommitCustom `json:"commit_custom"`
		Network_fee   int64              `json:"network_fee"`
	}
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, err
	}
	return &Brc20TransferTransaction{
		Transaction:  r.Commit_psbt,
		NetworkFee:   r.Network_fee,
		CommitFee:    r.Commit_fee,
		CommitVsize:  r.Commit_vsize,
		CommitId:     r.Commit_id,
		CommitCustom: r.Commit_custom,
	}, nil
}
