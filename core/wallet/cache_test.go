package wallet

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestSSSss(t *testing.T) {
	InfoProvider = &FakeWalletInfoProvider{}

	wallet := Wallet{
		WalletId: "",
	}

	info := wallet.PolkaAccountInfo(44)

	account, err := info.Account()
	require.Nil(t, err)

	t.Log(account.Address())
	t.Log(info.Address())

	info = wallet.EthereumAccountInfo()
	t.Log(info.Address())
	t.Log(info.PublicKeyHex())
}

type FakeWalletInfoProvider struct {
}

func (f *FakeWalletInfoProvider) Mnemonic(walletId string) string {
	return testcase.M1
}
func (f *FakeWalletInfoProvider) Keystore(walletId string) string {
	return ""
}
func (f *FakeWalletInfoProvider) Password(walletId string) string {
	return ""
}
func (f *FakeWalletInfoProvider) PrivateKey(walletId string) string {
	return ""
}
