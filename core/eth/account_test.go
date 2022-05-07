package eth

import (
	"crypto/ecdsa"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
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

func TestETHWallet_Privatekey_Publickey_Address(t *testing.T) {
	// 从 coming 的 trust wallet 库计算的测试用例
	// private key = 0x8c3083c24062f065ff2ee71b21f665375b266cebffa920e8909ec7c48006725d
	// public key  = 0xc66cbe3908fda67d2fb229b13a63aa1a2d8428acef2ff67bc31f6a79f2e2085f // Curve25519
	// public key  = 0xb34ec4ec2ebc84b04d9170bed91f65306c7045863efb9175d721104a8ecc17f2 // Ed25519
	// public key  = 0x011e56a004e205db53ae3cc7291ffb8a28181aed4b4e95813c17b9a96db2d769 // Ed25519Blake2b
	// public key  = 0x04bd6d7af856d20188fcfdb8ff38b978bc7c72fd028b67a6fab3d2120dd9bd1db61c5d44e242001dce224188a8b88150e16e9748438703bbf2dc417135c4f9377e // Secp256k1 compressed false
	// public key  = 0x02bd6d7af856d20188fcfdb8ff38b978bc7c72fd028b67a6fab3d2120dd9bd1db6 // Secp256k1 compressed true
	// public key  = 0x027bcb5a6edf262eca9602b8343baa1cd5dd7811e540e850b05661b6524e504222 // Nist256p1
	// address     = 0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2

	phrase := "unaware oxygen allow method allow property predict various slice travel please priority"
	seed, _ := bip39.NewSeedWithErrorChecking(phrase, "")

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}

	path, _ := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	key := masterKey
	for _, n := range path {
		key, err = key.DeriveNonStandard(n)
		if err != nil {
			t.Fatal(err)
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		t.Fatal(err)
	}
	privateKeyECDSA := privateKey.ToECDSA()
	privateKeyHex := types.HexEncodeToString(privateKey.Serialize())
	t.Log("private key = ", privateKeyHex)

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal(".....")
	}

	data := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyHex := types.HexEncodeToString(data)
	t.Log("public key = ", publicKeyHex)

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	t.Log("address = ", address.Hex())
}
