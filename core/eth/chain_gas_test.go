package eth

import (
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func TestChain_EstimateGasLimitLayer2(t *testing.T) {
	l2gasPrice := big.NewInt(1e6)

	add := common.HexToAddress("0")
	msg := ethereum.CallMsg{
		From:     add,
		To:       &add,
		Value:    big.NewInt(1e14),
		GasPrice: l2gasPrice,
		Data:     nil,
	}

	chain := rpcs.optimismProd.Chain()
	gas, err := chain.EstimateGasLimitLayer2(&CallMsg{msg: msg})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(gas, gas.GasFee())
}

func TestChain_SuggestGasPriceEIP1559(t *testing.T) {
	tests := []struct {
		name    string
		rpcInfo rpcInfo
		checker string
		wantErr bool
	}{
		{
			name:    "eth prod gas price",
			rpcInfo: rpcs.ethereumProd,
			checker: "https://etherscan.io/gastracker",
		},
		{
			name:    "rinkeby gas price",
			rpcInfo: rpcs.rinkeby,
			checker: "",
		},
		{
			name:    "binance prod not support eip1559 yet",
			rpcInfo: rpcs.binanceProd,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.rpcInfo.Chain()
			got, err := c.SuggestGasPriceEIP1559()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("SuggestGasPriceEIP1559() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.BaseFee == "" || got.SuggestPriorityFee == "" {
				t.Errorf("SuggestGasPriceEIP1559() got an empty fee = %v", got)
			} else {
				t.Logf("SuggestGasPriceEIP1559() got %v, you can cheker at %v", got, tt.checker)
				t.Log(got.UseLow())
				t.Log(got)
				t.Log(got.UseHigh())
			}
		})
	}
}

func TestChain_EstimateGasLimit(t *testing.T) {
	addressZero := &common.Address{}
	e36Wei, _ := big.NewInt(0).SetString("1000000000000000000000000000000000000", 10)
	ethEnoughGasPrice := big.NewInt(100000000000) // 100 Gwei

	addressUSDTContract := common.HexToAddress(rpcs.ethereumProd.contracts.USDT)
	addressUSDTFrom := common.HexToAddress("0x22fFF189C37302C02635322911c3B64f80CE7203")
	usdtTransferData, _ := types.HexDecodeString("0xa9059cbb000000000000000000000000dac17f958d2ee523a2206206994597c13d831ec70000000000000000000000000000000000000000000000000000000000200b20")
	tests := []struct {
		name    string
		rpcInfo rpcInfo
		msg     *CallMsg
		wantErr bool
	}{
		{
			name:    "eth prod transfer very big amount",
			rpcInfo: rpcs.ethereumProd,
			msg: &CallMsg{msg: ethereum.CallMsg{
				From:     *addressZero,
				To:       addressZero,
				Value:    e36Wei, // 10^18 Ether
				GasPrice: ethEnoughGasPrice,
			}},
			wantErr: true,
		},
		{
			name:    "eth prod transfer nomal amount",
			rpcInfo: rpcs.ethereumProd,
			msg: &CallMsg{msg: ethereum.CallMsg{
				From:     *addressZero,
				To:       addressZero,
				Value:    big.NewInt(1000000000), // 1 Gwei
				GasPrice: ethEnoughGasPrice,
			}},
		},
		{
			name:    "eth prod transfer small amount",
			rpcInfo: rpcs.ethereumProd,
			msg: &CallMsg{msg: ethereum.CallMsg{
				From:     *addressZero,
				To:       addressZero,
				Value:    big.NewInt(1), // 1 wei
				GasPrice: ethEnoughGasPrice,
			}},
		},
		{
			name:    "eth prod usdt transfer",
			rpcInfo: rpcs.ethereumProd,
			msg: &CallMsg{msg: ethereum.CallMsg{
				From:     addressUSDTFrom,
				To:       &addressUSDTContract,
				Value:    big.NewInt(0),
				GasPrice: ethEnoughGasPrice,
				Data:     usdtTransferData,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.rpcInfo.Chain()
			got, err := c.EstimateGasLimit(tt.msg)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("EstimateGasLimit() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			t.Log("EstimateGasLimit() estimate gas limit =", got.Value)
		})
	}
}
