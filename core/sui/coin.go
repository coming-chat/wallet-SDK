package sui

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (t *Token) getCoins(address string) (cs Coins, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	cs = []*Coin{}

	cli, err := t.chain.client()
	if err != nil {
		return
	}
	addr, err := types.NewAddressFromHex(address)
	if err != nil {
		return
	}
	objects, err := cli.GetObjectsOwnedByAddress(context.Background(), *addr)
	if err != nil {
		return
	}

	type coinData struct {
		Fields struct {
			Balance uint64 `json:"balance"`
		} `json:"fields"`
	}
	cointype := t.coinType()
	var res *types.ObjectRead
	var bytes []byte
	for _, obj := range objects {
		if obj.Type != cointype {
			continue
		}

		res, err = cli.GetObject(context.Background(), *obj.ObjectId)
		if err != nil {
			return
		}
		if res.Status != types.ObjectStatusExists {
			continue
		}
		bytes, err = json.Marshal(res.Details.Data)
		if err != nil {
			return
		}
		coindata := coinData{}
		err = json.Unmarshal(bytes, &coindata)
		if err != nil {
			return
		}

		cs = append(cs, &Coin{
			ObjectId: *obj.ObjectId,
			Balance:  coindata.Fields.Balance,
		})
	}

	// sort by balance descend
	sort.Slice(cs, func(i, j int) bool {
		return cs[i].Balance > cs[j].Balance
	})
	return cs, nil
}

type Coin struct {
	ObjectId types.ObjectId
	Balance  uint64
}

type Coins []*Coin

func (cs Coins) Total() *big.Int {
	total := big.NewInt(0)
	for _, c := range cs {
		total = total.Add(total, big.NewInt(int64(c.Balance)))
	}
	return total
}

func (cs Coins) PickupTransferCoin(amount string) (*PickedCoins, error) {
	amountInt, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf(`Invalid transfer amount "%v".`, amount)
	}
	need := big.NewInt(0).Set(amountInt)

	estimateGasPerCoin := big.NewInt(MaxGasForMerge)
	total := big.NewInt(0)
	pickedCoins := Coins{}
	for _, coin := range cs {
		need = need.Add(need, estimateGasPerCoin)
		total = total.Add(total, big.NewInt(int64(coin.Balance)))
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
	Coins  Coins
	Total  *big.Int
	Amount *big.Int
}

func (cs *PickedCoins) EstimateGas() uint64 {
	return uint64(len(cs.Coins)) * MaxGasBudget
}

func (cs *PickedCoins) LastCoin() *Coin {
	if len(cs.Coins) == 0 {
		return nil
	}
	return cs.Coins[len(cs.Coins)-1]
}

func (cs *PickedCoins) LastCoinTransferAmount() uint64 {
	surplus := big.NewInt(0).Sub(cs.Total, cs.Amount).Uint64()
	return cs.LastCoin().Balance - surplus
}
