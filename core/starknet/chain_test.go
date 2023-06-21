package starknet

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func MainnetChain() *Chain {
	return NewChainWithRpc(BaseRpcUrlMainnet, NetworkMainnet)
}
func GoerliChain() *Chain {
	return NewChainWithRpc(BaseRpcUrlGoerli, NetworkGoerli)
}

func TestBalance(t *testing.T) {
	owner := "0x0023C4475F2f2355580f5994294997d3A18237ef62223D20C41876556327A05E"
	chain := GoerliChain()

	balance, err := chain.BalanceOf(owner, ETHTokenAddress)
	require.Nil(t, err)
	t.Log(balance.Total)
}

func TestDeployAccount(t *testing.T) {
	acc := M1Account(t)
	chain := GoerliChain()

	txn, err := chain.BuildDeployAccountTransaction(acc.PublicKeyHex())
	require.Nil(t, err)

	signedTxn, err := txn.SignedTransactionWithAccount(acc)
	require.Nil(t, err)

	hash, err := chain.SendSignedTransaction(signedTxn)
	require.Nil(t, err)

	t.Log(hash.Value)
}

func TestTransfer(t *testing.T) {
	acc := M1Account(t)
	chain := GoerliChain()

	token, err := chain.NewToken(ETHTokenAddress)
	require.Nil(t, err)

	txn, err := token.BuildTransfer(acc.Address(), acc.Address(), "10000000")
	require.Nil(t, err)

	gasFee, err := chain.EstimateTransactionFeeUseAccount(txn, acc)
	require.Nil(t, err)
	t.Log(gasFee.Value)

	// signedTxn, err := txn.SignedTransactionWithAccount(acc)
	// require.Nil(t, err)
	// hash, err := chain.SendSignedTransaction(signedTxn)
	// require.Nil(t, err)
	// t.Log(hash.Value)
}

func TestNonce(t *testing.T) {
	account := M1Account(t)
	chain := GoerliChain()

	address := account.Address()

	nonce, err := chain.gw.Nonce(context.Background(), address, "latest")
	require.Nil(t, err)
	t.Log(nonce.String())
}
