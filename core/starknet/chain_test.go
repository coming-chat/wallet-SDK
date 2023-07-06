package starknet

import (
	"context"
	"testing"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/stretchr/testify/require"
)

func MainnetChain() *Chain {
	c, _ := NewChainWithRpc(BaseRpcUrlMainnet, NetworkMainnet)
	return c
}
func GoerliChain() *Chain {
	c, _ := NewChainWithRpc(BaseRpcUrlGoerli, NetworkGoerli)
	return c
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

	txn, err := chain.BuildDeployAccountTransaction(acc.PublicKeyHex(), "")
	require.Nil(t, err)

	signedTxn, err := txn.SignedTransactionWithAccount(acc)
	require.Nil(t, err)

	hash, err := chain.SendSignedTransaction(signedTxn)
	require.Nil(t, err)

	t.Log(hash.Value)
}

func TestChain_IsContractAddressDeployed(t *testing.T) {
	addr := "0x63242861a734490bf31412bcb84a6ad37e370c99a5697de6dd3e8c2ebd40539"
	chain := GoerliChain()

	deployed, err := chain.IsContractAddressDeployed(addr)
	require.Nil(t, err)
	require.True(t, deployed.Value)

	deployed2, err := MainnetChain().IsContractAddressDeployed(addr)
	require.Nil(t, err)
	require.False(t, deployed2.Value)
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
	// xx := signedTxn.(*SignedTransaction)
	// xx.NeedAutoDeploy = true
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
	hash := "0x06bb5eaa8861e3bf95c61bd6723c758333b2675d9d80fee032f1aefef1bf9cbd"

	detail, err := chain.FetchTransactionDetail(hash)
	require.Nil(t, err)
	t.Log(detail)
}

func TestFetchTransactionDetail_Mainnet(t *testing.T) {
	chain := MainnetChain()
	hash := "0x0415fb8adcafec90b89bb8dadb0e8a5968ecba0040d4c560b1847bde9f392954"

	detail, err := chain.FetchTransactionDetail(hash)
	require.Nil(t, err)
	t.Log(detail)
}

func TestTransactionInfo(t *testing.T) {
	chain := GoerliChain()

	hash := "0x01982805b13bb3d661c0015df15c210a6166cbad720e554dbf62f2106102c849"

	txn, err := chain.gw.Transaction(context.Background(), gateway.TransactionOptions{
		TransactionHash: hash,
	})
	require.Nil(t, err)

	id, err := chain.gw.TransactionID(context.Background(), hash)
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

	t.Log(id)
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

func TestNotdeployedAccount(t *testing.T) {
	acc, err := AccountWithPrivateKey("0x123456fffff")
	require.Nil(t, err)

	chain := GoerliChain()
	token := chain.MainToken()

	txn, err := token.BuildTransfer(acc.Address(), acc.Address(), "100000000")
	require.Nil(t, err)

	_, err = chain.EstimateTransactionFeeUseAccount(txn, acc)
	t.Log(IsNotDeployedError(err))

	signedTxn, err := txn.SignedTransactionWithAccount(acc)
	require.Nil(t, err)

	_, err = chain.SendSignedTransaction(signedTxn)
	t.Log(IsNotDeployedError(err))
}
