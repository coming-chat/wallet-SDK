package btc

import (
	"bytes"
	"encoding/hex"
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

	CommitCustom *Brc20CommitCustom `json:"commit_custom"`

	NetworkFee  int64 `json:"network_fee"`
	SatpointFee int64 `json:"satpoint_fee"`
	ServiceFee  int64 `json:"service_fee"`
	CommitFee   int64 `json:"commit_fee"`
	CommitVsize int64 `json:"commit_vsize"`

	signedTxn *SignedPsbtTransaction
}

func NewBrc20MintTransactionWithJsonString(jsonStr string) (*Brc20MintTransaction, error) {
	jsonBytes := []byte(jsonStr)
	var res Brc20MintTransaction
	err := json.Unmarshal(jsonBytes, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (t *Brc20MintTransaction) RevealArray() *base.StringArray {
	return &base.StringArray{Values: t.Reveal}
}
func (t *Brc20MintTransaction) InscriptionArray() *base.StringArray {
	return &base.StringArray{Values: t.Inscription}
}

func (t *Brc20MintTransaction) IsSigned() bool {
	return t.signedTxn != nil
}

func (t *Brc20MintTransaction) SignWithAccount(account base.Account) (signedTxn *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}
func (t *Brc20MintTransaction) SignedTransactionWithAccount(account base.Account) (signedTxn base.SignedTransaction, err error) {
	txn, err := NewPsbtTransaction(t.Commit)
	if err != nil {
		return
	}
	signedPsbtTxn, err := txn.SignedTransactionWithAccount(account)
	if err != nil {
		return
	}
	t.signedTxn = signedPsbtTxn.(*SignedPsbtTransaction)
	return t, nil
}

func (t *Brc20MintTransaction) HexString() (res *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

func (t *Brc20MintTransaction) PsbtHexString() (*base.OptionalString, error) {
	packet := t.signedTxn.Packet
	if err := EnsurePsbtFinalize(&packet); err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	if err := packet.Serialize(&buff); err != nil {
		return nil, err
	}
	hexString := hex.EncodeToString(buff.Bytes())
	return &base.OptionalString{Value: hexString}, nil
}

// PublishWithChain
// @return hash string
func (t *Brc20MintTransaction) PublishWithChain(c *Chain) (s *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	host, err := zeroWalletHost(c.Chainnet)
	if err != nil {
		return
	}
	hexString, err := t.PsbtHexString()
	if err != nil {
		return
	}

	url := fmt.Sprintf("%v/ordinal/tx/BitBox", host)
	header := map[string]string{"Content-Type": "application/json"}
	requestBody := map[string]any{
		"commit": hexString.Value,
		"reveal": t.Reveal,
	}
	requestBytes, _ := json.Marshal(requestBody)
	resp, err := httpUtil.Request(http.MethodPut, url, header, requestBytes)
	if err != nil {
		return
	}
	if resp.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", resp.Code, string(resp.Body))
	}

	var resObj struct {
		Hash string `json:"orderId"`
	}
	err = json.Unmarshal(resp.Body, &resObj)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: resObj.Hash}, nil
}

// BuildBrc20MintTransaction
// @param op "mint" or "transfer"
func (c *Chain) BuildBrc20MintTransaction(sender, receiver string, op, ticker, amount string, feeRate int64) (txn *Brc20MintTransaction, err error) {
	return c.BuildBrc20MintWithPostage(sender, receiver, op, ticker, amount, feeRate, 546)
}

// BuildBrc20MintWithPostage
// @param postage default is 546 if less than 546
func (c *Chain) BuildBrc20MintWithPostage(sender, receiver string, op, ticker, amount string, feeRate int64, postage int64) (txn *Brc20MintTransaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if postage < 546 {
		postage = 546
	}

	host, err := comingOrdHost(c.Chainnet)
	if err != nil {
		return
	}
	url := fmt.Sprintf("%v/mintWithPostage", host)
	header := map[string]string{"Content-Type": "application/json"}
	requestBody := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "mintWithPostage",
		"params": map[string]any{
			"source":      sender,
			"fee_rate":    feeRate,
			"content":     fmt.Sprintf(`{"p":"brc-20","op":"%v","tick":"%v","amt":"%v"}`, op, ticker, amount),
			"destination": receiver,
			"extension":   ".txt",
			"repeat":      1,

			"target_postage": postage,
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
