package eth

import "testing"

func TestToken_EstimateGasLimit(t1 *testing.T) {
	type args struct {
		receiver string
		amount   string
	}
	defaultArgs := args{"0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2", "100"}

	tests := []struct {
		name    string
		rpcInfo rpcInfo
		from    string
		args    args
		wantErr bool
	}{

		{
			name:    "eth USDT",
			rpcInfo: rpcs.ethereumProd,
			from:    "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			args:    defaultArgs,
			wantErr: true, // there is no eth balance, it's will got an error: insufficient funds for transfer
		},
		{
			name:    "binance test USDC",
			rpcInfo: rpcs.binanceTest,
			from:    "0x9a576ec81b75ab1a00baeb976441e34db23882fe",
			args:    defaultArgs,
		},
		{
			name:    "sherpax prod USB",
			rpcInfo: rpcs.sherpaxProd,
			from:    "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			args:    defaultArgs,
			wantErr: true,
		},
		{
			name:    "eth error rpc",
			rpcInfo: rpcInfo{url: rpcs.ethereumProd.url + "s"},
			args:    defaultArgs,
			from:    "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			chain := NewChainWithRpc(tt.rpcInfo.url)
			token := chain.MainEthToken()
			gasPrice, err := chain.SuggestGasPrice()
			if err != nil {
				if !tt.wantErr {
					t1.Errorf("EstimateGasLimit() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			got, err := token.EstimateGasLimit(tt.from, tt.args.receiver, gasPrice.Value, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t1.Errorf("EstimateGasLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				t1.Log(got)
			}
		})
	}
	t1.Log("Note: We can't assert the gas limit.")
}
