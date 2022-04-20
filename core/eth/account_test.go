package eth

import (
	"reflect"
	"testing"
)

func TestNewAccountWithMnemonic(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		want     *Account
		wantErr  bool
	}{
		{
			name:     "valid account 1",
			mnemonic: "unaware oxygen allow method allow property predict various slice travel please priority",
			want: &Account{
				privateKey: "0x8c3083c24062f065ff2ee71b21f665375b266cebffa920e8909ec7c48006725d",
				address:    "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			},
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAccountWithMnemonic() got = %v, want %v", got, tt.want)
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
			name:    "invalid address, but we cant check",
			address: "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c1",
			want:    true,
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
				t.Errorf("IsValidAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
