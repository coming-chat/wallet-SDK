package multi_signature_check

import (
	"reflect"
	"testing"
)

func TestAddressAmountKey(t *testing.T) {
	type args struct {
		address string
		amount  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddressAmountKey(tt.args.address, tt.args.amount); got != tt.want {
				t.Errorf("AddressAmountKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestJsonRpc(t *testing.T) {
	type args struct {
		url    string
		method string
		params []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RequestJsonRpc(tt.args.url, tt.args.method, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequestJsonRpc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RequestJsonRpc() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXGatewayBitcoinVerifyTxValid(t *testing.T) {
	type args struct {
		url           string
		rawTx         string
		withdrawalIds string
		isFullAmount  bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := XGatewayBitcoinVerifyTxValid(tt.args.url, tt.args.rawTx, tt.args.withdrawalIds, tt.args.isFullAmount)
			if (err != nil) != tt.wantErr {
				t.Errorf("XGatewayBitcoinVerifyTxValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("XGatewayBitcoinVerifyTxValid() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXGatewayCommonWithdrawalListWithFeeInfo(t *testing.T) {
	type args struct {
		url      string
		assertId int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				url:      "https://testnet3.chainx.org/rpc",
				assertId: 1,
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := XGatewayCommonWithdrawalListWithFeeInfo(tt.args.url, tt.args.assertId)
			if (err != nil) != tt.wantErr {
				t.Errorf("XGatewayCommonWithdrawalListWithFeeInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
			//if got != tt.want {
			//	t.Errorf("XGatewayCommonWithdrawalListWithFeeInfo() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
