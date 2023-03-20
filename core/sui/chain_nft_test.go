package sui

import (
	"testing"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/stretchr/testify/require"
)

func TestFetchNfts(t *testing.T) {
	// owner := "0xd059ab4c5c7d2be6537101f76c41f25cdf5cc26e"
	owner := M1Account(t).Address()
	chain := DevnetChain()

	nfts, err := chain.FetchNFTs(owner)
	require.Nil(t, err)
	for name, group := range nfts {
		t.Log("=======================================")
		t.Logf("group: %v, count: %v", name, len(group))
		for idx, nft := range group {
			t.Logf("%4v: %v", idx, nft)
		}
	}
}

func TestMintNFT(t *testing.T) {
	account := M1Account(t)
	chain := DevnetChain()

	var (
		timeNow = time.Now().Format("06-01-02 15:04")
		nftName = "ComingChat NFT at " + timeNow
		nftDesc = "This is a NFT created by ComingChat"
		nftUrl  = "https://coming.chat/favicon.ico"
	)
	txn, err := chain.MintNFT(account.Address(), nftName, nftDesc, nftUrl, "", MaxGasBudget)
	require.Nil(t, err)
	signedTxn, err := txn.SignWithAccount(account)
	require.Nil(t, err)
	hash, err := chain.SendRawTransaction(signedTxn.Value)
	require.Nil(t, err)
	t.Log("mint nft success, hash = ", hash)
}

func TestTransferNFT(t *testing.T) {
	account := M1Account(t)
	receiver := M2Account(t).Address()

	chain := TestnetChain()

	nfts, err := chain.FetchNFTs(account.Address())
	require.Nil(t, err)
	var nft *base.NFT
out:
	for _, group := range nfts {
		for _, n := range group {
			nft = n
			break out
		}
	}
	require.NotNil(t, nft)

	txn, err := chain.TransferNFT(account.Address(), receiver, nft.Id, "", MaxGasBudget)
	require.Nil(t, err)
	signedTxn, err := txn.SignWithAccount(account)
	require.Nil(t, err)
	hash, err := chain.SendRawTransaction(signedTxn.Value)
	require.Nil(t, err)
	t.Log("transfer nft success, hash = ", hash)
}
