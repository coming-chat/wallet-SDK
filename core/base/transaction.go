package base

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
}
