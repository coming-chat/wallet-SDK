package sui

import (
	"context"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const gasBudget = MaxGasBudget

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
	suiToken := NewTokenMain(c)
	coins, err := suiToken.getCoins(address)
	if err != nil {
		return
	}
	coin, err := coins.PickCoinNoLess(gasBudget)
	if err != nil {
		return
	}
	gasInt := types.NewSafeSuiBigInt[uint64](gasBudget)
	tx, err := client.MoveCall(
		context.Background(),
		*addr,
		*packageIdHex,
		module,
		funcName,
		typArgs,
		arg,
		&coin.CoinObjectId,
		gasInt,
	)
	if err != nil {
		return
	}
	return &Transaction{
		Txn:          *tx,
		MaxGasBudget: gasBudget,
	}, nil
}
