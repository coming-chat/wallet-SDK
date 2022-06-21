package eth

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/params"
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

type OptimismLayer2Gas struct {
	L1GasLimit string
	L1GasPrice string
	L2GasLimit string
	L2GasPrice string
}

// l1GasLimit * l1GasPrice + l2Gaslimit * l2GasPrice
func (g *OptimismLayer2Gas) GasFee() string {
	l1Limit, ok := big.NewInt(0).SetString(g.L1GasLimit, 10)
	if !ok {
		l1Limit = big.NewInt(0)
	}
	l1Price, ok := big.NewInt(0).SetString(g.L1GasPrice, 10)
	if !ok {
		l1Price = big.NewInt(0)
	}
	l2Limit, ok := big.NewInt(0).SetString(g.L2GasLimit, 10)
	if !ok {
		l2Limit = big.NewInt(0)
	}
	l2Price, ok := big.NewInt(0).SetString(g.L2GasPrice, 10)
	if !ok {
		l2Price = big.NewInt(0)
	}
	l1Fee := big.NewInt(0).Mul(l1Limit, l1Price)
	l2Fee := big.NewInt(0).Mul(l2Limit, l2Price)
	return big.NewInt(0).Add(l1Fee, l2Fee).String()
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
	gasStr := strconv.FormatUint(gasLimit, 10)

	return &base.OptionalString{Value: gasStr}, nil
}

func (c *Chain) EstimateGasLimitLayer2(msg *CallMsg) (*OptimismLayer2Gas, error) {
	l2GasLimitString, err := c.EstimateGasLimit(msg)
	if err != nil {
		return nil, err
	}
	l2GasPrice := msg.msg.GasPrice

	// We need fetch the ethereum mainnet Gas Price
	ethMainRpc := "https://geth-mainnet.coming.chat"
	l1GasPriceString, err := NewChainWithRpc(ethMainRpc).SuggestGasPrice()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(msg.msg)
	if err != nil {
		return nil, err
	}
	l1GasLimit := calculateL1GasLimit(data, overhead)

	return &OptimismLayer2Gas{
		L1GasLimit: l1GasLimit.String(),
		L1GasPrice: l1GasPriceString.Value,
		L2GasLimit: l2GasLimitString.Value,
		L2GasPrice: l2GasPrice.String(),
	}, nil
}

const overhead uint64 = 200 * params.TxDataNonZeroGasEIP2028

func calculateL1GasLimit(data []byte, overhead uint64) *big.Int {
	zeroes, ones := zeroesAndOnes(data)
	zeroesCost := zeroes * params.TxDataZeroGas
	onesCost := ones * params.TxDataNonZeroGasEIP2028
	gasLimit := zeroesCost + onesCost + overhead
	return new(big.Int).SetUint64(gasLimit)
}

func zeroesAndOnes(data []byte) (uint64, uint64) {
	var zeroes uint64
	var ones uint64
	for _, byt := range data {
		if byt == 0 {
			zeroes++
		} else {
			ones++
		}
	}
	return zeroes, ones
}
