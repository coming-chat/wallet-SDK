package aptos

import (
	"os"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func M1Account(t *testing.T) *Account {
	acc, err := NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)
	return acc
}

var (
	PriMartian1 = os.Getenv("PriMartian1")
	PriMartian2 = os.Getenv("PriMartian2")
	PriPetra1   = os.Getenv("PriPetra1")
)

func TestAccount(t *testing.T) {
	mnemonic := testcase.M1
	account, err := NewAccountWithMnemonic(mnemonic)
	require.Nil(t, err)

	prihex, _ := account.PrivateKeyHex()
	acc2, err := AccountWithPrivateKey(prihex)
	require.Nil(t, err)

	require.Equal(t, account.PublicKey(), acc2.PublicKey())
	require.Equal(t, account.Address(), acc2.Address())
	require.Equal(t, account.Address(), "0x11dd2037a613716fdc7cdbd96390b6450bce6754e46b9251cd3c8cd7733683bd")

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
