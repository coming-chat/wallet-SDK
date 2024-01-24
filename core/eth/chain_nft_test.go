package eth

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestChain_TransferNFT_Erc721(t *testing.T) {
	mn := testcase.M1
	sender, err := NewAccountWithMnemonic(mn)
	require.Nil(t, err)

	receiver := sender.Address()
	nftId := "0"
	nftContract := "0x199Dcb0132a66b05723882259832e240fF735810"
	nftStandard := "erc-721"

	chain := NewChainWithRpc("https://canary-testnet.bevm.io/")

	txn, err := chain.TransferNFTParams(sender.Address(), receiver,
		nftId, nftContract, nftStandard)
	require.Nil(t, err)

	signedTx, err := chain.BuildTransferTxWithAccount(sender, txn)
	require.Nil(t, err)

	run := false
	if run {
		txHash, err := chain.SendRawTransaction(signedTx.Value)
		require.Nil(t, err)
		t.Log(txHash)
	}
}
