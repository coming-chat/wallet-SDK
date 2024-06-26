package wallet

import (
	"testing"
	"time"

	"github.com/coming-chat/wallet-SDK/core/btc"
	"github.com/coming-chat/wallet-SDK/core/cosmos"
	"github.com/coming-chat/wallet-SDK/core/doge"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	{
		wallet := NewCacheWallet(&m1Wallet)
		t.Logf("wallet %v's type = %v", wallet.key(), wallet.WalletType())

		polkaAddress, err := wallet.PolkaAccountInfo(44).Address()
		require.Nil(t, err)
		require.Equal(t, polkaAddress.Value, "5Rdo22BWmvxoFRgLcny1unddon7VErbF1NMhXG1ukskgvnix")

		aptosAddress, err := wallet.AptosAccountInfo().Address()
		require.Nil(t, err)
		require.Equal(t, aptosAddress.Value, "0x11dd2037a613716fdc7cdbd96390b6450bce6754e46b9251cd3c8cd7733683bd")

		ethereumAddress, err := wallet.EthereumAccountInfo().Address()
		require.Nil(t, err)
		require.Equal(t, ethereumAddress.Value, "0x8687B640546744b5338AA292fBbE881162dd5bAe")

		bitcoinAddress, err := wallet.BitcoinAccountInfo("mainnet", 0).Address()
		require.Nil(t, err)
		require.Equal(t, bitcoinAddress.Value, "bc1pn8rv6lqakzlz7yvflc24qjv0pjx8w0plcfq34jl6scs7pqcac2wqjymurp")

		bitcoinAddress2, err := wallet.BitcoinAccountInfo("mainnet", 2).Address()
		require.Nil(t, err)
		require.Equal(t, bitcoinAddress2.Value, "3CZw6pGHD8VVqAhhJZFQKXLLkNkkhqTRHS")

		dogecoinAddress, err := wallet.DogecoinAccountInfo("mainnet").Address()
		require.Nil(t, err)
		require.Equal(t, dogecoinAddress.Value, "DNEJHtF6qjTsWH9muEH5BoQbAmUkyvvDC9")

		timeStart := time.Now()
		times := 10000
		for i := 0; i < times; i++ {
			ethereumAddress, err := wallet.EthereumAccountInfo().Address()
			require.Nil(t, err)
			require.Equal(t, ethereumAddress.Value, "0x8687B640546744b5338AA292fBbE881162dd5bAe")
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

		bitcoinAddress, err := wallet.BitcoinAccountInfo("mainnet", 0).Address()
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
		require.Equal(t, polkaAddress.Value, "5Rdo22BWmvxoFRgLcny1unddon7VErbF1NMhXG1ukskgvnix")
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

func TestDecimalPrivatekey(t *testing.T) {
	decimalKey := "2522809042406563759994430227158123451351879727850354505141300412651234567890"
	argWallet := WalletStore{
		cacheKey:   "decimal_privatekey",
		privateKey: decimalKey,
	}

	wallet := NewCacheWallet(&argWallet)
	t.Log(wallet.PolkaAccountInfo(44).Address())
	t.Log(wallet.BitcoinAccountInfo(btc.ChainMainnet, btc.AddressTypeComingTaproot).Address())
	t.Log(wallet.EthereumAccountInfo().Address())
	t.Log(wallet.CosmosAccountInfo(cosmos.CosmosAtom.Cointype, cosmos.CosmosAtom.Prefix).Address())
	t.Log(wallet.DogecoinAccountInfo(doge.ChainMainnet).Address())
	t.Log(wallet.SolanaAccountInfo().Address())
	t.Log(wallet.AptosAccountInfo().Address())
	t.Log(wallet.SuiAccountInfo().Address())
	t.Log(wallet.StarcoinAccountInfo().Address())
	t.Log(wallet.StarknetAccountInfo(false).Address())
}
