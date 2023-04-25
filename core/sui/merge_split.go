package sui

import (
	"context"
	"math/big"
	"sort"
	"strconv"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type MergeCoinRequest struct {
	Owner        string
	CoinType     string
	TargetAmount string

	// queried coins
	Coins      types.Coins
	CoinsCount int
	// If the transaction is executed, owner will receive a coin amount that is not less than this.
	EstimateAmount string
	// Will the goal of merging coins of a specified amount be achieved?
	WillBeAchieved bool
}

type MergeCoinPreview struct {
	// The original request
	Request *MergeCoinRequest

	// The merge coins transaction
	Transaction *Transaction
	// Did the simulated transaction execute successfully?
	SimulateSuccess bool
	EstimateGasFee  int64

	// If the transaction is executed, owner will receive a coin amount that is not less than this.
	EstimateAmount string
	// Due to the results obtained through simulated execution, we may know that the balance may increase and the value of this state may be inconsistent with the value in the request.
	WillBeAchieved bool
}

func (c *Chain) BuildMergeCoinPreview(request *MergeCoinRequest) (preview *MergeCoinPreview, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	ownerAddr, err := types.NewAddressFromHex(request.Owner)
	if err != nil {
		return
	}
	cli, err := c.Client()
	if err != nil {
		return
	}

	totalAmount := big.NewInt(0)
	mergeIds := make([]types.ObjectId, 0)
	for _, coin := range request.Coins {
		totalAmount.Add(totalAmount, big.NewInt(0).SetUint64(coin.Balance.Uint64()))
		mergeIds = append(mergeIds, coin.CoinObjectId)
	}

	txn, err := c.EstimateTransactionFeeAndRebuildTransaction(MaxGasForPay, func(gasBudget uint64) (*Transaction, error) {
		var txnBytes *types.TransactionBytes
		gasInt := types.NewSafeSuiBigInt(gasBudget)
		if request.CoinType == SUI_COIN_TYPE {
			txnBytes, err = cli.PayAllSui(context.Background(), *ownerAddr, *ownerAddr, mergeIds, gasInt)
		} else {
			txnBytes, err = cli.Pay(context.Background(), *ownerAddr, mergeIds,
				[]types.Address{*ownerAddr},
				[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(totalAmount.Uint64())},
				nil, gasInt)
		}
		if err != nil {
			return nil, err
		}
		return &Transaction{Txn: *txnBytes}, nil
	})

	return &MergeCoinPreview{
		Request: request,

		Transaction:     txn,
		SimulateSuccess: true,
		EstimateGasFee:  txn.EstimateGasFee,
		EstimateAmount:  request.EstimateAmount,
		WillBeAchieved:  request.WillBeAchieved,
	}, nil
}

// @param coinType Default is `SUI_COIN_TYPE`
func (c *Chain) BuildMergeCoinRequest(owner, coinType, targetAmount string) (req *MergeCoinRequest, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if coinType == "" {
		coinType = SUI_COIN_TYPE
	}
	ownerAddr, err := types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	amountInt, err := strconv.ParseUint(targetAmount, 10, 64)
	if err != nil {
		return
	}
	cli, err := c.Client()
	if err != nil {
		return
	}

	pageCoins, err := cli.GetCoins(context.Background(), *ownerAddr, &coinType, nil, MAX_INPUT_COUNT_MERGE)
	if err != nil {
		return
	}
	if len(pageCoins.Data) <= 0 {
		return nil, ErrNoCoinsFound
	}
	coins := pageCoins.Data
	sort.Slice(coins, func(i, j int) bool {
		return coins[i].Balance.Uint64() > coins[j].Balance.Uint64()
	})

	amountBigInt := big.NewInt(0).SetUint64(amountInt)
	totalAmount := big.NewInt(0)
	mergingCoins := []types.Coin{}
	for _, coin := range coins {
		if coin.Balance.Uint64() >= amountBigInt.Uint64() {
			return nil, ErrNoNeedMergeCoin
		}
		totalAmount.Add(totalAmount, big.NewInt(0).SetUint64(coin.Balance.Uint64()))
		mergingCoins = append(mergingCoins, coin)
		if totalAmount.Cmp(amountBigInt) >= 0 {
			break
		}
	}
	if len(mergingCoins) == 1 {
		return nil, ErrMergeOneCoin
	}

	return &MergeCoinRequest{
		Owner:        owner,
		CoinType:     coinType,
		TargetAmount: targetAmount,

		Coins:          mergingCoins,
		CoinsCount:     len(mergingCoins),
		EstimateAmount: totalAmount.String(),
		WillBeAchieved: totalAmount.Cmp(amountBigInt) >= 0,
	}, nil
}

// @param coinType Default is `SUI_COIN_TYPE`
func (c *Chain) BuildSplitCoinTransaction(owner, coinType, targetAmount string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if coinType == "" {
		coinType = SUI_COIN_TYPE
	}
	ownerAddr, err := types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	amountInt, err := strconv.ParseUint(targetAmount, 10, 64)
	if err != nil {
		return
	}
	cli, err := c.Client()
	if err != nil {
		return
	}

	pageCoins, err := cli.GetCoins(context.Background(), *ownerAddr, &coinType, nil, MAX_INPUT_COUNT_MERGE)
	if err != nil {
		return
	}
	pickedCoins, err := types.PickupCoins(pageCoins, *big.NewInt(0).SetUint64(amountInt), MAX_INPUT_COUNT_MERGE, false)
	if err != nil {
		return
	}

	return c.EstimateTransactionFeeAndRebuildTransaction(MaxGasForPay, func(gasBudget uint64) (*Transaction, error) {
		var txnBytes *types.TransactionBytes
		gasInt := types.NewSafeSuiBigInt(gasBudget)
		if coinType == SUI_COIN_TYPE && (pickedCoins.Count() > 1 || len(pageCoins.Data) == 1) {
			txnBytes, err = cli.PaySui(context.Background(), *ownerAddr,
				pickedCoins.CoinIds(),
				[]types.Address{*ownerAddr},
				[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountInt)},
				gasInt)
		} else if pickedCoins.Count() > 1 {
			txnBytes, err = cli.Pay(context.Background(), *ownerAddr,
				pickedCoins.CoinIds(),
				[]types.Address{*ownerAddr},
				[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountInt)},
				nil, gasInt)
		} else {
			theCoin := pickedCoins.Coins[0]
			anotherAmount := theCoin.Balance.Uint64() - amountInt
			txnBytes, err = cli.SplitCoin(context.Background(), *ownerAddr,
				theCoin.CoinObjectId,
				[]types.SafeSuiBigInt[uint64]{
					types.NewSafeSuiBigInt(amountInt),
					types.NewSafeSuiBigInt(anotherAmount),
				}, nil, gasInt)
		}
		if err != nil {
			return nil, err
		}

		return &Transaction{
			Txn: *txnBytes,
		}, nil
	})
}
