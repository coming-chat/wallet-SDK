package wallet

import "testing"

func TestGenMnemonic(t *testing.T) {
	mnemonic, err := GenMnemonic()
	if err != nil {
		t.Error(err)
	}
	t.Logf("mnemonic: %s", mnemonic)
}
