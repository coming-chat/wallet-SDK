package wallet

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestAccounts(t *testing.T) {
	m := testcase.M1
	m = "only admit tumble endorse swear argue copy frozen favorite climb obscure palm"
	wal, err := NewWalletWithMnemonic(m)
	require.Nil(t, err)

	var account base.Account

	t.Log("Chainx, Minix, Sherpax")
	account, _ = wal.GetOrCreatePolkaAccount(44)
	showAccount(t, account)
	require.Equal(t, account.Address(), "5Ud8ArhmuPnwaF61HCS2QhsC4yUxttaxriXdUSf3XMLmy5Ge")

	t.Log("Polkadot")
	account, _ = wal.GetOrCreatePolkaAccount(0)
	showAccount(t, account)
	require.Equal(t, account.Address(), "15z6oCt2zbJxhfHbRmiyo4AVddC2Upt4ZL9rQNMo9jBKaoGh")

	t.Log("Kusama")
	account, _ = wal.GetOrCreatePolkaAccount(2)
	showAccount(t, account)
	require.Equal(t, account.Address(), "HZRKBxqmB4R1n6XEqV2YrhLvbUcbC96wDG7djeQ5SNJ9VS4")

	t.Log("Ethereum")
	account, _ = wal.GetOrCreateEthereumAccount()
	showAccount(t, account)
	require.Equal(t, account.Address(), "0x4be5b6c8657dAe87031B6fF1906A08953d4204E5")

	t.Log("Bitcoin")
	account, _ = wal.GetOrCreateBitcoinAccount("mainnet", 0)
	showAccount(t, account)
	require.Equal(t, account.Address(), "bc1p5s8866n4h959679ylqdpkhcthld6y6dhp0phru5eaqpwxnfdxnaqp9g9jl")

	t.Log("Bitcoin signet")
	account, _ = wal.GetOrCreateBitcoinAccount("signet", 0)
	showAccount(t, account)
	require.Equal(t, account.Address(), "tb1p5s8866n4h959679ylqdpkhcthld6y6dhp0phru5eaqpwxnfdxnaqkd72gs")

	t.Log("Dogecoin")
	account, _ = wal.GetOrCreateDogeAccount("mainnet")
	showAccount(t, account)
	require.Equal(t, account.Address(), "DJRP7zgrd4h26TQPpuA3vQM8REjxh6TYMk")

	t.Log("Cosmos")
	account, _ = wal.GetOrCreateCosmosAccount()
	showAccount(t, account)
	require.Equal(t, account.Address(), "cosmos1jxylpm0twp7zgj3pk3qvww640nv4cppw6rp3vx")

	t.Log("Terra")
	account, _ = wal.GetOrCreateCosmosTypeAccount(330, "terra")
	showAccount(t, account)
	require.Equal(t, account.Address(), "terra1qnjrfufhmm72s2dtave86uz8quexrp2nkfqpy0")

	t.Log("Solana")
	account, _ = wal.GetOrCreateSolanaAccount()
	showAccount(t, account)
	require.Equal(t, account.Address(), "AEygNEH37EeHMq8MH19B1ZzrHuy3wQEfmszykYecf3Dt")

	t.Log("Aptos")
	account, _ = wal.GetOrCreateAptosAccount()
	showAccount(t, account)
	require.Equal(t, account.Address(), "0xb34f8adce502d60bff9d04d568977e50fd7720a6ee59de605df6261152ca09c0")

	t.Log("Sui")
	account, _ = wal.GetOrCreateSuiAccount()
	showAccount(t, account)
	require.Equal(t, account.Address(), "0x109986847ea978ebab6ed15604d10b91e772e412")

	t.Log("Starcoin")
	account, _ = wal.GetOrCreateStarcoinAccount()
	showAccount(t, account)
	require.Equal(t, account.Address(), "0x551c986364613236Cf1810D2C70c83b9")
}

func showAccount(t *testing.T, acc base.Account) {
	t.Log(acc.PrivateKeyHex())
	t.Log("public", acc.PublicKeyHex())
	t.Log("address", acc.Address(), "\n")
}
