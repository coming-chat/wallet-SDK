package eth

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

/// Test the chain `linea`

func LineaMainnetChain() *Chain {
	// scan https://lineascan.build
	return NewChainWithRpc("https://rpc.linea.build")
}
func LineaTestnetChain() *Chain {
	// scan https://goerli.lineascan.build
	return NewChainWithRpc("https://rpc.goerli.linea.build")
}

func TestLineaBalance(t *testing.T) {
	owner := "0x8687B640546744b5338AA292fBbE881162dd5bAe"

	chain := LineaTestnetChain()
	balance, err := chain.BalanceOfAddress(owner)
	require.Nil(t, err)
	t.Log(balance.Total)
}

func TestLineaErc20TokenBalance(t *testing.T) {
	tokenAddress := "0x83240E55e35147B095e8958103a4fd4B32700a3C"
	chain := LineaTestnetChain()
	token := chain.Erc20Token(tokenAddress)

	tokenInfo, err := token.TokenInfo()
	require.Nil(t, err)
	t.Log(tokenInfo)

	owner := "0x422f72B27819798986F41c1bede24e76114DE584"
	balance, err := token.BalanceOfAddress(owner)
	require.Nil(t, err)
	t.Log(balance)
}

func TestLineaTransfer(t *testing.T) {
	sender, err := NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)
	require.Equal(t, sender.Address(), "0x8687B640546744b5338AA292fBbE881162dd5bAe")

	toAddress := sender.Address()
	amount := ETH(0.01).String()

	chain := LineaTestnetChain()

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

func TestLineaFetchTransactionDetail(t *testing.T) {
	hash := "0x6736f02ab900694e324d963365c0200e8131f1a5d0c547c0e23a5628d5bbb3bd"

	chain := LineaTestnetChain()

	detail, err := chain.FetchTransactionDetail(hash)
	require.Nil(t, err)
	require.Equal(t, detail.FromAddress, "0x8687B640546744b5338AA292fBbE881162dd5bAe")
	require.Equal(t, detail.Amount, ETH(0.01).String())
	require.Equal(t, detail.FinishTimestamp, int64(1696836892))

	t.Log(detail.JsonString())
}
