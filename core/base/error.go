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

	ErrEstimateGasNeedPublicKey = errors.New("the estimated fee should invoking `EstimateTransactionFeeUsePublicKey`")
)
