package sui

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/assert"
)

func TestAccount(t *testing.T) {
	account, err := NewAccountWithMnemonic(testcase.M1)
	assert.Nil(t, err)
	assert.Equal(t, account.Address(), "0x6c5d2cd6e62734f61b4e318e58cbfd1c4b99dfaf")

	t.Log(account.PrivateKeyHex())
	t.Log(account.PublicKeyHex())
	t.Log(account.Address())
}

func TestPublicKeyToAddress(t *testing.T) {
	pub := "0x1cec19ef9a036d27a055e8ad49e8c37cdc16ab2fb3270b73424a971af9039604"
	addr, err := EncodePublicKeyToAddress(pub)
	assert.Nil(t, err)
	assert.Equal(t, addr, "0x0bd43fc3aa4f62e8943d16f66beb7546fafb2bac")
}

// Account of os environment M1
func M1Account() *Account {
	account, _ := NewAccountWithMnemonic(testcase.M1)
	return account
}

// Account of chrome wallet extension
func ChromeAccount() *Account {
	m := "crack coil okay hotel glue embark all employ east impact stomach cigar"
	account, _ := NewAccountWithMnemonic(m)
	return account
}
