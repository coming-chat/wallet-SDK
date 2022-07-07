package solana

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
)

func TestAccount(t *testing.T) {
	mnemonic := testcase.M1
	acc, _ := NewAccountWithMnemonic(mnemonic)
	t.Log(acc.PrivateKeyHex())
	t.Log(acc.PublicKeyHex())
	t.Log(acc.Address())

	prihex, _ := acc.PrivateKeyHex()
	acc2, _ := NewAccountWithPrivateKey(prihex)
	t.Log(acc2.PrivateKeyHex())
	t.Log(acc2.PublicKeyHex())
	t.Log(acc2.Address())
}

func TestValidAddress(t *testing.T) {
	// AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqPY
	addr := "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r1yT3vCxqPy"

	b := IsValidAddress(addr)
	t.Log(b)
}
