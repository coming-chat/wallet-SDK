package starknet

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestTokenInfo(t *testing.T) {
	chain := GoerliChain()
	token := chain.MainToken()

	info, err := token.TokenInfo()
	require.Nil(t, err)
	t.Log(info)
}

func TestTokenBalance_Mainnet(t *testing.T) {
	chain := MainnetChain()

	USDTTokenAddress := "0x068f5c6a61780768455de69077e07e89787839bf8166decfbf92b645209c0fb8"
	token, err := NewToken(chain, USDTTokenAddress)
	require.Nil(t, err)

	owner := "0x06f163a072ee7e3894d973c6fd05ada06b1669174ba2aecefbb8b578f4e14003"
	balance, err := token.BalanceOfAddress(owner)
	require.Nil(t, err)
	t.Log(balance.Total)

	info, err := token.TokenInfo()
	require.Nil(t, err)
	t.Log(base.JsonString(info))
}

func TestTokenBalance_Goerli(t *testing.T) {
	chain := GoerliChain()

	BTCTokenAddress := "0x072df4dc5b6c4df72e4288857317caf2ce9da166ab8719ab8306516a2fddfff7"
	token, err := NewToken(chain, BTCTokenAddress)
	require.Nil(t, err)

	owner := "0x02ad6ae0a72c2f083f9e1b33057f8b35c643023c54f41be5b03f807277fcd88c"
	balance, err := token.BalanceOfAddress(owner)
	require.Nil(t, err)
	t.Log(balance.Total)

	info, err := token.TokenInfo()
	require.Nil(t, err)
	t.Log(base.JsonString(info))
}

func TestToken_Transfer(t *testing.T) {
	chain := GoerliChain()

	mn := testcase.M1
	acc, err := NewAccountWithMnemonic(mn)
	require.Nil(t, err)
	t.Log(acc.Address())

	tokenAddr := "0x005a643907b9a4bc6a55e9069c4fd5fd1f5c79a22470690f75556c4736e34426" // usdc
	token, err := chain.NewToken(tokenAddr)
	require.Nil(t, err)
	info, err := token.TokenInfo()
	require.Nil(t, err)
	t.Log(base.JsonString(info))

	balance, err := token.BalanceOfAddress(acc.Address())
	require.Nil(t, err)
	t.Log(balance.Total)
	if len(balance.Total) < 3 {
		return // balance not enough.
	}

	transferAmount := balance.Total[0 : len(balance.Total)-1]
	txn, err := token.BuildTransfer(acc.Address(), acc.Address(), transferAmount)
	require.Nil(t, err)
	signedTxn, err := txn.SignedTransactionWithAccount(acc)
	require.Nil(t, err)
	hash, err := chain.SendSignedTransaction(signedTxn)
	require.Nil(t, err)
	t.Log(hash.Value)
}
