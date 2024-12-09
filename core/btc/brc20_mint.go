package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type Brc20MintTransaction struct {
	NetworkFee  int64 `json:"network_fee"`
	SatpointFee int64 `json:"satpoint_fee"`
	ServiceFee  int64 `json:"service_fee"`
	CommitFee   int64 `json:"commit_fee"`
	CommitVsize int64 `json:"commit_vsize"`

	CommitId    string            `json:"commit_id"`
	Commit      string            `json:"commit"`
	Reveal      *base.StringArray `json:"reveal"`
	Inscription *base.StringArray `json:"inscription"`

	CommitCustom *Brc20CommitCustom `json:"commit_custom"`

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
	url := fmt.Sprintf("%v/MintInscription", host)
	header := map[string]string{"Content-Type": "application/json"}
	requestBody := map[string]any{
		"commit_fee_rate": strconv.FormatInt(feeRate, 10),
		"reveal_fee_rate": strconv.FormatInt(feeRate, 10),
		"source":          sender,
		"destination":     receiver,
		"inscriptions": []string{
			fmt.Sprintf(`{"p":"brc-20","op":"%v","tick":"%v","amt":"%v"}`, op, ticker, amount),
		},
		"repeats": []int{
			1,
		},
		"postage": postage,
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
		NetworkFee  int64 `json:"network_fee"`
		SatpointFee int64 `json:"satpoint_fee"`
		ServiceFee  int64 `json:"service_fee"`
		CommitFee   int64 `json:"commit_fee"`
		CommitVsize int64 `json:"commit_vsize"`

		CommitId   string            `json:"commit_id"`
		CommitPsbt string            `json:"commit_psbt"`
		RevealTxs  *base.StringArray `json:"reveal_txs"`
		RevealIds  *base.StringArray `json:"reveal_ids"`

		CommitCustom *Brc20CommitCustom `json:"commit_custom"`
	}
	err = json.Unmarshal(resp.Body, &r)
	if err != nil {
		return nil, err
	}
	return &Brc20MintTransaction{
		CommitId:    r.CommitId,
		Commit:      r.CommitPsbt,
		Reveal:      r.RevealTxs,
		Inscription: r.RevealIds,

		CommitCustom: r.CommitCustom,

		NetworkFee:  r.NetworkFee,
		SatpointFee: r.SatpointFee,
		ServiceFee:  r.ServiceFee,
		CommitFee:   r.CommitFee,
		CommitVsize: r.CommitVsize,
	}, nil
}
