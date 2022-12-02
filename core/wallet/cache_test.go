package wallet

import (
	"testing"
	"time"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	InfoProvider = &FakeWalletInfoProvider{}

	{
		wallet := NewCacheWallet("m1")
		t.Logf("wallet %v's type = %v", wallet.WalletId, wallet.WalletType())

		polkaAddress, err := wallet.PolkaAccountInfo(44).Address()
		require.Nil(t, err)
		require.Equal(t, polkaAddress.Value, testcase.Accounts.Polka44.Address)

		aptosAddress, err := wallet.AptosAccountInfo().Address()
		require.Nil(t, err)
		require.Equal(t, aptosAddress.Value, testcase.Accounts.Aptos.Address)

		ethereumAddress, err := wallet.EthereumAccountInfo().Address()
		require.Nil(t, err)
		require.Equal(t, ethereumAddress.Value, testcase.Accounts.Ethereum.Address)

		timeStart := time.Now()
		times := 10000
		for i := 0; i < times; i++ {
			ethereumAddress, err := wallet.EthereumAccountInfo().Address()
			require.Nil(t, err)
			require.Equal(t, ethereumAddress.Value, testcase.Accounts.Ethereum.Address)
		}
		timeSpent := time.Since(timeStart)
		t.Logf("%v times get address spent = %v", times, timeSpent)
	}

	{
		wallet := NewCacheWallet("private_l64")
		t.Logf("wallet %v's type = %v", wallet.WalletId, wallet.WalletType())

		polkaAccount, err := wallet.PolkaAccountInfo(44).Account()
		require.Nil(t, err)
		t.Log(polkaAccount.Address())

		ethAccount, err := wallet.EthereumAccountInfo().Account()
		require.Nil(t, err)
		t.Log(ethAccount.Address())

		aptos, err := wallet.AptosAccountInfo().Account()
		require.Nil(t, err)
		t.Log(aptos.Address())
	}

	{
		wallet := NewCacheWallet("watch")
		t.Logf("wallet %v's type = %v", wallet.WalletId, wallet.WalletType())

		polkaAddress, err := wallet.PolkaAccountInfo(0).Address()
		require.Nil(t, err)
		require.Equal(t, polkaAddress.Value, "0x33214838821")

		bitcoinAddress, err := wallet.BitcoinAccountInfo("mainnet").Address()
		require.Nil(t, err)
		require.Equal(t, bitcoinAddress.Value, "0x33214838821")
	}

	{
		wallet := NewCacheWallet("invalid wallet id")
		t.Logf("wallet %v's type = %v", wallet.WalletId, wallet.WalletType())

		_, err := wallet.PolkaAccountInfo(2).Account()
		require.Equal(t, err, ErrWalletInfoNotExist)

		_, err = wallet.SuiAccountInfo().Address()
		require.Equal(t, err, ErrWalletInfoNotExist)
	}

}

type WalletStore struct {
	mnemonic     string
	keystore     string
	password     string
	privateKey   string
	watchAddress string
}

var wallets = map[string]WalletStore{
	"m1":          {mnemonic: testcase.M1},
	"keystore":    {keystore: "TODO", password: "TODO"},
	"watch":       {watchAddress: "0x33214838821"},
	"private_l64": {privateKey: "0x0000000000000000000000000000000000000000000000000000000000000001"},
}

type FakeWalletInfoProvider struct {
}

func (f *FakeWalletInfoProvider) Mnemonic(walletId string) string {
	return wallets[walletId].mnemonic
}
func (f *FakeWalletInfoProvider) Keystore(walletId string) string {
	return wallets[walletId].keystore
}
func (f *FakeWalletInfoProvider) Password(walletId string) string {
	return wallets[walletId].password
}
func (f *FakeWalletInfoProvider) PrivateKey(walletId string) string {
	return wallets[walletId].privateKey
}
func (f *FakeWalletInfoProvider) WatchAddress(walletId string) string {
	return wallets[walletId].watchAddress
}
