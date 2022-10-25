package aptos

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOfferAndCancelTokenTransactionParams(t *testing.T) {
	privateKey := os.Getenv("PriMartian2")
	sender, err := AccountWithPrivateKey(privateKey)
	require.Nil(t, err)
	t.Log("sender address:", sender.Address())

	privateKey = os.Getenv("PriPetra1")
	receiver, err := AccountWithPrivateKey(privateKey)
	require.Nil(t, err)
	t.Log("receiver address:", receiver.Address())

	creator := "0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5"
	collection := "Martian Testnet98901"
	tokenName := "Martian NFT #98901"
	chain := NewChainWithRestUrl(testnetRestUrl)

	// Offer token
	offerTxn, err := chain.OfferTokenTransactionParams(sender, receiver.Address(), creator, collection, tokenName)
	require.Nil(t, err)
	// if true { /*
	if false { /**/
		txHash, err := chain.SendRawTransaction(offerTxn.Value)
		require.Nil(t, err)
		t.Log("offer token send success: hash =", txHash)

		time.Sleep(10 * time.Second)
	}

	// Cancel token
	cancelTxn, err := chain.CancelTokenOffer(sender, receiver.Address(), creator, collection, tokenName)
	require.Nil(t, err)
	// if true { /*
	if false { /**/
		txHash, err := chain.SendRawTransaction(cancelTxn.Value)
		require.Nil(t, err)
		t.Log("cancel offer send succeed: hash =", txHash)
	}
}

func TestClaimTokenFromHash(t *testing.T) {
	privateKey := os.Getenv("PriPetra1")
	receiver, err := AccountWithPrivateKey(privateKey)
	require.Nil(t, err)
	t.Log("receiver address:", receiver.Address())

	offerHash := "0x8a9673937e7d4f01a7f305dfd8ad18d29998540eb216e4aa304f4c68a3717f46"

	chain := NewChainWithRestUrl(testnetRestUrl)
	signedTxn, err := chain.ClaimTokenFromHash(receiver, offerHash)
	require.Nil(t, err)

	// if true { /*
	if false { /**/
		txHash, err := chain.SendRawTransaction(signedTxn.Value)
		require.Nil(t, err)
		t.Log("claim send succeed: hash =", txHash)
	}
}
