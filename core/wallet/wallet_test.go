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

	t.Log("Polkadot")
	account, _ = wal.GetOrCreatePolkaAccount(0)
	showAccount(t, account)

	t.Log("Kusama")
	account, _ = wal.GetOrCreatePolkaAccount(2)
	showAccount(t, account)

	t.Log("Ethereum")
	account, _ = wal.GetOrCreateEthereumAccount()
	showAccount(t, account)

	t.Log("Bitcoin")
	account, _ = wal.GetOrCreateBitcoinAccount("mainnet")
	showAccount(t, account)

	t.Log("Bitcoin signet")
	account, _ = wal.GetOrCreateBitcoinAccount("signet")
	showAccount(t, account)

	t.Log("Dogecoin")
	account, _ = wal.GetOrCreateDogeAccount("mainnet")
	showAccount(t, account)

	t.Log("Cosmos")
	account, _ = wal.GetOrCreateCosmosAccount()
	showAccount(t, account)

	t.Log("Terra")
	account, _ = wal.GetOrCreateCosmosTypeAccount(330, "terra")
	showAccount(t, account)

	t.Log("Solana")
	account, _ = wal.GetOrCreateSolanaAccount()
	showAccount(t, account)

	t.Log("Aptos")
	account, _ = wal.GetOrCreateAptosAccount()
	showAccount(t, account)

	t.Log("Sui")
	account, _ = wal.GetOrCreateSuiAccount()
	showAccount(t, account)

	t.Log("Starcoin")
	account, _ = wal.GetOrCreateStarcoinAccount()
	showAccount(t, account)
}

func showAccount(t *testing.T, acc base.Account) {
	t.Log(acc.PrivateKeyHex())
	t.Log("public", acc.PublicKeyHex())
	t.Log("address", acc.Address(), "\n")
}
