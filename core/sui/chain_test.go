package sui

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

const DevnetRpcUrl = "https://fullnode.devnet.sui.io"

func DevnetChain() *Chain {
	return NewChainWithRpcUrl(DevnetRpcUrl)
}

func TestTransfer(t *testing.T) {
	account := M1Account(t)
	chain := DevnetChain()
	token := NewTokenMain(chain)

	// toAddress := "0x0c61c2622b77e2a9a3c953690e915ab82d6370d9"
	// amount := "8000000"
	toAddress := account.Address()
	amount := "15000000"

	signedTxn, err := token.BuildTransferTxWithAccount(account, toAddress, amount)
	require.Nil(t, err)

	hash, err := chain.SendRawTransaction(signedTxn.Value)
	require.Nil(t, err)

	t.Log(hash)
}

func TestFetchTransactionDetail(t *testing.T) {
	// digest := "4nMHqXi60PLxj/DxLCWwkiO3L41kIz89qMDEpStRdP8="
	digest := "hPOfmwiRRsxleD0JGA67bWFBur+z1BdbLo6vYxzB+9w="
	chain := DevnetChain()

	detail, err := chain.FetchTransactionDetail(digest)
	require.Nil(t, err)

	t.Log(detail)
}

func TestSplit(t *testing.T) {
	account := M1Account(t)
	chain := DevnetChain()

	client, err := chain.client()
	require.Nil(t, err)

	signer, _ := types.NewAddressFromHex(account.Address())
	coin := "0x03149662d06e9427a67777092e03701b91af24a7"
	coinID, err := types.NewHexData(coin)
	require.Nil(t, err)

	txn, err := client.SplitCoinEqual(context.Background(), *signer, *coinID, 5, nil, 2000)
	signedTxn := txn.SignWith(account.account.PrivateKey)

	detail, err := client.ExecuteTransaction(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	require.Nil(t, err)
	t.Log(detail)
}
