package sui

import (
	"context"

	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

// @param maxGasBudget Default `MinGasBudget` if is 0.
func (c *Chain) BaseMoveCall(address, packageId, module, funcName string, typArgs []string, arg []any, maxGasBudget uint64) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	client, err := c.Client()
	if err != nil {
		return
	}
	addr, err := sui_types.NewAddressFromHex(address)
	if err != nil {
		return
	}
	packageIdHex, err := sui_types.NewObjectIdFromHex(packageId)
	if err != nil {
		return
	}
	if maxGasBudget == 0 {
		maxGasBudget = MinGasBudget
	}
	return c.EstimateTransactionFeeAndRebuildTransaction(maxGasBudget, func(gasBudget uint64) (*Transaction, error) {
		gasInt := types.NewSafeSuiBigInt(gasBudget)
		tx, err := client.MoveCall(
			context.Background(),
			*addr,
			*packageIdHex,
			module,
			funcName,
			typArgs,
			arg,
			nil,
			gasInt,
		)
		if err != nil {
			return nil, err
		}
		return &Transaction{
			Txn: *tx,
		}, nil
	})

}
