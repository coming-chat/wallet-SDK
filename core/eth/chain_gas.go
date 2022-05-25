package eth

import (
	"context"
	"math/big"
)

type GasPrice struct {
	// Latest block baseFee in wei.
	BaseFee string

	// Suggest maxPriorityFeePerGas * 1 in wei
	PriorityFeeLow string
	// Suggest maxPriorityFeePerGas * 1.5 in wei
	PriorityFeeAverage string
	// Suggest maxPriorityFeePerGas * 2 in wei
	PriorityFeeHigh string
}

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

	priorityFloat := big.NewFloat(0).SetInt(priorityFee)
	averageFloat := big.NewFloat(0).Mul(priorityFloat, big.NewFloat(1.5))
	averageInt, _ := averageFloat.Int(nil)
	priorityHigh := big.NewInt(0).Mul(priorityFee, big.NewInt(2))

	return &GasPrice{
		BaseFee: pendingBaseFeeInt.String(),

		PriorityFeeLow:     priorityFee.String(),
		PriorityFeeAverage: averageInt.String(),
		PriorityFeeHigh:    priorityHigh.String(),
	}, nil
}
