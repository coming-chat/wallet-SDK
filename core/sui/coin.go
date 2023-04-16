package sui

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (t *Token) getCoins(address string) (coins types.Coins, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := t.chain.Client()
	if err != nil {
		return
	}
	addr, err := types.NewAddressFromHex(address)
	if err != nil {
		return
	}
	coinType := t.rType.ShortString()
	pageCoins, err := cli.GetCoins(context.Background(), *addr, &coinType, nil, 0)
	if err != nil {
		return
	}
	coins = pageCoins.Data

	// sort by balance descend
	sort.Slice(coins, func(i, j int) bool {
		return coins[i].Balance.Uint64() > coins[j].Balance.Uint64()
	})
	return coins, nil
}

func pickupTransferCoin(coins types.Coins, amount string) (*PickedCoins, error) {
	amountInt, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf(`Invalid transfer amount "%v".`, amount)
	}
	need := big.NewInt(0).Set(amountInt)

	estimateGasPerCoin := big.NewInt(MaxGasForPay)
	total := big.NewInt(0)
	pickedCoins := types.Coins{}
	for _, coin := range coins {
		need = need.Add(need, estimateGasPerCoin)
		total = total.Add(total, big.NewInt(coin.Balance.Int64()))
		pickedCoins = append(pickedCoins, coin)
		if total.Cmp(need) >= 0 {
			return &PickedCoins{
				Coins:  pickedCoins,
				Total:  total,
				Amount: amountInt,
			}, nil
		}
	}
	return nil, fmt.Errorf(`Insufficient account balance "%v"`, total.String())
}

type PickedCoins struct {
	Coins  types.Coins
	Total  *big.Int
	Amount *big.Int
}

func (cs *PickedCoins) CoinIds() []types.ObjectId {
	coinIds := []types.ObjectId{}
	for _, coin := range cs.Coins {
		coinIds = append(coinIds, coin.CoinObjectId)
	}
	return coinIds
}

func (cs *PickedCoins) EstimateTotalGas() uint64 {
	return MaxGasForTransfer
}

func (cs *PickedCoins) EstimateMergeGas() uint64 {
	return MaxGasForTransfer
}
