package sui

import (
	"context"
	"math/big"
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
	targetAmount, err := strconv.ParseUint(request.TargetAmount, 10, 64)
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

	var txnBytes *types.TransactionBytes
	gasBudget := types.NewSafeSuiBigInt[uint64](MaxGasForPay)
	if request.CoinType == SUI_COIN_TYPE {
		txnBytes, err = cli.PayAllSui(context.Background(), *ownerAddr, *ownerAddr, mergeIds, gasBudget)
	} else {
		txnBytes, err = cli.Pay(context.Background(), *ownerAddr, mergeIds,
			[]types.Address{*ownerAddr},
			[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(totalAmount.Uint64())},
			nil, gasBudget)
	}
	if err != nil {
		return
	}

	rawTxn := &Transaction{
		Txn:          *txnBytes,
		MaxGasBudget: gasBudget.Int64(),
	}
	simulate, err := cli.DryRunTransaction(context.Background(), txnBytes)
	if err != nil || !simulate.Effects.Data.IsSuccess() {
		return &MergeCoinPreview{
			Request:        request,
			Transaction:    rawTxn,
			EstimateAmount: totalAmount.String(),
			WillBeAchieved: totalAmount.Uint64() > targetAmount,
		}, nil
	}

	balanceChange := int64(0)
	rawTxn.EstimateGasFee = simulate.Effects.Data.GasFee()
	for _, c := range simulate.BalanceChanges {
		if c.CoinType == request.CoinType {
			balanceChange, _ = strconv.ParseInt(c.Amount, 10, 64)
			break
		}
	}
	estimateAmount := big.NewInt(0).Add(totalAmount, big.NewInt(balanceChange))

	return &MergeCoinPreview{
		Request: request,

		Transaction:     rawTxn,
		SimulateSuccess: true,
		EstimateGasFee:  simulate.Effects.Data.GasFee(),
		EstimateAmount:  estimateAmount.String(),
		WillBeAchieved:  estimateAmount.Uint64() > targetAmount,
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

	// We will try to merge all the coins as much as possible.
	count := base.Min(len(pageCoins.Data), MAX_INPUT_COUNT_MERGE)
	coins := types.Coins(pageCoins.Data[0:count])
	totalAmount := coins.TotalBalance()

	return &MergeCoinRequest{
		Owner:        owner,
		CoinType:     coinType,
		TargetAmount: targetAmount,

		Coins:          coins,
		CoinsCount:     len(coins),
		EstimateAmount: totalAmount.String(),
		WillBeAchieved: totalAmount.Uint64() > amountInt,
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

	var txnBytes *types.TransactionBytes
	gasBudget := types.NewSafeSuiBigInt[uint64](MaxGasForPay)
	if coinType == SUI_COIN_TYPE && (pickedCoins.Count() > 1 || len(pageCoins.Data) == 1) {
		txnBytes, err = cli.PaySui(context.Background(), *ownerAddr,
			pickedCoins.CoinIds(),
			[]types.Address{*ownerAddr},
			[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountInt)},
			gasBudget)
	} else if pickedCoins.Count() > 1 {
		txnBytes, err = cli.Pay(context.Background(), *ownerAddr,
			pickedCoins.CoinIds(),
			[]types.Address{*ownerAddr},
			[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountInt)},
			nil, gasBudget)
	} else {
		theCoin := pickedCoins.Coins[0]
		anotherAmount := theCoin.Balance.Uint64() - amountInt
		txnBytes, err = cli.SplitCoin(context.Background(), *ownerAddr,
			theCoin.CoinObjectId,
			[]types.SafeSuiBigInt[uint64]{
				types.NewSafeSuiBigInt(amountInt),
				types.NewSafeSuiBigInt(anotherAmount),
			}, nil, gasBudget)
	}
	if err != nil {
		return
	}

	return &Transaction{
		Txn:          *txnBytes,
		MaxGasBudget: MaxGasForPay,
	}, nil
}
