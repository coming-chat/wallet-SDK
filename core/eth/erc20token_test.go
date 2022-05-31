package eth

import (
	"reflect"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
)

func TestBatch(t *testing.T) {
	c := NewChainWithRpc(rpcs.binanceTest.url)
	array := []string{
		// "0xE7e312dfC08e060cda1AF38C234AEAcc7A982143", // 报错
		// "0x4B53739D798EF0BEa5607c254336b40a93c75b52", // 报错
		// "0x935CC842f220CF3A7D10DA1c99F01B1A6894F7C5", // 报错
		// "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56", // 报错
		"0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee", // 可请求余额
	}

	address := "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"
	result, err := c.BatchErc20TokenBalance(array, address)
	if err != nil {
		t.Fatal(err)
	}

	for i, x := range result {
		t.Log(i, x)
	}
}

func TestSdkBatch(t *testing.T) {
	e, _ := NewEthChain().CreateRemote(rpcs.binanceTest.url)
	contract := "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"
	address := "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"
	result, err := e.SdkBatchTokenBalance(contract, address)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}

func TestErc20Token_TokenInfo(t1 *testing.T) {
	tests := []struct {
		name     string
		rpc      string
		contract string
		want     *base.TokenInfo
		wantErr  bool
	}{
		{
			name:     "eth USDT",
			rpc:      rpcs.ethereumProd.url,
			contract: rpcs.ethereumProd.contracts.USDT,
			want:     &base.TokenInfo{Name: "Tether USD", Symbol: "USDT", Decimal: 6},
		},
		{
			name:     "binance prod USDC",
			rpc:      rpcs.binanceProd.url,
			contract: rpcs.binanceProd.contracts.USDC,
			want:     &base.TokenInfo{Name: "USD Coin", Symbol: "USDC", Decimal: 18},
		},
		{
			name:     "sherpax test BUSD",
			rpc:      rpcs.sherpaxTest.url,
			contract: rpcs.sherpaxTest.contracts.BUSD,
			want:     &base.TokenInfo{Name: "Binance BUSD Token", Symbol: "BUSD", Decimal: 18},
		},
		{
			name:     "eth error rpc",
			rpc:      rpcs.ethereumProd.url + "s",
			contract: rpcs.ethereumProd.contracts.USDT,
			wantErr:  true,
		},
		{
			name:     "eth error contract",
			rpc:      rpcs.ethereumProd.url,
			contract: rpcs.sherpaxTest.contracts.BUSD, // not a eth contract address
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			chain := NewChainWithRpc(tt.rpc)
			token := chain.Erc20Token(tt.contract)
			got, err := token.TokenInfo()
			if (err != nil) != tt.wantErr {
				t1.Errorf("TokenInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("TokenInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErc20Token_BalanceOfAddress(t1 *testing.T) {
	tests := []struct {
		name     string
		rpcInfo  rpcInfo
		contract string
		address  string
		wantErr  bool
	}{
		{
			name:     "eth USDT",
			rpcInfo:  rpcs.ethereumProd,
			contract: rpcs.ethereumProd.contracts.USDT,
			address:  "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
		},
		{
			name:     "binance test USDC",
			rpcInfo:  rpcs.binanceTest,
			contract: rpcs.binanceTest.contracts.BUSD,
			address:  "0x9a576ec81b75ab1a00baeb976441e34db23882fe",
		},
		{
			name:     "sherpax prod USB",
			rpcInfo:  rpcs.sherpaxProd,
			contract: rpcs.sherpaxProd.contracts.USB,
			address:  "0x66E96c1238Bfb68C9D88C85F8C4b5FFb44a65472",
		},
		{
			name:     "eth error rpc",
			rpcInfo:  rpcInfo{url: rpcs.ethereumProd.url + "s"},
			contract: rpcs.ethereumProd.contracts.USDT,
			wantErr:  true,
		},
		{
			name:     "sherpax prod error contract",
			rpcInfo:  rpcs.sherpaxProd,
			contract: rpcs.ethereumProd.contracts.USDT,
			address:  "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			chain := NewChainWithRpc(tt.rpcInfo.url)
			token := chain.Erc20Token(tt.contract)
			got, err := token.BalanceOfAddress(tt.address)
			if (err != nil) != tt.wantErr {
				t1.Errorf("BalanceOfAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				scanAddress := tt.rpcInfo.scan + "/token/" + tt.contract + "?a=" + tt.address
				t1.Log(got.Total)
				t1.Log("We can't assert the balance, Maybe you can query on scan url: ", scanAddress)
			}
		})
	}
}

func TestErc20Token_EstimateGasLimit(t1 *testing.T) {
	enoughGasPrice := "100000000000" // 100 Gwei
	haveEthUsdtAddress := "0x22fFF189C37302C02635322911c3B64f80CE7203"

	tests := []struct {
		name     string
		rpcInfo  rpcInfo
		contract string
		from     string
		to       string
		amount   string
		wantErr  bool
	}{

		{
			name:     "eth USDT",
			rpcInfo:  rpcs.ethereumProd,
			contract: rpcs.ethereumProd.contracts.USDT,
			from:     haveEthUsdtAddress,
			to:       "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			amount:   "100",
		},
		{
			name:     "binance test USDC",
			rpcInfo:  rpcs.binanceTest,
			contract: rpcs.binanceTest.contracts.BUSD,
			from:     "0x7Da8a0276627fa857f5459f4B1A9D8161226d604",
			to:       "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			amount:   "100",
		},
		{
			name:     "sherpax prod USB",
			rpcInfo:  rpcs.sherpaxProd,
			contract: rpcs.sherpaxProd.contracts.USB,
			from:     "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			to:       "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			amount:   "100",
			wantErr:  true, // ther is no balance.
		},
		{
			name:     "eth error rpc",
			rpcInfo:  rpcInfo{url: rpcs.ethereumProd.url + "s"},
			contract: rpcs.ethereumProd.contracts.USDT,
			from:     "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			to:       "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			amount:   "100",
			wantErr:  true,
		},
		{
			name:     "sherpax prod error contract",
			rpcInfo:  rpcs.sherpaxProd,
			contract: rpcs.ethereumProd.contracts.USDT,
			from:     "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			to:       "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			amount:   "100",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			chain := NewChainWithRpc(tt.rpcInfo.url)
			token := chain.Erc20Token(tt.contract)
			got, err := token.EstimateGasLimit(tt.from, tt.to, enoughGasPrice, tt.amount)
			if (err != nil) != tt.wantErr {
				t1.Errorf("EstimateGasLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				t1.Log(got)
			}
		})
	}
}
