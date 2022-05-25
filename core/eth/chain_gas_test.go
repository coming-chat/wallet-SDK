package eth

import (
	"testing"
)

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
			if got.BaseFee == "" && got.PriorityFeeLow == "" {
				t.Errorf("SuggestGasPriceEIP1559() got an empty fee = %v", got)
			} else {
				t.Logf("SuggestGasPriceEIP1559() got %v, you can cheker at %v", got, tt.checker)
			}
		})
	}
}
