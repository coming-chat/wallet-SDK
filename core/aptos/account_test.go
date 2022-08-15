package aptos

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
)

func TestAccount(t *testing.T) {
	mnemonic := testcase.M1
	account, err := NewAccountWithMnemonic(mnemonic)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(account.PrivateKeyHex())
	t.Log(account.PublicKeyHex())
	t.Log(account.Address())

	prihex, _ := account.PrivateKeyHex()
	acc2, err := AccountWithPrivateKey(prihex)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(acc2.PrivateKeyHex())
	t.Log(acc2.PublicKeyHex())
	t.Log(acc2.Address())
}

func TestIsValidAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			address: "0x1",
			want:    true,
		},
		{
			address: "0x1234567890abcdefABCDEF",
			want:    true,
		},
		{
			address: "0X1234567890123456789012345678901234567890123456789012345678901234",
			want:    true,
		},
		{
			address: "012345aabcdF",
			want:    true,
		},
		{address: "1x23239444"},
		{address: "0x1fg"},
		{address: "0X12345678901234567890123456789012345678901234567890123456789012345"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidAddress(tt.address); got != tt.want {
				t.Errorf("IsValidAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
