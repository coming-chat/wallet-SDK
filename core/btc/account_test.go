package btc

import (
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
	"github.com/tyler-smith/go-bip39"
)

type TestAccountCase struct {
	mnemonic    string
	privateKey  string
	publicKey   string
	addrMainnet string
	addrSignet  string
}

var accountCase = &TestAccountCase{
	mnemonic:    "unaware oxygen allow method allow property predict various slice travel please priority",
	privateKey:  "0xc7fceb75bafba7aa10ffe10315352bfc523ac733f814e6a311bc736873df8923",
	publicKey:   "0x04a721f170043daafde0fa925ab6caf5d2abcdadd2249291b1840e3d99a3f41149e13185ef52451eef2e7cc0c5fe4180b64ca2d17eb886b2328518f6aed684719a",
	addrMainnet: "bc1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ysunz3sg",
	addrSignet:  "tb1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ystm5728",
}
var errorCase = &TestAccountCase{
	mnemonic: "unaware oxygen allow method allow property predict various slice travel please check",
}

func TestNewAccountWithMnemonic(t *testing.T) {
	type args struct {
		mnemonic string
		chainnet string
	}
	tests := []struct {
		name        string
		args        args
		wantAddress string // If the generated address can match, there is no problem.
		wantErr     bool
	}{
		{
			name:        "mainnet nomal",
			args:        args{mnemonic: accountCase.mnemonic, chainnet: ChainMainnet},
			wantAddress: accountCase.addrMainnet,
		},
		{
			name:        "coming bitcoin nomal",
			args:        args{mnemonic: accountCase.mnemonic, chainnet: ChainBitcoin},
			wantAddress: accountCase.addrMainnet,
		},
		{
			name:        "signet nomal",
			args:        args{mnemonic: accountCase.mnemonic, chainnet: ChainSignet},
			wantAddress: accountCase.addrSignet,
		},
		{
			name:    "error chainnet",
			args:    args{mnemonic: accountCase.mnemonic, chainnet: "xxxxxxx"},
			wantErr: true,
		},
		{
			name:    "error mnemonic",
			args:    args{mnemonic: errorCase.mnemonic, chainnet: ChainSignet},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithMnemonic(tt.args.mnemonic, tt.args.chainnet)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got.address != tt.wantAddress {
				t.Errorf("NewAccountWithMnemonic() got = %v, want %v", got.address, tt.wantAddress)
			}
		})
	}
}

func TestAccount_DeriveAccountAt(t *testing.T) {
	baseAccount, err := NewAccountWithMnemonic(accountCase.mnemonic, ChainMainnet)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name        string
		chainnet    string
		wantAddress string
		wantErr     bool
	}{
		{name: "same as mainnet", chainnet: ChainMainnet, wantAddress: accountCase.addrMainnet},
		{name: "change signet", chainnet: ChainSignet, wantAddress: accountCase.addrSignet},
		{name: "error net", chainnet: "signet2", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := baseAccount.DeriveAccountAt(tt.chainnet)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeriveAccountAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got.address != tt.wantAddress {
				t.Errorf("DeriveAccountAt() got = %v, want %v", got.address, tt.wantAddress)
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
		{name: "normal test", mnemonic: accountCase.mnemonic, want: accountCase.privateKey},
		{name: "invalid mnemonic", mnemonic: errorCase.mnemonic, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			account, err := NewAccountWithMnemonic(tt.mnemonic, ChainMainnet)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("PrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			got, err := account.PrivateKeyHex()
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

func TestAccountWithPrivateKey(t *testing.T) {
	acc, err := NewAccountWithMnemonic(testcase.M1, ChainMainnet)
	require.Nil(t, err)
	privateKey, err := acc.PrivateKeyHex()
	require.Nil(t, err)

	acc2, err := AccountWithPrivateKey(privateKey, ChainMainnet)
	require.Nil(t, err)
	require.Equal(t, acc.privateKey, acc2.privateKey)
	require.Equal(t, acc.address, acc2.address)
}

func TestIsValidAddress(t *testing.T) {
	type args struct {
		chainnet string
		address  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "mainnet valid",
			args: args{chainnet: ChainMainnet, address: accountCase.addrMainnet},
			want: true,
		},
		{
			name: "signet valid",
			args: args{chainnet: ChainSignet, address: accountCase.addrSignet},
			want: true,
		},
		{
			name: "mainnet valid check in signet",
			args: args{chainnet: ChainSignet, address: accountCase.addrMainnet},
			want: true,
		},
		{
			name: "signet valid check in mainnet",
			args: args{chainnet: ChainMainnet, address: accountCase.addrSignet},
			want: true,
		},
		{
			name: "error address",
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

func TestBTCWallet_Privatekey_Publickey_Address(t *testing.T) {
	// 从 coming 的 musig 库计算的测试用例
	// private key = 0xc7fceb75bafba7aa10ffe10315352bfc523ac733f814e6a311bc736873df8923
	// public key = 0x04a721f170043daafde0fa925ab6caf5d2abcdadd2249291b1840e3d99a3f41149e13185ef52451eef2e7cc0c5fe4180b64ca2d17eb886b2328518f6aed684719a
	// mainnet address = bc1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ysunz3sg
	// signet address = tb1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ystm5728

	phrase := "unaware oxygen allow method allow property predict various slice travel please priority"
	data, _ := bip39.NewSeedWithErrorChecking(phrase, "")

	pri, pub := btcec.PrivKeyFromBytes(data)
	priHex := types.HexEncodeToString(pri.Serialize())
	pubHex := types.HexEncodeToString(pub.SerializeUncompressed())
	t.Log("private key = ", priHex)
	t.Log("public key = ", pubHex)

	pubData := pub.SerializeUncompressed()
	addressHash, err := btcutil.NewAddressTaproot(pubData[1:33], &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("mainnet address = ", addressHash.EncodeAddress())

	addressHash, err = btcutil.NewAddressTaproot(pubData[1:33], &chaincfg.SigNetParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("signet address = ", addressHash.EncodeAddress())
}

func TestAccountWithPrivatekey(t *testing.T) {
	mnemonic := testcase.M1
	accountFromMnemonic, err := NewAccountWithMnemonic(mnemonic, ChainMainnet)
	require.Nil(t, err)
	privateKey, err := accountFromMnemonic.PrivateKeyHex()
	require.Nil(t, err)

	accountFromPrikey, err := AccountWithPrivateKey(privateKey, ChainMainnet)
	require.Nil(t, err)

	require.Equal(t, accountFromMnemonic.Address(), accountFromPrikey.Address())
}
