package eth

import (
	"testing"
)

func TestChain_BalanceOfAddress(t *testing.T) {
	tests := []struct {
		name    string
		net     rpcInfo
		address string
		wantErr bool
	}{
		{
			name:    "eth normal",
			net:     rpcs.ethereumProd,
			address: "0x62c3aF16954fba6D920835ec56f7b63139daAa6e",
		},
		{
			name:    "eth black hole",
			net:     rpcs.ethereumProd,
			address: "0x0000000000000000000000000000000000000000",
		},
		{
			name:    "binance-prod normal",
			net:     rpcs.binanceProd,
			address: "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
		},
		{
			name:    "binance-prod error address altered one char", // but is can queryed
			net:     rpcs.binanceProd,
			address: "0x62c3aF16954fba6D920835ec56f7b63139daAa6d",
		},
		{
			name:    "eth error eip55 address", // but is can queryed
			net:     rpcs.ethereumProd,
			address: "0x62c3aF16954fba6D920835ec56f7b63139daAA6E",
		},
		{
			name:    "eth error address short",
			net:     rpcs.ethereumProd,
			address: "0x62c3aF16954fba6D920835ec56f",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.net.Chain().BalanceOfAddress(tt.address)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("BalanceOfAddress() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			t.Log("queryed balance is ", got.Total)
			t.Log("Unable to verify balance, maybe you should check with this address which may be useful: " + tt.net.scan + "/address/" + tt.address)
		})
	}
}
