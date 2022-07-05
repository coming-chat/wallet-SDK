package eth

import (
	"reflect"
	"testing"
)

func TestNewTransactionFromHex(t *testing.T) {
	type args struct {
		hexData string
	}
	tests := []struct {
		name    string
		args    args
		want    *Transaction
		wantErr bool
	}{
		{
			name: "EIP1559 tx",
			args: args{"02f8708081dd84861c468084938580c882b437941717a0d5c8705ee89a8ad6e808268d6a826c97a480b844095ea7b30000000000000000000000007e88c5e7134e4589f6316636ca8fe8cc9f8ed50500000000000000000000000000000000000000000000000000000000000f4240c0808080"},
			want: &Transaction{
				Nonce:                "221",
				GasPrice:             "2475000008",
				GasLimit:             "46135",
				To:                   "0x1717A0D5C8705EE89A8aD6E808268D6A826C97A4",
				Value:                "0",
				Data:                 "095ea7b30000000000000000000000007e88c5e7134e4589f6316636ca8fe8cc9f8ed50500000000000000000000000000000000000000000000000000000000000f4240",
				MaxPriorityFeePerGas: "2250000000",
			},
			wantErr: false,
		},
		{
			name: "Legancy tx",
			args: args{"f86981dd84938580c882b437941717a0d5c8705ee89a8ad6e808268d6a826c97a480b844095ea7b30000000000000000000000007e88c5e7134e4589f6316636ca8fe8cc9f8ed5050000000000000000000000000000000000000000000000000000000005f5e100808080"},
			want: &Transaction{
				Nonce:                "221",
				GasPrice:             "2475000008",
				GasLimit:             "46135",
				To:                   "0x1717A0D5C8705EE89A8aD6E808268D6A826C97A4",
				Value:                "0",
				Data:                 "095ea7b30000000000000000000000007e88c5e7134e4589f6316636ca8fe8cc9f8ed5050000000000000000000000000000000000000000000000000000000005f5e100",
				MaxPriorityFeePerGas: "",
			},
			wantErr: false,
		},
		{
			name:    "error example",
			args:    args{"580c882b437941717a0d5c8705ee89a8ad6e808268d6a826c97a480b844095ea7b30000000000000000000000007e88c5e7134e4589f6316636ca8fe8cc9f8ed5050000000000000000000000000000000000000000000000000000000005f5e100808080"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransactionFromHex(tt.args.hexData)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTransactionFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionFromHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
