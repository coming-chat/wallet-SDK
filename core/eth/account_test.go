package eth

import (
	"testing"
)

type TestAccountCase struct {
	mnemonic string
	address  string
}

var accountCase1 = &TestAccountCase{
	mnemonic: "unaware oxygen allow method allow property predict various slice travel please priority",
	address:  "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
}
var accountCase2 = &TestAccountCase{
	mnemonic: "police saddle quote salon run split notice taxi expand uniform zone excess",
	address:  "0xD32D26054099DbB5A14387d0cF15Df4452EFE4a9",
}
var errorCase = &TestAccountCase{mnemonic: "unaware oxygen allow method allow property predict various slice travel please wrong"}

func TestNewAccountWithMnemonic(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		address  string // If the generated address can match, there is no problem.
		wantErr  bool
	}{
		{name: "valid account 1", mnemonic: accountCase1.mnemonic, address: accountCase1.address},
		{name: "valid account 2", mnemonic: accountCase2.mnemonic, address: accountCase2.address},
		{name: "error mnemonic", mnemonic: errorCase.mnemonic, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithMnemonic(tt.mnemonic)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got.Address() != tt.address {
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
			name:    "valid address case1",
			address: accountCase1.address,
			want:    true,
		},
		{
			name:    "valid address no 0x",
			address: "7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
			want:    true,
		},
		{
			name:    "valid address all caps",
			address: "0x52908400098527886E0F7030069857D2E4169EE7",
			want:    true,
		},
		{
			name:    "valid address all lower",
			address: "0x27b1fdb04752bbc536007a920d24acb045561c26",
			want:    true,
		},
		{
			name:    "invalid address of eip55, alter a char",
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
