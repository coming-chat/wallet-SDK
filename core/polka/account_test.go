package polka

import (
	"testing"
)

type TestAccountCase struct {
	mnemonic string
	keystore string
	password string

	privateKey string
	publicKey  string
	address0   string
	address2   string
	address44  string
}

var accountCase = &TestAccountCase{
	mnemonic:   "unaware oxygen allow method allow property predict various slice travel please priority",
	privateKey: "0xa6ddfbcac3ff93fbbcff52b064d951e48f1c3828bc5c9d69030f4adbba445f60",
	publicKey:  "0x4cebba1cf615cfefc3bf44117a4a64ec827555fa2a3120b729286b8f7bddc93c",
	address0:   "12jrfZLTddDxRQAjoSkWurDyEPxPdkhPcgU2AGxFHbgBpyHZ",
	address2:   "EKBBYRGQCyQjWyfcWWZfekpXNEyk7xRzZaHPeErDJsAPeiD",
	address44:  "5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnH",
}
var keystoreCase = &TestAccountCase{
	keystore:  "{\"encoded\":\"5zmfXmtpiz8sryDmupYcoFDDCRj0ufe1Fx1EfGFLQoMAgAAAAQAAAAgAAADajJFtVRycQELlG4KibfgTOX4zexng/E3oj+I+ND9GYQIcHnIrEfAu1Ptcoi1HLiM8GfKuzcmMg9ZEvhywWF1Hau4XThv8pk8xGQUyMn2iMQtV8JA/5SGL/w5r5bT9vPOsidQEkc4Q5RvEsqjeU0hCkGKQXIui/9DqFR02Dq9pn3KYK3EQNjkNZplBJ59h4pG+E6SNMG8XuKqDMn+b\",\"encoding\":{\"content\":[\"pkcs8\",\"sr25519\"],\"type\":[\"scrypt\",\"xsalsa20-poly1305\"],\"version\":\"3\"},\"address\":\"5UczqUVGsoQpZnBCZkDtxvLxJ42KnUfaGTzPkQmZeAAug4s9\",\"meta\":{\"genesisHash\":\"0x96675ae0e91fe7d102f8eebc4ee4fbb9241b483bc6645ac975864684d1c222ff\",\"isHardware\":false,\"name\":\"wallet test\",\"tags\":[],\"whenCreated\":1645428018341}}",
	password:  "111",
	address0:  "15yyTpfXxzvqhCNniKWrMGeFrhjPNQxfy5ccgLUKGY1THbTW",
	address44: "5UczqUVGsoQpZnBCZkDtxvLxJ42KnUfaGTzPkQmZeAAug4s9",
}
var oldPrivateKeyCase = &TestAccountCase{
	mnemonic:   "rookie october miracle crisp invest grace birth exile black attitude bitter napkin",
	privateKey: "0xba865d03c9f6f27871d4eddd8baffe2b16c444945388b39adb0a0966020bbbbe",
}

func TestNewAccountWithMnemonic(t *testing.T) {
	type args struct {
		mnemonic string
		network  int
	}
	tests := []struct {
		name    string
		args    args
		address string // If the generated address can match, there is no problem.
		wantErr bool
	}{
		{
			name:    "chainx test",
			args:    args{mnemonic: accountCase.mnemonic, network: 44},
			address: accountCase.address44,
		},
		{
			name:    "polkadot test",
			args:    args{mnemonic: accountCase.mnemonic, network: 0},
			address: accountCase.address0,
		},
		{
			name:    "kusama test",
			args:    args{mnemonic: accountCase.mnemonic, network: 2},
			address: accountCase.address2,
		},
		{
			name:    "error mnemonic",
			args:    args{mnemonic: "unaware oxygen allow method allow property predict various ", network: 44},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithMnemonic(tt.args.mnemonic, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got.Address() != tt.address {
				t.Errorf("NewAccountWithMnemonic() got = %v, want %v", got.Address(), tt.address)
			}
		})
	}
}

func TestNewAccountWithKeystore(t *testing.T) {
	type args struct {
		keystore string
		password string
		network  int
	}
	tests := []struct {
		name        string
		args        args
		wantAddress string
		wantErr     bool
	}{
		{
			name:        "normal case net 0",
			args:        args{keystore: keystoreCase.keystore, password: keystoreCase.password, network: 0},
			wantAddress: keystoreCase.address0,
		},
		{
			name:        "normal case net 44",
			args:        args{keystore: keystoreCase.keystore, password: keystoreCase.password, network: 44},
			wantAddress: keystoreCase.address44,
		},
		{
			name:    "error password",
			args:    args{keystore: keystoreCase.keystore, password: "xxxxx", network: 0},
			wantErr: true,
		},
		{
			name:    "error keystore",
			args:    args{keystore: "{encoded...}", password: keystoreCase.password, network: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithKeystore(tt.args.keystore, tt.args.password, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountWithKeystore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got.Address() != tt.wantAddress {
				t.Errorf("NewAccountWithKeystore() got = %v, want %v", got.Address(), tt.wantAddress)
			}
		})
	}
}

func TestAccount_PrivateKey(t *testing.T) {
	tests := []struct {
		name     string
		mnemonic string
		want     string
		wantErr  bool
	}{
		{name: "valid1", mnemonic: accountCase.mnemonic, want: accountCase.privateKey},
		{name: "old wallet case", mnemonic: oldPrivateKeyCase.mnemonic, want: oldPrivateKeyCase.privateKey},
		{name: "error menmonic", mnemonic: "rookie october miracle crisp invest grace ", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccountWithMnemonic(tt.mnemonic, 44)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			got, err := account.PrivateKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got != tt.want {
				t.Errorf("PrivateKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeAddressToPublicKey(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    string
		wantErr bool
	}{
		{name: "valid address net 0", address: accountCase.address0, want: accountCase.publicKey},
		{name: "valid address net 2", address: accountCase.address2, want: accountCase.publicKey},
		{name: "valid address net 44", address: accountCase.address44, want: accountCase.publicKey},
		{name: "invalid address 44, alter a char", address: "5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnB", wantErr: true},
		{name: "empty address", address: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeAddressToPublicKey(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeAddressToPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got != tt.want {
				t.Errorf("DecodeAddressToPublicKey() got = %v, want %v", got, tt.want)
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
		{name: "valid address net 0", address: accountCase.address0, want: true},
		{name: "valid address net 2", address: accountCase.address2, want: true},
		{name: "valid address net 44", address: accountCase.address44, want: true},
		{name: "invalid address 44, alter a char", address: "5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnA", want: false},
		{name: "empty address", address: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidAddress(tt.address); got != tt.want {
				t.Errorf("IsValidAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
