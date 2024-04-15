package btc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"reflect"
	"testing"
)

func TestDerivation(t *testing.T) {
	type args struct {
		mnemonic string
		path     string
		network  *chaincfg.Params
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				mnemonic: "path version decrease world crawl prefer horror version spare water deputy piece",
				path:     "m/84'/0'/0'/0/0",
				network:  &chaincfg.TestNet3Params,
			},
			want:    "cPxupEEx8ir3ZERTpz73knTvmRRbhJCo6Z5BBW6Mk8erDL2uxn1d",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Derivation(tt.args.mnemonic, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Derivation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			wif, err := PrivateKeyToWIF(got, tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyToWIF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(wif, tt.want) {
				t.Errorf("Derivation() got = %v, want %v", got, tt.want)
			}
		})
	}
}
