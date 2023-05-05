package sui

import (
	"context"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (c *Chain) BaseMoveCall(address, packageId, module, funcName string, typArgs []string, arg []any) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	client, err := c.Client()
	if err != nil {
		return
	}
	addr, err := types.NewAddressFromHex(address)
	if err != nil {
		return
	}
	packageIdHex, err := types.NewHexData(packageId)
	if err != nil {
		return
	}
	return c.EstimateTransactionFeeAndRebuildTransaction(MinGasBudget, func(gasBudget uint64) (*Transaction, error) {
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
