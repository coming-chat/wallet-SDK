package sui

import (
	"context"
	"errors"
	"math/big"
	"strconv"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const MAX_INPUT_COUNT = 256

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

	// Actually, the maximum count of input coins is `MAX_INPUT_COUNT-1`
	maxInput := MAX_INPUT_COUNT - 1

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

	pageCoins, err := cli.GetCoins(context.Background(), *ownerAddr, &coinType, nil, uint(maxInput))
	if err != nil {
		return
	}
	if len(pageCoins.Data) <= 0 {
		return nil, errors.New("You do not have this coin yet")
	}
	count := base.Min(len(pageCoins.Data), maxInput)

	// We will try to merge all the coins as much as possible.
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

	pageCoins, err := cli.GetCoins(context.Background(), *ownerAddr, &coinType, nil, MAX_INPUT_COUNT)
	if err != nil {
		return
	}
	if len(pageCoins.Data) <= 0 {
		return nil, errors.New("You do not have this coin yet")
	}

	biggestCoin := pageCoins.Data[0]
	for _, coin := range pageCoins.Data {
		if biggestCoin.Balance.Uint64() < coin.Balance.Uint64() {
			biggestCoin = coin
		}
	}
	if biggestCoin.Balance.Uint64() <= amountInt {
		return nil, errors.New("no coin found that can split to the target amount.")
	}

	anotherAmount := biggestCoin.Balance.Uint64() - amountInt

	var txnBytes *types.TransactionBytes
	gasBudget := types.NewSafeSuiBigInt[uint64](MaxGasForPay)
	if coinType == SUI_COIN_TYPE && len(pageCoins.Data) == 1 {
		// only have one sui coin, it cannot be both split and used as a gas fee at the same time.
		txnBytes, err = cli.PaySui(context.Background(), *ownerAddr,
			[]types.ObjectId{biggestCoin.CoinObjectId},
			[]types.Address{*ownerAddr},
			[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountInt)},
			gasBudget)
	} else {
		txnBytes, err = cli.SplitCoin(context.Background(), *ownerAddr, biggestCoin.CoinObjectId,
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
