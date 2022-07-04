package eth

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func TestChain_EstimateGasLimitLayer2(t *testing.T) {
	address := "0x0000000000000000000000000000000000000000"
	l2gasPrice := strconv.FormatInt(1e6, 10)
	amount := strconv.FormatInt(1e14, 10)

	msg := CallMsg{}
	msg.SetFrom(address)
	msg.SetTo(address)
	msg.SetValue(amount)
	msg.SetGasPrice(l2gasPrice)
	msg.SetGasLimit("21000")

	chain := rpcs.optimismProd.Chain()
	gas, err := chain.MainEthToken().EstimateGasFeeLayer2(&msg)
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
		{
			name:    "avax prod gas price",
			rpcInfo: rpcs.avaxProd,
			checker: "https://snowtrace.io/gastracker",
		},
		{
			name:    "polygon prod",
			rpcInfo: rpcs.polygonProd,
			checker: "https://polygonscan.com/gastracker",
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
	e36Wei, _ := big.NewInt(0).SetString("1000000000000000000000000000000000000", 10)
	tests := []struct {
		name    string
		rpcInfo rpcInfo
		amount  *big.Int
		msg     *CallMsg
		wantErr bool
	}{
		{
			name:    "eth prod transfer very big amount",
			rpcInfo: rpcs.ethereumProd,
			amount:  e36Wei, // 10^18 Ether
			wantErr: true,
		},
		{
			name:    "eth prod transfer nomal amount",
			rpcInfo: rpcs.ethereumProd,
			amount:  big.NewInt(1e9), // 1 Gwei
		},
		{
			name:    "eth prod transfer small amount",
			rpcInfo: rpcs.ethereumProd,
			amount:  big.NewInt(1), // 1 wei
		},
		{
			name:    "avax prod normal",
			rpcInfo: rpcs.avaxProd,
			amount:  big.NewInt(1e9),
		},
		{
			name:    "avax test normal",
			rpcInfo: rpcs.avaxTest,
			amount:  big.NewInt(1e9),
		},
		{
			name:    "polygon prod normal",
			rpcInfo: rpcs.polygonProd,
			amount:  big.NewInt(1e9),
		},
	}
	ethEnoughGasPrice := big.NewInt(1e11) // 100 Gwei
	addressZero := common.Address{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.rpcInfo.Chain()
			msg := &CallMsg{msg: ethereum.CallMsg{
				From:     addressZero,
				To:       &addressZero,
				Value:    tt.amount,
				GasPrice: ethEnoughGasPrice,
			}}
			got, err := c.EstimateGasLimit(msg)
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

func TestChain_EstimateGasLimit_Erc20(t *testing.T) {
	tests := []struct {
		name    string
		rpcInfo rpcInfo
		from    string
		to      string
		wantErr bool
	}{
		{
			name:    "eth prod usdt transfer",
			rpcInfo: rpcs.ethereumProd,
			from:    "0x22fFF189C37302C02635322911c3B64f80CE7203",
			to:      rpcs.ethereumProd.contracts.USDT,
		},
		{
			name:    "avax prod erc20 transfer",
			rpcInfo: rpcs.avaxProd,
			from:    "0xCF08fB3925900bdbc07B9CC8dF87efD7BCD3BdbA",
			to:      "0xB97EF9Ef8734C71904D8002F8b6Bc66Dd9c48a6E",
		},
		{
			name:    "polygon prod erc20 transfer",
			rpcInfo: rpcs.polygonProd,
			from:    "0xbB9480FDB174063a2DeDF31b50A4f06aE13ed964",
			to:      "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
		},
	}

	valueZero := big.NewInt(0)
	ethEnoughGasPrice := big.NewInt(1e11) // 100 Gwei
	erc20TransferData, _ := types.HexDecodeString("0xa9059cbb000000000000000000000000dac17f958d2ee523a2206206994597c13d831ec70000000000000000000000000000000000000000000000000000000000200b20")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.rpcInfo.Chain()
			msg := &CallMsg{msg: ethereum.CallMsg{
				From:     common.HexToAddress(tt.from),
				To:       addressPointer(tt.to),
				Value:    valueZero,
				GasPrice: ethEnoughGasPrice,
				Data:     erc20TransferData,
			}}
			got, err := c.EstimateGasLimit(msg)
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

func addressPointer(hex string) *common.Address {
	a := common.HexToAddress(hex)
	return &a
}
