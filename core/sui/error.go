package sui

import (
	"github.com/coming-chat/go-sui/types"
)

var (
	ErrNoCoinsFound        = types.ErrNoCoinsFound
	ErrInsufficientBalance = types.ErrInsufficientBalance
	ErrNeedMergeCoin       = types.ErrNeedMergeCoin
	ErrNeedSplitGasCoin    = types.ErrNeedSplitGasCoin
)

func IsMergeError(err error) bool {
	return err.Error() == ErrNeedMergeCoin.Error()
}

func IsSplitError(err error) bool {
	return err.Error() == ErrNeedSplitGasCoin.Error()
}
