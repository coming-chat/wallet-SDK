package base

import "errors"

var (
	ErrUnsupportedFunction = errors.New("this method is not supported")

	ErrInvalidChainType       = errors.New("invalid chain type")
	ErrInvalidAccountType     = errors.New("invalid account type")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)
