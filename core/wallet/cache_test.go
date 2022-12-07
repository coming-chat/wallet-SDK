package wallet

import (
	"testing"
	"time"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	{
		wallet := NewCacheWallet(&m1Wallet)
		t.Logf("wallet %v's type = %v", wallet.key(), wallet.WalletType())

		polkaAddress, err := wallet.PolkaAccountInfo(44).Address()
		require.Nil(t, err)
		require.Equal(t, polkaAddress.Value, testcase.Accounts.Polka44.Address)

		aptosAddress, err := wallet.AptosAccountInfo().Address()
		require.Nil(t, err)
		require.Equal(t, aptosAddress.Value, testcase.Accounts.Aptos.Address)

		ethereumAddress, err := wallet.EthereumAccountInfo().Address()
		require.Nil(t, err)
		require.Equal(t, ethereumAddress.Value, testcase.Accounts.Ethereum.Address)

		bitcoinAddress, err := wallet.BitcoinAccountInfo("mainnet").Address()
		require.Nil(t, err)
		require.Equal(t, bitcoinAddress.Value, testcase.Accounts.BtcMainnet.Address)

		dogecoinAddress, err := wallet.DogecoinAccountInfo("mainnet").Address()
		require.Nil(t, err)
		require.Equal(t, dogecoinAddress.Value, testcase.Accounts.DogeMainnet.Address)

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
		wallet := NewCacheWallet(&private_l64Wallet)
		t.Logf("wallet %v's type = %v", wallet.key(), wallet.WalletType())

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
		wallet := NewCacheWallet(&watchWallet)
		t.Logf("wallet %v's type = %v", wallet.key(), wallet.WalletType())

		polkaAddress, err := wallet.PolkaAccountInfo(0).Address()
		require.Nil(t, err)
		require.Equal(t, polkaAddress.Value, watchWallet.watchAddress)

		bitcoinAddress, err := wallet.BitcoinAccountInfo("mainnet").Address()
		require.Nil(t, err)
		require.Equal(t, bitcoinAddress.Value, watchWallet.watchAddress)
	}

	{
		wallet := NewCacheWallet(&emptyWallet)
		t.Logf("wallet %v's type = %v", wallet.key(), wallet.WalletType())

		_, err := wallet.PolkaAccountInfo(2).Account()
		require.Equal(t, err, ErrWalletInfoNotExist)

		_, err = wallet.SuiAccountInfo().Address()
		require.Equal(t, err, ErrWalletInfoNotExist)
	}

	{
		// test m1 again.
		wallet := NewCacheWallet(&m1Wallet)
		t.Logf("wallet %v's type = %v", wallet.key(), wallet.WalletType())

		polkaAddress, err := wallet.PolkaAccountInfo(44).Address()
		require.Nil(t, err)
		require.Equal(t, polkaAddress.Value, testcase.Accounts.Polka44.Address)
	}

}

type WalletStore struct {
	cacheKey     string
	mnemonic     string
	keystore     string
	password     string
	privateKey   string
	watchAddress string
}

func (s *WalletStore) SDKCacheKey() string {
	return s.cacheKey
}
func (s *WalletStore) SDKMnemonic() string {
	return s.mnemonic
}
func (s *WalletStore) SDKKeystore() string {
	return s.keystore
}
func (s *WalletStore) SDKPassword() string {
	return s.password
}
func (s *WalletStore) SDKPrivateKey() string {
	return s.privateKey
}
func (s *WalletStore) SDKWatchAddress() string {
	return s.watchAddress
}

var (
	m1Wallet = WalletStore{
		cacheKey: "m1",
		mnemonic: testcase.M1,
	}
	keystoreWallet = WalletStore{
		cacheKey: "keystore",
		keystore: "TODO",
		password: "TODO",
	}
	watchWallet = WalletStore{
		cacheKey:     "watch",
		watchAddress: "0x33214838821",
	}
	private_l64Wallet = WalletStore{
		cacheKey:   "private_l64",
		privateKey: "0x0000000000000000000000000000000000000000000000000000000000000001",
	}
	emptyWallet = WalletStore{
		cacheKey: "empty",
	}
)
