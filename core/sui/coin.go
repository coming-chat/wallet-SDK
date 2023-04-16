package sui

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (t *Token) getCoins(address string, limit uint) (coins types.Coins, err error) {
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
	pageCoins, err := cli.GetCoins(context.Background(), *addr, &coinType, nil, limit)
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

// @params coins we assume that it is sorted by balance descending.
func pickupTransferCoin(coins types.Coins, amount uint64, isSUI bool) (*PickedCoins, error) {
	amountInt := big.NewInt(0).SetUint64(amount)
	amountAddGas := big.NewInt(0).Add(amountInt, big.NewInt(MaxGasForPay))

	total := big.NewInt(0)
	pickedCoins := types.Coins{}
	enough := false
	meetSmaller := false
	for _, coin := range coins {
		balance := coin.Balance.Uint64()

		if balance < amount {
			meetSmaller = true
		} else if balance == amount {
			return &PickedCoins{
				Coins:  []types.Coin{coin},
				Total:  *amountInt,
				Amount: *amountInt,

				CanUseTransferObject: true,
			}, nil
		}

		if !enough {
			total = total.Add(total, big.NewInt(0).SetUint64(balance))
			pickedCoins = append(pickedCoins, coin)
			if isSUI {
				enough = total.Cmp(amountAddGas) >= 0
			} else {
				enough = total.Cmp(amountInt) >= 0
			}
		}
		if enough && meetSmaller {
			break
		}
	}
	if enough {
		return &PickedCoins{
			Coins:  pickedCoins,
			Total:  *total,
			Amount: *amountInt,

			CanUseTransferObject: false,
		}, nil
	}
	return nil, fmt.Errorf(`insufficient account balance`)
}

type PickedCoins struct {
	Coins  types.Coins
	Total  big.Int
	Amount big.Int

	CanUseTransferObject bool
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
