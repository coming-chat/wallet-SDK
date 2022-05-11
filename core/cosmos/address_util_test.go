package cosmos

import (
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestIsValidAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		prefix  string
		want    bool
	}{
		{
			name:    "cosmos normal",
			address: accountCase1.address,
			prefix:  accountCase1.prefix,
			want:    true,
		},
		{
			name:    "terra normal",
			address: accountTerra.address,
			prefix:  accountTerra.prefix,
			want:    true,
		},
		{
			name:    "cosmos error prefix",
			address: accountCase1.address,
			prefix:  accountTerra.prefix,
			want:    false,
		},
		{
			name:    "cosmos error address",
			address: accountCase1.address + "s",
			prefix:  accountCase1.prefix,
			want:    false,
		},
		{
			name:    "terra error address",
			address: accountTerra.address[1:],
			prefix:  accountTerra.prefix,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidAddress(tt.address, tt.prefix)
			if got != tt.want {
				t.Errorf("IsValidAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivatePublicKeyDataAndStringTransform(t *testing.T) {
	account, _ := NewCosmosAccountWithMnemonic(accountCase1.mnemonic)
	priBytes, _ := account.PrivateKey()
	priHex, _ := account.PrivateKeyHex()
	pubBytes := account.PublicKey()
	pubHex := account.PublicKeyHex()
	address := account.Address()
	t.Log(priBytes)
	t.Log(priHex)
	t.Log(pubBytes)
	t.Log(pubHex)
	t.Log(address)

	t.Log("============================")
	priData, err := types.HexDecodeString(priHex)
	if err != nil {
		t.Fatal(err)
	}
	priKey := secp256k1.PrivKey{Key: priData}
	if priKey.Equals(account.privKey) {
		t.Log("private key restore success!!")
	} else {
		t.Fatal("private key resotre failured.")
	}

	t.Log("============================")
	pubData, err := types.HexDecodeString(pubHex)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := secp256k1.PubKey{Key: pubData}
	if pubKey.Equals(account.privKey.PubKey()) {
		t.Log("public key restore success!!")
	} else {
		t.Fatal("public key resotre failured.")
	}

	t.Log("============================")
	accAddress, err := AccAddressFromBech32(address, accountCase1.prefix)
	originAddress := sdk.AccAddress(account.privKey.PubKey().Address())
	if accAddress.Equals(originAddress) {
		t.Log("address key restore success!!")
	} else {
		t.Fatal("address key resotre failured.")
	}
}
