package sui

import (
	"context"

	"github.com/coming-chat/go-sui/types"
)

const gasBudget = MaxGasBudget

func (c *Chain) BaseMoveCall(address, packageId, module, funcName string, typArgs []string, arg []any) (*Transaction, error) {
	client, err := c.Client()
	if err != nil {
		return nil, err
	}
	addr, err := types.NewAddressFromHex(address)
	if err != nil {
		return nil, err
	}
	packageIdHex, err := types.NewHexData(packageId)
	if err != nil {
		return nil, err
	}
	suiToken := NewTokenMain(c)
	coins, err := suiToken.getCoins(address)
	if err != nil {
		return nil, err
	}
	coin, err := coins.PickCoinNoLess(gasBudget)
	if err != nil {
		return nil, err
	}
	tx, err := client.MoveCall(
		context.Background(),
		*addr,
		*packageIdHex,
		module,
		funcName,
		typArgs,
		arg,
		&coin.CoinObjectId,
		gasBudget,
	)
	if err != nil {
		return nil, err
	}
	return &Transaction{
		Txn:          *tx,
		MaxGasBudget: gasBudget,
	}, nil
}
