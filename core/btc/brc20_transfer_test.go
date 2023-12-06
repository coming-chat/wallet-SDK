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
	idArray := base.StringArray{AnyArray: idList}

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
	idArray := base.StringArray{AnyArray: idList}

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
	require.Greater(t, arr.CurrentCount(), 0, "the user has no inscription")

	idList := []string{}
	for _, item := range arr.Items {
		idList = append(idList, item.InscriptionId)
	}
	if len(idList) > maxCount {
		return idList[:maxCount]
	}
	return idList
}

func Test_Marshal_Brc20TransferTransaction(t *testing.T) {
	jsonStr := `
	{
		"transaction": "70736274ff0100e701",
		"commit_custom": [
			"02220200000000000",
			"8d3fd6c21cf9242110d1646d5ae313a1233c2f4e8b597c0115d0f2a798334f13",
			"0",
			"ad683b00dc5e63ed0034263ee7cb92db972f550ce861dc6785316e1e095041cb",
			"1"
		],
		"network_fee": 795
	}
	`

	txn, err := NewBrc20TransferTransactionWithJsonString(jsonStr)
	require.Nil(t, err)
	cmt := txn.CommitCustom
	require.Equal(t, cmt.BaseTx, "02220200000000000")
	require.Equal(t, cmt.Utxos.Count(), 2)
	require.Equal(t, cmt.Utxos.ValueAt(0).Txid, "8d3fd6c21cf9242110d1646d5ae313a1233c2f4e8b597c0115d0f2a798334f13")
	require.Equal(t, cmt.Utxos.ValueAt(1).Index, int64(1))
}
