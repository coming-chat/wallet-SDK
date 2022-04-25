package base

import "errors"

var (
	ErrUnsupportedFunction = errors.New("This chain does not support this feature.")
)
