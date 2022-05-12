package cosmos

import (
	"testing"
)

type TestAccountCase struct {
	mnemonic string
	cointype int64
	prefix   string
	address  string
}

var accountCase1 = &TestAccountCase{
	mnemonic: "unaware oxygen allow method allow property predict various slice travel please priority",
	cointype: 118,
	prefix:   "cosmos",
	address:  "cosmos19jwusy7lm8v5kqay8qjml79hs6e30t8j7ygm8r",
}
var accountCase2 = &TestAccountCase{
	mnemonic: "wild claw cabin cupboard update cheap thumb blanket float rare change inhale",
	cointype: 118,
	prefix:   "cosmos",
	address:  "cosmos10d2wkfl7y8rpgyxkcwa8urwt8muuc9aqcq9vys",
}
var accountTerra = &TestAccountCase{
	mnemonic: "canyon young easy visa antenna address zone maple captain garden faith crawl tomorrow left risk identify impose miss baby whale nest assume clap trial",
	cointype: 330,
	prefix:   "terra",
	address:  "terra1swy7k7r0jv4rmyjslp35pf0dfp0cs92c8mdwlr",
}

func TestNewAccountWithMnemonic(t *testing.T) {
	errorcase := *accountCase1
	errorcase.mnemonic = errorcase.mnemonic[1:]
	tests := []struct {
		name    string
		acase   TestAccountCase
		wantErr bool
	}{
		{name: "normal case1", acase: *accountCase1},
		{name: "normal case2", acase: *accountCase2},
		{name: "terra normal", acase: *accountTerra},
		{name: "error mnemonic", acase: errorcase, wantErr: true},
		{name: "error empty mnemonic", acase: TestAccountCase{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithMnemonic(tt.acase.mnemonic, tt.acase.cointype, tt.acase.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got.Address() != tt.acase.address {
				t.Errorf("NewAccountWithMnemonic() got = %v, want %v", got, tt.acase.address)
			}
		})
	}
}
