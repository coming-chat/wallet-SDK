package btc

import (
	"testing"
)

func TestNewAccountWithMnemonic(t *testing.T) {
	type args struct {
		mnemonic string
		chainnet string
	}
	mnemonic := "unaware oxygen allow method allow property predict various slice travel please priority"
	// privateKey := "0xc7fceb75bafba7aa10ffe10315352bfc523ac733f814e6a311bc736873df8923"
	// publicKey := "0x04a721f170043daafde0fa925ab6caf5d2abcdadd2249291b1840e3d99a3f41149e13185ef52451eef2e7cc0c5fe4180b64ca2d17eb886b2328518f6aed684719a"
	addressMainnet := "bc1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ysunz3sg"
	addressSignet := "tb1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ystm5728"
	tests := []struct {
		name        string
		args        args
		wantAddress string // If the generated address can match, there is no problem.
		wantErr     bool
	}{
		{
			name:        "mainnet nomal",
			args:        args{mnemonic: mnemonic, chainnet: ChainMainnet},
			wantAddress: addressMainnet,
		},
		{
			name:        "coming bitcoin nomal",
			args:        args{mnemonic: mnemonic, chainnet: ChainBitcoin},
			wantAddress: addressMainnet,
		},
		{
			name:        "signet nomal",
			args:        args{mnemonic: mnemonic, chainnet: ChainSignet},
			wantAddress: addressSignet,
		},
		{
			name: "error chainnet",
			args: args{
				mnemonic: mnemonic,
				chainnet: "xxxxxxx",
			},
			wantErr: true,
		},
		{
			name: "error mnemonic",
			args: args{
				mnemonic: "unaware oxygen allow method allow property predict various slice travel please fake",
				chainnet: ChainSignet,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithMnemonic(tt.args.mnemonic, tt.args.chainnet)
			if err != nil {
				if tt.wantErr {
					t.Log(tt.name, ": get a error", err)
				} else {
					t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.address != tt.wantAddress {
				t.Errorf("NewAccountWithMnemonic() got = %v, want %v", got, tt.wantAddress)
			}
		})
	}
}

func TestIsValidAddress(t *testing.T) {
	type args struct {
		chainnet string
		address  string
	}
	addressMainnet := "bc1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ysunz3sg"
	addressSignet := "tb1p4fwg0qlcsm94y90gnkwr0zkfsv9gxjlq43mpegf4cmn9xed02xcq3n0386"
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "mainnet valid",
			args: args{chainnet: ChainMainnet, address: addressMainnet},
			want: true,
		},
		{
			name: "signet valid",
			args: args{chainnet: ChainSignet, address: addressSignet},
			want: true,
		},
		{
			name: "mainnet valid check in signet",
			args: args{chainnet: ChainSignet, address: addressMainnet},
			want: true,
		},
		{
			name: "signet valid check in mainnet",
			args: args{chainnet: ChainMainnet, address: addressSignet},
			want: true,
		},
		{
			name: "mainnet invalid",
			args: args{chainnet: ChainMainnet, address: "bc1p5uslzuqy8k40mc86jfdtdjh4624umtw"},
			want: false,
		},
		{
			name: "empty address",
			args: args{chainnet: ChainSignet, address: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidAddress(tt.args.address, tt.args.chainnet); got != tt.want {
				t.Errorf("IsValidAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
