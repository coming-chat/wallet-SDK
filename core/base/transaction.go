package base

import (
	"encoding/json"
	"strings"
)

type TransactionStatus = SDKEnumInt

const (
	TransactionStatusNone    TransactionStatus = 0
	TransactionStatusPending TransactionStatus = 1
	TransactionStatusSuccess TransactionStatus = 2
	TransactionStatusFailure TransactionStatus = 3
)

type Transaction struct {
}

// Transaction details that can be fetched from the chain
type TransactionDetail struct {
	// hash string on chain
	HashString string

	// transaction amount
	Amount string

	EstimateFees string

	// sender's address
	FromAddress string
	// receiver's address
	ToAddress string

	Status TransactionStatus
	// transaction completion timestamp (s), 0 if Status is in Pending
	FinishTimestamp int64
	// failure message
	FailureMessage string

	// If this transaction is a CID transfer, its value will be the CID, otherwise it is empty
	CIDNumber string
	// If this transaction is a NFT transfer, its value will be the Token name, otherwise it is empty
	TokenName string
}

// Check the `CIDNumber` is not empty.
func (d *TransactionDetail) IsCIDTransfer() bool {
	return strings.TrimSpace(d.CIDNumber) != ""
}

// Check the `TokenName` is not empty.
func (d *TransactionDetail) IsNFTTransfer() bool {
	return strings.TrimSpace(d.TokenName) != ""
}

func (d *TransactionDetail) JsonString() string {
	b, err := json.Marshal(d)
	if err != nil {
		return ""
	}
	return string(b)
}
