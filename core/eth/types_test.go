package eth

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUtils_EncodeErc721TransferFrom(t *testing.T) {
	sender := "0x151e446ca01b57e495a31d53bc622ac33bd7a0be"
	receiver := "0x2c32bd5f7d3eab4bc9d968c90c82debb1bdcced9"
	nftId := "1"

	// Reference: https://scan-canary-testnet.bevm.io/tx/0x2c763fc26b1021340edc5614c7411a8f4d2220d22fecd805557fef4536268ef8
	data, err := EncodeErc721TransferFrom(sender, receiver, nftId)
	require.Nil(t, err)
	dataHex := hex.EncodeToString(data)
	require.Equal(t, dataHex,
		"23b872dd000000000000000000000000151e446ca01b57e495a31d53bc622ac33bd7a0be0000000000000000000000002c32bd5f7d3eab4bc9d968c90c82debb1bdcced90000000000000000000000000000000000000000000000000000000000000001")
}

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
