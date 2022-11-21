package cosmos

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

type TestAccountCase struct {
	mnemonic string
	cointype int64
	prefix   string
	address  string
}

var accountCase1 = &TestAccountCase{
	mnemonic: "unaware oxygen allow method allow property predict various slice travel please priority",
	cointype: CosmosCointype,
	prefix:   CosmosPrefix,
	address:  "cosmos19jwusy7lm8v5kqay8qjml79hs6e30t8j7ygm8r",
}
var accountCase2 = &TestAccountCase{
	mnemonic: "wild claw cabin cupboard update cheap thumb blanket float rare change inhale",
	cointype: CosmosCointype,
	prefix:   CosmosPrefix,
	address:  "cosmos10d2wkfl7y8rpgyxkcwa8urwt8muuc9aqcq9vys",
}
var accountTerra = &TestAccountCase{
	mnemonic: "canyon young easy visa antenna address zone maple captain garden faith crawl tomorrow left risk identify impose miss baby whale nest assume clap trial",
	cointype: TerraCointype,
	prefix:   TerraPrefix,
	address:  "terra1swy7k7r0jv4rmyjslp35pf0dfp0cs92c8mdwlr",
}
var accountTerra2 = &TestAccountCase{
	mnemonic: "chronic crater bronze frown since repeat wonder lazy skull extend view later van copper result fun fantasy unaware author regular dizzy hood swamp sail",
	cointype: TerraCointype,
	prefix:   TerraPrefix,
	address:  "terra1ugmzgw4m89mv887suxn2070k24dhu9xynrxhg8",
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
		{name: "terra normal2", acase: *accountTerra2},
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

func TestAccountPrivatekeyAtDifferenceChain(t *testing.T) {
	mnemonic := testcase.M1

	cosmosAccount, err := NewAccountWithMnemonic(mnemonic, CosmosCointype, CosmosPrefix)
	require.Nil(t, err)
	terraAccount, err := NewAccountWithMnemonic(mnemonic, TerraCointype, TerraPrefix)
	require.Nil(t, err)
	p1, _ := cosmosAccount.PrivateKeyHex()
	p2, _ := terraAccount.PrivateKeyHex()
	t.Logf("cosmos private key = %v", p1)
	t.Logf("terra private key = %v", p2)
	t.Log(cosmosAccount.Address(), "\n", terraAccount.Address())
	// result: same mnemonic generate difference private key at difference cosmos chain
}

func TestAccountWithPrivatekey(t *testing.T) {
	mnemonic := testcase.M1
	accountFromMnemonic, err := NewAccountWithMnemonic(mnemonic, CosmosCointype, CosmosPrefix)
	require.Nil(t, err)
	privateKey, err := accountFromMnemonic.PrivateKeyHex()
	require.Nil(t, err)

	accountFromPrikey, err := AccountWithPrivateKey(privateKey, CosmosCointype, CosmosPrefix)
	require.Nil(t, err)

	require.Equal(t, accountFromMnemonic.Address(), accountFromPrikey.Address())
}
