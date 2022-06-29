package doge

import (
	"fmt"
	"math/big"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type TransactionInput struct {
	Addresses []string `json:"addresses"`
}

type TransactionOutput struct {
	Value     *big.Int `json:"value"`
	Addresses []string `json:"addresses"`
}

type Transaction struct {
	// Demo: https://api.blockcypher.com/v1/doge/main/txs/7bc313903372776e1eb81d321e3fe27c9721ce8e71a9bcfee1bde6baea31b5c2
	Hash          string               `json:"hash"`
	Total         *big.Int             `json:"total"`
	Fees          *big.Int             `json:"fees"`
	Received      *time.Time           `json:"received"`
	Confirmed     *time.Time           `json:"confirmed"`
	Confirmations int64                `json:"confirmations"`
	Inputs        []*TransactionInput  `json:"inputs"`
	Outputs       []*TransactionOutput `json:"outputs"`
	OpReturn      string               `json:"data_protocol"`
}

func (t *Transaction) From() string {
	if len(t.Inputs) == 0 || len(t.Inputs[0].Addresses) == 0 {
		return ""
	}
	return t.Inputs[0].Addresses[0]
}

func (t *Transaction) ToAddressAndTransferAmount() (string, *big.Int) {
	if len(t.Outputs) == 0 {
		return "", nil
	}
	from := t.From()
	to := ""
	total := big.NewInt(0).SetBytes(t.Total.Bytes())
	for idx, out := range t.Outputs {
		address := ""
		if len(out.Addresses) > 0 {
			address = out.Addresses[0]
		}
		if address == from {
			total.Sub(total, out.Value)
			continue
		}
		if idx == 0 {
			to = address
		} else {
			to = fmt.Sprintf("%v, %v", to, address)
		}
	}
	if to == "" {
		// If there are no other recipients, we consider the user to transfer to himself
		to = from
	}

	return to, total
}

func (t *Transaction) Status() base.TransactionStatus {
	if t.Confirmations >= 1 {
		return base.TransactionStatusSuccess
	} else {
		return base.TransactionStatusPending
	}
}

func (t *Transaction) SdkDetail() *base.TransactionDetail {
	to, amount := t.ToAddressAndTransferAmount()
	var finishTime int64
	if t.Confirmed != nil {
		finishTime = t.Confirmed.Unix()
	}
	var amountString string
	if amount != nil {
		amountString = amount.String()
	}
	return &base.TransactionDetail{
		HashString:      t.Hash,
		Amount:          amountString,
		EstimateFees:    t.Fees.String(),
		FromAddress:     t.From(),
		ToAddress:       to,
		Status:          t.Status(),
		FinishTimestamp: finishTime,
	}
}
