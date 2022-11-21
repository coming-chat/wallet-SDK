package doge

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

type TestAccountCase struct {
	mnemonic    string
	privateKey  string
	publicKey   string
	addrMainnet string
	addrTestnet string
}

var accountCase = &TestAccountCase{
	mnemonic:    "unaware oxygen allow method allow property predict various slice travel please priority",
	privateKey:  "0xc7fceb75bafba7aa10ffe10315352bfc523ac733f814e6a311bc736873df8923",
	publicKey:   "0x04a721f170043daafde0fa925ab6caf5d2abcdadd2249291b1840e3d99a3f41149e13185ef52451eef2e7cc0c5fe4180b64ca2d17eb886b2328518f6aed684719a",
	addrMainnet: "DJhF8ahvTfGhqcLEn7sN4gJMJVVbmfwxkU",
	addrTestnet: "nhkJrbSqPdjRiauRowWpK5teYMstkMp4M6",
}
var errorCase = &TestAccountCase{
	mnemonic: "unaware oxygen allow method allow property predict various slice travel please check",
}

const (
	// https://shibe.technology/
	returnAddress = "nbMFaHF9pjNoohS4fD1jefKBgDnETK9uPu"
)

func TestDoge(t *testing.T) {
	account, err := NewAccountWithMnemonic(accountCase.mnemonic, ChainTestnet)
	require.Nil(t, err)
	t.Log(account.PrivateKeyHex())
	t.Log(account.PublicKeyHex())
	t.Log(account.Address())
}

func TestAccountWithPrivatekey(t *testing.T) {
	mnemonic := testcase.M1
	accountFromMnemonic, err := NewAccountWithMnemonic(mnemonic, ChainMainnet)
	require.Nil(t, err)
	privateKey, err := accountFromMnemonic.PrivateKeyHex()
	require.Nil(t, err)

	accountFromPrikey, err := AccountWithPrivateKey(privateKey, ChainMainnet)
	require.Nil(t, err)

	require.Equal(t, accountFromMnemonic.Address(), accountFromPrikey.Address())
}
