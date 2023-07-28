package btc

import "errors"

var (
	ErrUnsupportedChain  = errors.New("Unsupported BTC chainnet")
	ErrHttpResponseParse = errors.New("Network data parsing error")

	ErrDecodeAddress = errors.New("Btc cannot support decode address to public key")

	ErrPsbtEncode                 = errors.New("unsupported Psbt tx encode")
	ErrPsbtUnsupportedAccountType = errors.New("account psbt sign unsupported")
)
