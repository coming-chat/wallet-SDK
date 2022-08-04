package eth

import (
	"context"
	"errors"
	"math/big"
	"strconv"

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

	header, err := client.RemoteRpcClient.HeaderByNumber(ctx, big.NewInt(-1))
	if err != nil {
		return nil, err
	}
	if header.BaseFee == nil {
		return nil, errors.New("The specified chain does not yet support EIP1559")
	}

	priorityFee, err := client.RemoteRpcClient.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}
	if priorityFee.Cmp(big.NewInt(0)) == 0 {
		priorityFee = big.NewInt(1e9) // 1 Gwei
	}

	return (&GasPrice{
		BaseFee:            header.BaseFee.String(),
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

func (c *Chain) EstimateGasLimit(msg *CallMsg) (gas *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if len(msg.msg.Data) > 0 {
		// any contract transaction
		gas = &base.OptionalString{Value: DEFAULT_CONTRACT_GAS_LIMIT}
	} else {
		// nomal transfer
		gas = &base.OptionalString{Value: DEFAULT_ETH_GAS_LIMIT}
	}

	client, err := GetConnection(c.RpcUrl)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()
	gasLimit, err := client.RemoteRpcClient.EstimateGas(ctx, msg.msg)
	if err != nil {
		return
	}
	gasString := ""
	if len(msg.msg.Data) > 0 {
		gasFloat := big.NewFloat(0).SetUint64(gasLimit)
		gasFloat = gasFloat.Mul(gasFloat, big.NewFloat(client.gasFactor()))
		gasInt, _ := gasFloat.Int(nil)
		gasString = gasInt.String()
	} else {
		gasString = strconv.FormatUint(gasLimit, 10)
	}

	return &base.OptionalString{Value: gasString}, nil
}
