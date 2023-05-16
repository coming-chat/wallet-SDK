package sui

import (
	"errors"

	"github.com/coming-chat/go-sui/v2/types"
)

var (
	ErrNoCoinsFound        = types.ErrNoCoinsFound
	ErrInsufficientBalance = types.ErrInsufficientBalance
	ErrNeedMergeCoin       = types.ErrNeedMergeCoin
	ErrNeedSplitGasCoin    = types.ErrNeedSplitGasCoin

	ErrNoNeedMergeCoin = errors.New("existing coins exceed the target amount, no need to merge coins")
	ErrMergeOneCoin    = errors.New("only one coin does not need to merge coins")
)

func IsMergeError(err error) bool {
	return err != nil && err.Error() == ErrNeedMergeCoin.Error()
}

func IsSplitError(err error) bool {
	return err != nil && err.Error() == ErrNeedSplitGasCoin.Error()
}
