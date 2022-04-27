package eth

import (
	"testing"
)

func TestNewAccountWithMnemonic(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		address  string // If the generated address can match, there is no problem.
		wantErr  bool
	}{
		{
			name:     "valid account 1",
			mnemonic: "unaware oxygen allow method allow property predict various slice travel please priority",
			address:  "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
		},
		{
			name:     "error mnemonic",
			mnemonic: "unaware oxygen allow method allow property predict various slice travel please wrong",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithMnemonic(tt.mnemonic)
			if err != nil {
				if tt.wantErr {
					t.Log(tt.name, ": get a error", err)
				} else {
					t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.Address() != tt.address {
				t.Errorf("NewAccountWithMnemonic() got = %v, want %v", got, tt.address)
			}
		})
	}
}

func TestIsValidAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "valid address",
			address: "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			want:    true,
		},
		{
			name:    "valid address no 0x",
			address: "7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			want:    true,
		},
		{
			name:    "invalid address of eip55",
			address: "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c1",
			want:    false,
		},
		{
			name:    "invalid address short length",
			address: "0x7161ada3EA6e53E5652A4",
			want:    false,
		},
		{
			name:    "empty address",
			address: "",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidAddress(tt.address); got != tt.want {
				t.Errorf("IsValidAddress(%v) = %v, want %v", tt.address, got, tt.want)
			}
		})
	}
}
