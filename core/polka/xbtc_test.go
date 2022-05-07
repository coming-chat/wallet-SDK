package polka

import (
	"testing"
)

func TestXBTCToken_BalanceOfAddress(t1 *testing.T) {
	tests := []struct {
		name    string
		rpcInfo rpcInfo
		address string
		wantErr bool
	}{
		{
			name:    "prod empty",
			rpcInfo: rpcs.chainxProd,
			address: accountCase.address44,
		},
		{
			name:    "prod normal",
			rpcInfo: rpcs.chainxProd,
			address: "5QUEkPdwqUCdnTb7C89ay46kp7zKN4gwbaJ4cEF7caE58MKn",
		},
		{
			name:    "test normal",
			rpcInfo: rpcs.chainxTest,
			address: "5Qeua99vLDXrv9THQyM3Rtcq95oadoXfktFsAhVggiwxTNfZ",
		},
		{
			name:    "use XBTC at error chain",
			rpcInfo: rpcs.sherpaxProd,
			address: accountCase.address44,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			chain, _ := tt.rpcInfo.Chain()
			t := chain.XBTCToken()
			got, err := t.BalanceOfAddress(tt.address)
			if (err != nil) != tt.wantErr {
				t1.Errorf("BalanceOfAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				url := tt.rpcInfo.realScan + "/account/" + tt.address
				t1.Log("BalanceOfAddress() result: ", got, ", Maybe you should verify via the link: ", url)
			}
		})
	}
}
