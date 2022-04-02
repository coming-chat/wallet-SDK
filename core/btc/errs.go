package btc

import "errors"

var (
	ErrUnsupportedChain  = errors.New("Unsupported chain name")
	ErrHttpResponseParse = errors.New("Network data parsing error")
)
