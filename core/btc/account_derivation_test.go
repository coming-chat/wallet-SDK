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

func TestUnisatCase(t *testing.T) {
	mn := "certain tonight subway hazard parade security manual define maple magnet fix erosion"
	type args struct {
		mnemonic string
		network  string
		addrType AddressType
	}
	tests := []struct {
		name     string
		args     args
		wantAddr string
		wantErr  bool
	}{
		{
			name: "testnet P2WPKH",
			args: args{
				mnemonic: mn,
				network:  ChainTestnet,
				addrType: AddressTypeNativeSegwit,
			},
			wantAddr: "tb1qkevvrpwpydhk023l4eq0tlx4n9xam88j7cl8td",
		},
		{
			name: "testnet P2SH-P2WPKH",
			args: args{
				mnemonic: mn,
				network:  ChainTestnet,
				addrType: AddressTypeNestedSegwit,
			},
			wantAddr: "2MufMqgUTkHVUXnsRoYbQFbE5a7zMhpTHkU",
		},
		{
			name: "testnet P2TR",
			args: args{
				mnemonic: mn,
				network:  ChainTestnet,
				addrType: AddressTypeTaproot,
			},
			wantAddr: "tb1p5w4203pxzyz92glhj53lmwgvgcmmxx80khc9pgctpukx3pyrytaqt2saex",
		},
		{
			name: "testnet P2PKH",
			args: args{
				mnemonic: mn,
				network:  ChainTestnet,
				addrType: AddressTypeLegacy,
			},
			wantAddr: "mki7vqGxPmZcu1Uq3xyKmeoWetpCL9epGP",
		},
		{
			name: "testnet Coming Taproot",
			args: args{
				mnemonic: mn,
				network:  ChainTestnet,
				addrType: AddressTypeComingTaproot,
			},
			wantAddr: "tb1pk6sdtvvhshml36sghlrw55ppvxtp72mrwmgnu8zyw6ha7vxxaufqljj6yk",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, err := NewAccountWithMnemonic(tt.args.mnemonic, tt.args.network, tt.args.addrType)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && acc.Address() != tt.wantAddr {
				t.Errorf("NewAccountWithMnemonic() got = %v, want %v", acc.Address(), tt.wantAddr)
			}
		})
	}
}
