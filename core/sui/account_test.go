package sui

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	account := M1Account(t)
	require.Equal(t, account.Address(), "0x6c5d2cd6e62734f61b4e318e58cbfd1c4b99dfaf")

	t.Log(account.PrivateKeyHex())
	t.Log(account.PublicKeyHex())
	t.Log(account.Address())
}

func TestPublicKeyToAddress(t *testing.T) {
	pub := "0x1cec19ef9a036d27a055e8ad49e8c37cdc16ab2fb3270b73424a971af9039604"
	addr, err := EncodePublicKeyToAddress(pub)
	require.Nil(t, err)
	require.Equal(t, addr, "0x0bd43fc3aa4f62e8943d16f66beb7546fafb2bac")
}

// Account of os environment M1
func M1Account(t *testing.T) *Account {
	account, err := NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)
	return account
}

func M2Account(t *testing.T) *Account {
	account, err := NewAccountWithMnemonic(testcase.M2)
	require.Nil(t, err)
	return account
}
