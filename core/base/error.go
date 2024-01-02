package base

import "errors"

var (
	ErrUnsupportedFunction = errors.New("this method is not supported")

	ErrInvalidPrivateKey = errors.New("invalid private key")
	ErrInvalidPublicKey  = errors.New("invalid public key")
	ErrInvalidAddress    = errors.New("invalid address")

	ErrInvalidChainType       = errors.New("invalid chain type")
	ErrInvalidAccountType     = errors.New("invalid account type")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
	ErrInvalidTransactionData = errors.New("invalid transaction data")
	ErrInvalidTransactionHash = errors.New("invalid transaction hash")

	ErrInvalidAccountAddress = errors.New("invalid account address")
	ErrInvalidAmount         = errors.New("invalid amount")

	ErrMissingTransaction = errors.New("missing transaction information")

	ErrNotCoinTransferTxn = errors.New("not a coin transfer transaction")

	ErrEstimateGasNeedPublicKey = errors.New("the estimated fee should invoking `EstimateTransactionFeeUsePublicKey`")

	ErrInsufficientBalance = errors.New("insufficient account balance")
)
