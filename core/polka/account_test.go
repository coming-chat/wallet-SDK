package polka

import (
	"testing"
)

func TestNewAccountWithMnemonic(t *testing.T) {
	mnemonic := "unaware oxygen allow method allow property predict various slice travel please priority"
	chainxAddress := "5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnH"
	polkaAddress := "12jrfZLTddDxRQAjoSkWurDyEPxPdkhPcgU2AGxFHbgBpyHZ"
	kusamaAddress := "EKBBYRGQCyQjWyfcWWZfekpXNEyk7xRzZaHPeErDJsAPeiD"
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
			args:    args{mnemonic: mnemonic, network: 44},
			address: chainxAddress,
		},
		{
			name:    "polkadot test",
			args:    args{mnemonic: mnemonic, network: 0},
			address: polkaAddress,
		},
		{
			name:    "kusama test",
			args:    args{mnemonic: mnemonic, network: 2},
			address: kusamaAddress,
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
			if err != nil {
				if tt.wantErr {
					t.Log(tt.name, ": get a error", err)
				} else {
					t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.Address() != tt.address {
				t.Errorf("NewAccountWithMnemonic() got = %v, want %v", got.Address(), tt.address)
			}
		})
	}
}
