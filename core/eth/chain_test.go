package eth

import (
	"reflect"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
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

func TestChain_FetchTransactionDetail(t *testing.T) {
	type args struct {
		on   rpcInfo
		hash string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.TransactionDetail
		wantErr bool
	}{
		{
			name: "sherpax-prod contract token USB succeed",
			args: args{rpcs.sherpaxProd, "0x004eaae28f7a7f947c6e8a159f4b74a3aa719248ca4a47d9e5bbf32c394b460f"},
			want: &base.TransactionDetail{
				HashString:      "0x004eaae28f7a7f947c6e8a159f4b74a3aa719248ca4a47d9e5bbf32c394b460f",
				Amount:          "100000000000000000",
				EstimateFees:    "22895100000000000",
				FromAddress:     "0x3AA9c65C4393920E46B2B022841f3EaB7f49f7BC",
				ToAddress:       "0x5fD3d526A946DdB67C810f1F1C4A8c9214da17ef",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1647598998,
			},
		},
		{
			name: "binance-test BNB failured out of gas",
			args: args{rpcs.binanceTest, "0x5841f924fd76434f7f17ef8faf192142dbb5a363b9295eda0cc9f385e22399c7"},
			want: &base.TransactionDetail{
				HashString:      "0x5841f924fd76434f7f17ef8faf192142dbb5a363b9295eda0cc9f385e22399c7",
				Amount:          "1000000000000000000",
				EstimateFees:    "378000000000000",
				FromAddress:     "0xaa25Aa7a19f9c426E07dee59b12f944f4d9f1DD3",
				ToAddress:       "0x6cd2Bf22B3CeaDfF6B8C226487265d81164396C5",
				Status:          base.TransactionStatusFailure,
				FinishTimestamp: 1649840877,
				FailureMessage:  "out of gas",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := tt.args.on.Chain()
			got, err := chain.FetchTransactionDetail(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchTransactionDetail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChain_FetchTransactionDetail_Cover_Multi_Rpcs(t *testing.T) {
	type args struct {
		on   rpcInfo
		hash string
	}
	tests := []struct {
		name     string
		args     args
		wantTime int64 // Pass if timestamp is correct
		wantErr  bool
	}{
		{
			name:     "ethereum contract USDT succeed",
			args:     args{rpcs.ethereumProd, "0xdde5debd0856f295f85f6ee9ae13c84519fb01ef02cb85fa0e5194b618857898"},
			wantTime: 1650816611,
		},
		{
			name:     "sherpax-test KSX success",
			args:     args{rpcs.sherpaxTest, "0xf8ce72bbdfaa51a94bfe34a8bb8f973e35a0e4927711d9b247fb519ad5fb36ed"},
			wantTime: 1650872358,
		},
		{
			name:     "sherpax-prod contract token USB succeed",
			args:     args{rpcs.sherpaxProd, "0x004eaae28f7a7f947c6e8a159f4b74a3aa719248ca4a47d9e5bbf32c394b460f"},
			wantTime: 1647598998,
		},
		{
			name:     "binance-test BNB failured out of gas",
			args:     args{rpcs.binanceTest, "0x5841f924fd76434f7f17ef8faf192142dbb5a363b9295eda0cc9f385e22399c7"},
			wantTime: 1649840877,
		},
		{
			name:     "binance-prod contract USDT succeed",
			args:     args{rpcs.binanceProd, "0xdb3fd18da286f1a2fa90ec9b224b0eb0163a2ab98b9e701da66dc3d8d1d38c14"},
			wantTime: 1651217162,
		},
		{
			name:    "ethereum EIP1559 transfer succeed, but not yet supported",
			args:    args{rpcs.ethereumProd, "0xbba7921fda55dae8423d105d74acdbfea346f15969f6c924552d0b062128f271"},
			wantErr: true,
		},
		{
			name:    "binance-prod error hash",
			args:    args{rpcs.binanceProd, "0x5841f924fd76434f7f17ef8faf192142dbb5a363b9295eda0cc9f385e22399c6"},
			wantErr: true,
		},
		{
			name:    "ethereum-prod error hash",
			args:    args{rpcs.binanceProd, "0x5841f924fd76434f7f17ef8faf19214"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := tt.args.on.Chain()
			got, err := chain.FetchTransactionDetail(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got.FinishTimestamp != tt.wantTime {
				t.Errorf("FetchTransactionDetail() got = %v, want %v", got, tt.wantTime)
			}
		})
	}
}

func TestChain_CallContract(t *testing.T) {
	msg := NewCallMsg()
	msg.SetTo("0x37088186089c7d6bcd556d9a15087dfae3ba0c32")
	msg.SetDataHex("0x70a082310000000000000000000000008de5ff2eded4d897da535ab0f379ec1b9257ebab")

	tests := []struct {
		name    string
		on      rpcInfo
		msg     *CallMsg
		block   int64
		wantErr bool
	}{
		{
			name:  "binance prod pending",
			on:    rpcs.binanceProd,
			msg:   msg,
			block: -1,
		},
		{
			name:  "binance prod latest",
			on:    rpcs.binanceProd,
			msg:   msg,
			block: -2,
		},
		{
			name:    "sherpax prod pending",
			on:      rpcs.sherpaxProd,
			msg:     msg,
			block:   -1,
			wantErr: true, // sherpax not support pending call
		},
		{
			name:  "sherpax prod latest",
			on:    rpcs.sherpaxProd,
			msg:   msg,
			block: -2,
		},
		{
			name:  "sherpax prod normal block",
			on:    rpcs.sherpaxProd,
			msg:   msg,
			block: 1670400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := tt.on.Chain()
			got, err := chain.CallContract(tt.msg, tt.block)
			if (err != nil) != tt.wantErr {
				t.Errorf("CallContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				t.Logf("CallContract() got = %v", got)
			}
		})
	}
}
