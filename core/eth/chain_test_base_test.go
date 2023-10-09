package eth

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

/// Test the chain `base`

func BaseMainnetChain() *Chain {
	// scan https://basescan.org/
	return NewChainWithRpc("https://mainnet.base.org")
}
func BaseTestnetChain() *Chain {
	// scan https://goerli.basescan.org/
	return NewChainWithRpc("https://goerli.base.org")
}

func TestBaseBalance(t *testing.T) {
	owner := "0x8687B640546744b5338AA292fBbE881162dd5bAe"
	chain := BaseTestnetChain()
	balance, err := chain.BalanceOfAddress(owner)
	require.Nil(t, err)
	t.Log(balance)
}

func TestBaseErc20TokenBalance(t *testing.T) {
	tokenAddress := "0xDd1351a0d3BB3c9D4516aEd7adCcf8814c7A193B"
	chain := BaseTestnetChain()
	token := chain.Erc20Token(tokenAddress)

	tokenInfo, err := token.TokenInfo()
	require.Nil(t, err)
	t.Log(tokenInfo)

	owner := "0x14acba2BAB926C6BFb64239C120C466424217477"
	balance, err := token.BalanceOfAddress(owner)
	require.Nil(t, err)
	t.Log(balance)
}

func TestBaseTransfer(t *testing.T) {
	sender, err := NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)
	require.Equal(t, sender.Address(), "0x8687B640546744b5338AA292fBbE881162dd5bAe")

	toAddress := sender.Address()
	amount := ETH(0.01).String()

	chain := BaseTestnetChain()

	gasPrice, err := chain.SuggestGasPrice()
	require.Nil(t, err)

	token := chain.MainEthToken()
	gasLimit, err := token.EstimateGasLimit(sender.Address(), toAddress, gasPrice.Value, amount)
	require.Nil(t, err)

	transaction := NewTransaction("", gasPrice.Value, gasLimit, toAddress, amount, "")
	signedTx, err := token.BuildTransferTxWithAccount(sender, transaction)
	require.Nil(t, err)

	if false {
		hash, err := chain.SendRawTransaction(signedTx.Value)
		require.Nil(t, err)

		t.Log("txn send success, hash = ", hash)
	}
}

func TestBaseFetchTransactionDetail(t *testing.T) {
	hash := "0x6a86326feaada152b9c125c04904122bb8bdd1c07357c384f4818ebe3bf91f8f"

	chain := BaseTestnetChain()

	detail, err := chain.FetchTransactionDetail(hash)
	require.Nil(t, err)
	require.Equal(t, detail.FromAddress, "0x8687B640546744b5338AA292fBbE881162dd5bAe")
	require.Equal(t, detail.Amount, ETH(0.01).String())
	require.Equal(t, detail.FinishTimestamp, int64(1696824474))

	t.Log(detail.JsonString())
}
