package sui

import (
	"context"
	"math/big"
	"strconv"

	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/fardream/go-bcs/bcs"
)

type MergeCoinRequest struct {
	Owner        string
	CoinType     string
	TargetAmount string

	// queried coins
	Coins      types.PickedCoins
	CoinsCount int
	// If the transaction is executed, owner will receive a coin amount that is not less than this.
	EstimateAmount string
	// Will the goal of merging coins of a specified amount be achieved?
	WillBeAchieved bool
}

type MergeCoinPreview struct {
	EstimateGasFee int64

	// The original request
	Request *MergeCoinRequest

	// The merge coins transaction
	Transaction *Transaction

	// If the transaction is executed, owner will receive a coin amount that is not less than this.
	EstimateAmount string
	// Did the simulated transaction execute successfully?
	SimulateSuccess bool
	// Due to the results obtained through simulated execution, we may know that the balance may increase and the value of this state may be inconsistent with the value in the request.
	WillBeAchieved bool
}

func (c *Chain) BuildMergeCoinPreview(request *MergeCoinRequest) (preview *MergeCoinPreview, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	ownerAddr, err := sui_types.NewAddressFromHex(request.Owner)
	if err != nil {
		return
	}

	var pickedCoins *types.PickedCoins
	var pickedGasCoins *types.PickedCoins
	if request.CoinType == SUI_COIN_TYPE {
		pickedCoins = nil
		pickedGasCoins = &request.Coins
	} else {
		pickedCoins = &request.Coins
		pickedGasCoins, err = c.PickGasCoins(*ownerAddr, MaxGasForTransfer)
		if err != nil {
			return
		}
	}

	gasPrice, _ := c.CachedGasPrice()
	maxGasBudget := maxGasBudget(pickedGasCoins, MaxGasForTransfer)
	txn, err := c.EstimateTransactionFeeAndRebuildTransactionBCS(maxGasBudget, func(gasBudget uint64) (*Transaction, error) {
		ptb := sui_types.NewProgrammableTransactionBuilder()
		if request.CoinType == SUI_COIN_TYPE {
			err = ptb.PayAllSui(*ownerAddr)
		} else {
			err = ptb.Pay(
				pickedCoins.CoinRefs(),
				[]sui_types.SuiAddress{*ownerAddr},
				[]uint64{pickedCoins.TotalAmount.Uint64()},
			)
		}
		if err != nil {
			return nil, err
		}

		pt := ptb.Finish()
		tx := sui_types.NewProgrammable(*ownerAddr, pickedGasCoins.CoinRefs(), pt, gasBudget, gasPrice)
		txBytes, err := bcs.Marshal(tx)
		if err != nil {
			return nil, err
		}
		return &Transaction{TxnBytes: txBytes}, nil
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
	ownerAddr, err := sui_types.NewAddressFromHex(owner)
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
	pickedCoins := pickAllCoins(pageCoins)
	if len(pickedCoins.Coins) == 1 {
		return nil, ErrMergeOneCoin
	}
	willBeAchieved := pickedCoins.TotalAmount.Cmp(big.NewInt(0).SetUint64(amountInt)) >= 0

	return &MergeCoinRequest{
		Owner:        owner,
		CoinType:     coinType,
		TargetAmount: targetAmount,

		Coins:          *pickedCoins,
		CoinsCount:     len(pickedCoins.Coins),
		EstimateAmount: pickedCoins.TotalAmount.String(),
		WillBeAchieved: willBeAchieved,
	}, nil
}

// @param coinType Default is `SUI_COIN_TYPE`
func (c *Chain) BuildSplitCoinTransaction(owner, coinType, targetAmount string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if coinType == "" {
		coinType = SUI_COIN_TYPE
	}
	ownerAddr, err := sui_types.NewAddressFromHex(owner)
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
	// We'd better split a coin that can meet the target amount and the remaining coin value can be greater than 1SUI
	// so that the transaction can be executed smoothly.
	needAmount := amountInt + 1e9
	pickedCoins, err := types.PickupCoins(pageCoins, *big.NewInt(0).SetUint64(needAmount), MaxGasForPay, MAX_INPUT_COUNT_MERGE, 0)
	if err != nil {
		if err.Error() == ErrInsufficientBalance.Error() {
			if types.Coins(pageCoins.Data).TotalBalance().Uint64() < (amountInt + MinGasBudget*2) {
				return nil, ErrInsufficientBalance
			}
			pickedCoins = &types.PickedCoins{
				Coins: pageCoins.Data, // all coins should be used to merge
			}
			err = nil
		} else {
			return nil, err
		}
	}

	maxGasBudget := maxGasBudget(pickedCoins, MaxGasForPay)
	return c.EstimateTransactionFeeAndRebuildTransaction(maxGasBudget, func(gasBudget uint64) (*Transaction, error) {
		var txnBytes *types.TransactionBytes
		gasInt := types.NewSafeSuiBigInt(gasBudget)
		if coinType == SUI_COIN_TYPE && (pickedCoins.Count() > 1 || len(pageCoins.Data) == 1) {
			txnBytes, err = cli.PaySui(context.Background(), *ownerAddr,
				pickedCoins.CoinIds(),
				[]sui_types.SuiAddress{*ownerAddr},
				[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountInt)},
				gasInt)
		} else if pickedCoins.Count() > 1 {
			txnBytes, err = cli.Pay(context.Background(), *ownerAddr,
				pickedCoins.CoinIds(),
				[]sui_types.SuiAddress{*ownerAddr},
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
