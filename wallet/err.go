package wallet

import "errors"

var (
	errNilKey      = errors.New("no mnemonic or private key")
	errNilMetadata = errors.New("no metadata")
	errNotSigned   = errors.New("transaction not signed")
	errNoPublicKey = errors.New("transaction no public key")
)
