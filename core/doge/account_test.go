package doge

import "testing"

func TestDoge(t *testing.T) {
	mnemonic := "unaware oxygen allow method allow property predict various slice travel please priority"

	account, err := NewAccountWithMnemonic(mnemonic)
	if err != nil {
		t.Log(err)
	}
	t.Log(account.address, err)
}
