package sui

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

const DevnetRpcUrl = "https://gateway.devnet.sui.io:443"

func DevnetChain() *Chain {
	return NewChainWithRpcUrl(DevnetRpcUrl)
}

func TestTransfer(t *testing.T) {
	account := ChromeAccount()
	chain := DevnetChain()
	token := NewTokenMain(chain)

	toAddress := "0xc4173a804406a365e69dfb297d4eaaf002546ebd"
	amount := "10"
	// toAddress := M1Account().Address()
	// amount := "120000"

	signedTxn, err := token.BuildTransferTxWithAccount(account, toAddress, amount)
	require.Nil(t, err)

	hash, err := chain.SendRawTransaction(signedTxn.Value)
	require.Nil(t, err)

	t.Log(hash)
}

func TestFetchTransactionDetail(t *testing.T) {
	// digest := "4nMHqXi60PLxj/DxLCWwkiO3L41kIz89qMDEpStRdP8="
	digest := "RiP1hhhaNQKwJaEl+KixLtrkW1Z8WT8jtrzv8LLasA0="
	chain := DevnetChain()

	detail, err := chain.FetchTransactionDetail(digest)
	require.Nil(t, err)

	t.Log(detail)
}

func TestSplit(t *testing.T) {
	account := ChromeAccount()
	chain := DevnetChain()

	client, err := chain.client()
	require.Nil(t, err)

	signer, _ := types.NewAddressFromHex(account.Address())
	objId := "0xed763e483fe8a87a0e2568f9bf8c7b02c7034f5d"
	ID, err := types.NewHexData(objId)
	require.Nil(t, err)

	txn, err := client.SplitCoinEqual(context.Background(), *signer, *ID, 5, nil, 2000)
	signedTxn := txn.SignWith(account.account.PrivateKey)

	detail, err := client.ExecuteTransaction(context.Background(), *signedTxn)
	require.Nil(t, err)
	t.Log(detail)
}
