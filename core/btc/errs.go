package btc

import "errors"

var (
	ErrUnsupportedChain  = errors.New("Unsupported BTC chainnet")
	ErrHttpResponseParse = errors.New("Network data parsing error")
)
