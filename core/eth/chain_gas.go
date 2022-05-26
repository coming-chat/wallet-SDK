package eth

import (
	"context"
	"errors"
	"math/big"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type GasPrice struct {
	// Pending block baseFee in wei.
	BaseFee            string
	SuggestPriorityFee string

	MaxPriorityFee string
	MaxFee         string
}

// MaxPriorityFee = SuggestPriorityFee * 1.0
// MaxFee = (MaxPriorityFee + BaseFee) * 1.0
func (g *GasPrice) UseLow() *GasPrice {
	return g.UseRate(1.0, 1.0)
}

// MaxPriorityFee = SuggestPriorityFee * 1.5
// MaxFee = (MaxPriorityFee + BaseFee) * 1.11
func (g *GasPrice) UseAverage() *GasPrice {
	return g.UseRate(1.5, 1.11)
}

// MaxPriorityFee = SuggestPriorityFee * 2.0
// MaxFee = (MaxPriorityFee + BaseFee) * 1.5
func (g *GasPrice) UseHigh() *GasPrice {
	return g.UseRate(2.0, 1.5)
}

// MaxPriorityFee = SuggestPriorityFee * priorityRate
// MaxFee = (MaxPriorityFee + BaseFee) * maxFeeRate
func (g *GasPrice) UseRate(priorityRate, maxFeeRate float64) *GasPrice {
	suggestPriorityFloat, ok := big.NewFloat(0).SetString(g.SuggestPriorityFee)
	if !ok {
		suggestPriorityFloat = big.NewFloat(0)
	}
	baseFeeFloat, ok := big.NewFloat(0).SetString(g.BaseFee)
	if !ok {
		baseFeeFloat = big.NewFloat(0)
	}

	maxPriorityFloat := big.NewFloat(0).Mul(suggestPriorityFloat, big.NewFloat(priorityRate))
	sumFloat := big.NewFloat(0).Add(maxPriorityFloat, baseFeeFloat)
	maxFeeFloat := big.NewFloat(0).Mul(sumFloat, big.NewFloat(maxFeeRate))
	maxPriorityInt, _ := maxPriorityFloat.Int(nil)
	maxFeeInt, _ := maxFeeFloat.Int(nil)
	return &GasPrice{
		BaseFee:            g.BaseFee,
		SuggestPriorityFee: g.SuggestPriorityFee,
		MaxPriorityFee:     maxPriorityInt.String(),
		MaxFee:             maxFeeInt.String(),
	}
}

// The gas price use average grade default.
func (c *Chain) SuggestGasPriceEIP1559() (*GasPrice, error) {
	client, err := GetConnection(c.RpcUrl)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	// We best get the pending block base fee, but now, get the pending block will be crash
	// Waiting for new version support, We need call `HeaderByNumber(ctx, big.NewInt(-1))`
	header, err := client.RemoteRpcClient.HeaderByNumber(ctx, nil) // now we can only input `nil`
	if err != nil {
		return nil, err
	}
	if header.BaseFee == nil {
		return nil, errors.New("The specified chain does not yet support EIP1559")
	}
	// calculate formular refence from https://www.blocknative.com/blog/eip-1559-fees#determining-the-base-fee
	// nextBlock BaseFee = CurrentBaseFee * (0.875 + 0.25 * CurrentUsage / CurrentLimit )
	numerator := header.GasUsed
	if header.GasUsed > header.GasLimit {
		numerator = header.GasLimit
	}
	usageRate := big.NewFloat(float64(numerator) / float64(header.GasLimit))
	rate := big.NewFloat(0).Add(big.NewFloat(0.875), big.NewFloat(0).Mul(big.NewFloat(0.25), usageRate))
	pendingBaseFeeFloat := big.NewFloat(0).Mul(big.NewFloat(0).SetInt(header.BaseFee), rate)
	pendingBaseFeeInt, _ := pendingBaseFeeFloat.Int(nil)

	priorityFee, err := client.RemoteRpcClient.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}

	return (&GasPrice{
		BaseFee:            pendingBaseFeeInt.String(),
		SuggestPriorityFee: priorityFee.String(),
	}).UseAverage(), nil
}

func (c *Chain) SuggestGasPrice() (*base.OptionalString, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return nil, err
	}
	price, err := chain.SuggestGasPrice()
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: price}, nil
}
