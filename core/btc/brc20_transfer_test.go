package btc

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestBuildTransferTransaction(t *testing.T) {
	sender := "tb1p2hsjm57fsxrqcq5p42get87ttrw069kqa2ar444ma4ussquuaklqfsrknz"

	chain, err := NewChainWithChainnet(ChainTestnet)
	require.Nil(t, err)

	idList := inscriptionIds(t, chain, sender, "txtx", 2)
	idArray := base.StringArray{Values: idList}

	txn, err := chain.BuildBrc20TransferTransaction(sender, sender, &idArray, 1, "")
	require.Nil(t, err)

	psbtTxn, err := txn.ToPsbtTransaction()
	require.Nil(t, err)

	t.Log(psbtTxn)
}

func TestBrc20TransferTransaction_SignAndSend(t *testing.T) {
	account, err := NewAccountWithMnemonic(testcase.M1, ChainTestnet)
	require.Nil(t, err)
	chain, err := NewChainWithChainnet(ChainTestnet)
	require.Nil(t, err)
	user, err := account.TaprootAddress()
	require.Nil(t, err)

	idList := inscriptionIds(t, chain, user.Value, "txtx", 2)
	idArray := base.StringArray{Values: idList}

	txn, err := chain.BuildBrc20TransferTransaction(user.Value, user.Value, &idArray, 2, "")
	require.Nil(t, err)
	psbtTxn, err := txn.ToPsbtTransaction()
	require.Nil(t, err)
	signedTxn, err := psbtTxn.SignedTransactionWithAccount(account)
	require.Nil(t, err)

	if false {
		hash, err := chain.SendSignedTransaction(signedTxn)
		require.Nil(t, err)
		t.Log("transfer success:", hash.Value)
	}
}

func inscriptionIds(t *testing.T, chain *Chain, owner string, tick string, maxCount int) []string {
	arr, err := chain.FetchBrc20TransferableInscription(owner, tick)
	require.Nil(t, err)
	require.Greater(t, arr.ItemArray().Count(), 0, "the user has no inscription")

	idList := []string{}
	for _, item := range arr.Items {
		idList = append(idList, item.InscriptionId)
	}
	if len(idList) > maxCount {
		return idList[:maxCount]
	}
	return idList
}