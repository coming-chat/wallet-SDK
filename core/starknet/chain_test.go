package starknet

import (
	"context"
	"testing"

	"github.com/dontpanicdao/caigo/gateway"
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

func TestFetchTransactionDetail(t *testing.T) {
	chain := GoerliChain()
	hash := "0x01de50b64326c02a9830ea7bf825224103dd3ea4426309514039a01eaadcf5a4"

	detail, err := chain.FetchTransactionDetail(hash)
	require.Nil(t, err)
	t.Log(detail)
}

func TestTransactionInfo(t *testing.T) {
	chain := GoerliChain()

	hash := "0x47ba4a447e929987094289c27ecc3d37b5a02e580835083d87247c6a97a4e00"

	txn, err := chain.gw.Transaction(context.Background(), gateway.TransactionOptions{
		TransactionHash: hash,
	})
	require.Nil(t, err)

	block, err := chain.gw.Block(context.Background(), &gateway.BlockOptions{
		BlockHash: txn.BlockHash,
	})
	require.Nil(t, err)

	receipt, err := chain.gw.TransactionReceipt(context.Background(), hash)
	require.Nil(t, err)

	status, err := chain.gw.TransactionStatus(context.Background(), gateway.TransactionStatusOptions{
		TransactionHash: hash,
	})
	require.Nil(t, err)

	t.Log(txn)
	t.Log(block)
	t.Log(receipt)
	t.Log(status)
}

func TestFetchTransactionStatus(t *testing.T) {
	chain := GoerliChain()
	hash := "0x03ae12fb58a3f4a6dcd7d04ad10c4d3b2ab97d23ee167a6109db719ba703eed9"

	status := chain.FetchTransactionStatus(hash)
	t.Log(status)
}
