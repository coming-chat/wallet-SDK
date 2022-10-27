package aptos

import (
	"os"
	"testing"
	"time"

	txbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/lcs"
	"github.com/stretchr/testify/require"
)

func TestCIDPayload(t *testing.T) {
	cidBuilder := NewNFTPayloadBCSBuilder("")
	var payload txbuilder.TransactionPayloadEntryFunction

	if true { /*
			if false { /**/
		bytes, err := cidBuilder.CIDAllowDirectTransferPayload()
		require.Nil(t, err)
		err = lcs.Unmarshal(bytes, &payload)
		require.Nil(t, err)
	}

	if true { /*
			if false { /**/
		privateKey := os.Getenv("PriPetra1")
		receiver, err := AccountWithPrivateKey(privateKey)
		require.Nil(t, err)
		t.Log("receiver address:", receiver.Address())

		bytes, err := cidBuilder.CIDTokenTransferPayload(1234, receiver.Address())
		require.Nil(t, err)
		err = lcs.Unmarshal(bytes, &payload)
		require.Nil(t, err)
	}

	if true { /*
			if false { /**/
		bytes, err := cidBuilder.CIDRegister(2134)
		require.Nil(t, err)
		err = lcs.Unmarshal(bytes, &payload)
		require.Nil(t, err)
	}
}

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

	nftBuilder := NewNFTPayloadBCSBuilder("")

	// Offer token
	offerPayload, err := nftBuilder.OfferTokenTransactionParams(receiver.Address(), creator, collection, tokenName)
	require.Nil(t, err)
	// if true { /*
	if false { /**/
		txHash, err := chain.SubmitTransactionPayloadBCS(sender, offerPayload)
		require.Nil(t, err)
		t.Log("offer token send success: hash =", txHash)

		time.Sleep(10 * time.Second)
	}

	// Cancel token
	cancelPayload, err := nftBuilder.CancelTokenOffer(receiver.Address(), creator, collection, tokenName)
	require.Nil(t, err)
	// if true { /*
	if false { /**/
		txHash, err := chain.SubmitTransactionPayloadBCS(sender, cancelPayload)
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
	nftBuilder := NewNFTPayloadBCSBuilder("")

	claimPayload, err := nftBuilder.ClaimTokenFromHash(offerHash, chain, receiver.Address())
	require.Nil(t, err)

	// if true { /*
	if false { /**/
		txHash, err := chain.SubmitTransactionPayloadBCS(receiver, claimPayload)
		require.Nil(t, err)
		t.Log("claim send succeed: hash =", txHash)
	}
}

func TestChain_IsAccountAllowedDirectTransferToken(t *testing.T) {
	chain := NewChainWithRestUrl(testnetRestUrl)

	type args struct {
		address string
	}
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{
			name:    "error address",
			address: "0xopqrst",
			wantErr: true,
		},
		{
			name:    "may be is true",
			address: "0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5",
		},
		{
			name:    "maybe is false",
			address: "0x1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chain.IsAllowedDirectTransferToken(tt.address)
			require.Equal(t, err != nil, tt.wantErr, "Chain.IsAccountAllowedDirectTransferToken() error = %v, wantErr %v", err, tt.wantErr)
			if err == nil {
				t.Logf("allow = %v of address %v", got.Value, tt.address)
			}
		})
	}
}
