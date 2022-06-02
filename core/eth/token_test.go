package eth

import "testing"

func TestToken_EstimateGasLimit(t1 *testing.T) {
	addressZero := "0x0000000000000000000000000000000000000000"
	tests := []struct {
		name    string
		rpcInfo rpcInfo
		from    string
		to      string
		amount  string
		wantErr bool
	}{
		{
			name:    "eth prod",
			rpcInfo: rpcs.ethereumProd,
			from:    addressZero,
			to:      "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			amount:  "100",
		},
		{
			name:    "eth prod to invalid address and big amount",
			rpcInfo: rpcs.ethereumProd,
			from:    addressZero,
			to:      "0x7161ada3EA6e53E5652A45988DdfF1",
			amount:  "100000000",
		},
		{
			name:    "eth prod very big amount",
			rpcInfo: rpcs.ethereumProd,
			from:    addressZero,
			to:      "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			amount:  "1000000000000000000000000000000000000000000000000000000000000000000000000000000",
			wantErr: true, // the balance if not enough
		},
		// {
		// 	name:    "binance test",
		// 	rpcInfo: rpcs.binanceTest,
		// 	from:    addressZero,
		// 	to:      "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
		// 	amount:  "100",
		// },
		{
			name:    "sherpax prod",
			rpcInfo: rpcs.sherpaxProd,
			from:    "0xceE683Bb0F4815e1649db4adBC2d5c382Dd5079b",
			to:      "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			amount:  "100",
		},
		{
			name:    "eth error rpc",
			rpcInfo: rpcInfo{url: rpcs.ethereumProd.url + "s"},
			from:    addressZero,
			to:      "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			amount:  "100",
			wantErr: true,
		},
		{
			name:    "optimism prod",
			rpcInfo: rpcs.optimismProd,
			from:    addressZero,
			to:      accountCase1.address,
			amount:  "100000",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			chain := NewChainWithRpc(tt.rpcInfo.url)
			gasPrice, err := chain.SuggestGasPrice()
			if err != nil {
				if !tt.wantErr {
					t1.Errorf("EstimateGasLimit() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			tenfoldPrice := gasPrice.Value + "0"
			token := chain.MainEthToken()
			got, err := token.EstimateGasLimit(tt.from, tt.to, tenfoldPrice, tt.amount)
			if (err != nil) != tt.wantErr {
				t1.Errorf("EstimateGasLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				t1.Logf("EstimateGasLimit() %v", got)
			}
		})
	}
	t1.Log("Note: We can't assert the gas limit.")
}
