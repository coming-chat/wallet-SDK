package btc

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransaction_AddOpReturn(t *testing.T) {
	txn, err := NewTransaction(ChainSignet)
	require.NoError(t, err)
	err = txn.AddOpReturn("ComingChat")
	require.NoError(t, err)
	require.Equal(t, hex.EncodeToString(txn.msgTx.TxOut[0].PkScript), "6a0a436f6d696e6743686174")

	err = txn.AddOpReturn("len75hellohellohellohellohellohellohellohellohellohellohellohellohellohello")
	require.NoError(t, err)
	err = txn.AddOpReturn("len76hellohellohellohellohellohellohellohellohellohellohellohellohellohello#")
	require.Error(t, err) // error len > 75
}

func TestTransaction_Sign(t *testing.T) {
	var (
		mnemonic       = "antenna chaos arrive hungry distance human question history decade deal impose color"
		accTaproot, _  = NewAccountWithMnemonic(mnemonic, ChainSignet, AddressTypeTaproot)
		accComingTp, _ = NewAccountWithMnemonic(mnemonic, ChainSignet, AddressTypeComingTaproot)
		accNativeSg, _ = NewAccountWithMnemonic(mnemonic, ChainSignet, AddressTypeNativeSegwit)
		accNestedSg, _ = NewAccountWithMnemonic(mnemonic, ChainSignet, AddressTypeNestedSegwit)
		accLegacy, _   = NewAccountWithMnemonic(mnemonic, ChainSignet, AddressTypeLegacy)
	)
	t.Log("Taproot: ", accTaproot.Address())
	t.Log("ComingTaproot: ", accComingTp.Address())
	t.Log("NativeSegwit: ", accNativeSg.Address())
	t.Log("NestedSegwit: ", accNestedSg.Address())
	t.Log("Legacy: ", accLegacy.Address())

	// SendSignedTransaction
	// signedTxn, err := txn.SignedTransactionWithAccount(accTaproot)
	// hash, err := chain.SendSignedTransaction(signedTxn)
	//
	// SendRawTransaction
	// signedTxn, err := txn.SignWithAccount(accTaproot)
	// hash, err := chain.SendRawTransaction(signedTxn.Value)

	{
		from, to := accTaproot, accComingTp

		txn, err := NewTransaction(ChainSignet) // two input
		require.NoError(t, err)
		err = txn.AddInput("7ad33b3012aef3e153c65b57449ae4b279017aa99a6b581ea070903b3ca2b73f", 1, from.Address(), 6602)
		require.NoError(t, err)
		err = txn.AddInput("e69b40f50c74cfcd414e5c0c8a77146e2c679d9c2b521186d58d00ca17756713", 9, from.Address(), 1000000)
		require.NoError(t, err)
		err = txn.AddOutput(to.Address(), 1006300)
		require.NoError(t, err)

		txHex, err := txn.SignWithAccount(from)
		require.NoError(t, err)
		require.Equal(t, txHex.Value, "0x010000000001023fb7a23c3b9070a01e586b9aa97a0179b2e49a44575bc653e1f3ae12303bd37a0100000000ffffffff13677517ca008dd58611522b9c9d672c6e14778a0c5c4e41cdcf740cf5409be60900000000ffffffff01dc5a0f0000000000225120f78430ddf1178c9c04bd32e3e51d0aa72760756934f95173cb7926083e115062014095001ab2628b184b396fa6daf0d17cee878dddca7822c7b2d69c540ed8572ef82c11fc4daa6d51c8afdb436a1413f081015107d798ec2f4fe62ea6fcc5f3bd9401400b602a8a4d9c979e18b15f202fc7d4eb4c1ff410005a38cbc48dd6a7ef363eeb0d1833fdbbf949d01aa83f6b6873b81f021a9b02126d0b375b512df604bb105700000000")
		// txn detail: https://mempool.space/signet/tx/f926ce5fd50735c30c1f4de0c0d61d0ee14d24a4e3fc2210f410abbf1335cacf
	}

	{
		from, to := accComingTp, accNativeSg

		txn, err := NewTransaction(ChainSignet)
		require.NoError(t, err)
		err = txn.AddInput("6f846e1eb21a7bf8a232e4cb07eb6c3af7ca99e107be4e975d089db8f5d9b67e", 0, from.Address(), 1005600)
		require.NoError(t, err)
		err = txn.AddOutput(to.Address(), 1005400)
		require.NoError(t, err)
		err = txn.AddOpReturn("ComingChat Wallet SDK Test") // add op_return
		require.NoError(t, err)

		txHex, err := txn.SignWithAccount(from)
		require.NoError(t, err)
		require.Equal(t, txHex.Value, "0x010000000001017eb6d9f5b89d085d974ebe07e199caf73a6ceb07cbe432a2f87b1ab21e6e846f0000000000ffffffff0258570f0000000000160014b2f9ed398433cd438d5e88d62b2e28b440901cc500000000000000001c6a1a436f6d696e67436861742057616c6c65742053444b20546573740140767a1c1b1ea639b92c59f9b49a38e794526c2f301cf9000ad0e58ebe2ea832d448d3711a494f933a95ed72bfa1642418eaa1dcc092c0530b6351a80c5ea995b600000000")
		// txn detail: https://mempool.space/signet/tx/6ba14977701bf831615d2ecca5d1e5b020eb6fe8b999e2f5f1a9b84f636c1705
	}

	nativeSg_TxnHex := "0x0100000000010105176c634fb8a9f1f5e299b9e86feb20b0e5d1a5cc2e5d6131f81b707749a16b0000000000ffffffff0290560f000000000017a914fc910555e59c5449120f5ba0b5fdf75c5f7e3e718700000000000000001d6a1b436f6d696e67436861742057616c6c65742053444b20546573743202473044022037840ab4706872c141a60a065c1a0f65fd725bad1529d5149d1f241d5b44981d022047445550484ffe0b5db4316a3fd264437e2a9a7f8a6e26fdb2bc703044632d46012103c17dd7d56e18c38293a882dcd559f3e4e8ee8b5643f8d8e5841bb49447fb3a0f00000000"
	{
		from, to := accNativeSg, accNestedSg

		txn, err := NewTransaction(ChainSignet)
		require.NoError(t, err)
		err = txn.AddInput("6ba14977701bf831615d2ecca5d1e5b020eb6fe8b999e2f5f1a9b84f636c1705", 0, from.Address(), 1005400)
		require.NoError(t, err)
		err = txn.AddOutput(to.Address(), 1005200)
		require.NoError(t, err)
		err = txn.AddOutput("ComingChat Wallet SDK Test2", 0) // value 0 is op_return
		require.NoError(t, err)

		txHex, err := txn.SignWithAccount(from)
		require.NoError(t, err)
		require.Equal(t, txHex.Value, nativeSg_TxnHex)
		// txn detail: https://mempool.space/signet/tx/c0a4c2b41deadd4122a15c809079cf7d19c5b49dccb0247599427ff607295381
	}

	{
		from, to := accNestedSg, accLegacy

		txn, err := NewTransaction(ChainSignet) // two output
		require.NoError(t, err)
		// test input with prevTx
		err = txn.AddInput2("c0a4c2b41deadd4122a15c809079cf7d19c5b49dccb0247599427ff607295381", 0, nativeSg_TxnHex)
		require.NoError(t, err)
		err = txn.AddOutput(to.Address(), 1000000)
		require.NoError(t, err)
		err = txn.AddOutput(to.Address(), 4900)
		require.NoError(t, err)

		txHex, err := txn.SignWithAccount(from)
		require.NoError(t, err)
		require.Equal(t, txHex.Value, "0x0100000000010181532907f67f42997524b0cc9db4c5197dcf7990805ca12241ddea1db4c2a4c0000000001716001441d432f2b3d3daa1df08e21fbdbbcdfb0908001bffffffff0240420f00000000001976a9147a494419776de8f4cb846882a15df3cc836299cf88ac24130000000000001976a9147a494419776de8f4cb846882a15df3cc836299cf88ac024730440220731dfff44876d17a9b9cbf08a33a238e61f79f09fb73bd3a7541a93a9df4b450022040ec119a288c042077e2dd219e94e23def37bb93fc4f932fd1a6805e7b6681440121026612352433fe00fed2c41213b79be357daeac671f51812918d057e35a7a36cc000000000")
		// txn detail: https://mempool.space/signet/tx/87db2ca3d888c74d96fef6a516b94b33be24e09e4a5e6418a96ecb5ab6cf959c
	}

	{
		from, to := accLegacy, "tb1pdq423fm5dv00sl2uckmcve8y3w7guev8ka6qfweljlu23mmsw63qk6w2v3"

		txn, err := NewTransaction(ChainSignet)
		require.NoError(t, err)
		err = txn.AddInput("87db2ca3d888c74d96fef6a516b94b33be24e09e4a5e6418a96ecb5ab6cf959c", 0, from.Address(), 1000000)
		require.NoError(t, err)
		err = txn.AddInput("87db2ca3d888c74d96fef6a516b94b33be24e09e4a5e6418a96ecb5ab6cf959c", 1, from.Address(), 4900)
		require.NoError(t, err)
		err = txn.AddOutput(to, 1004500)
		require.NoError(t, err)

		txHex, err := txn.SignWithAccount(from)
		require.NoError(t, err)
		require.Equal(t, txHex.Value, "0x01000000029c95cfb65acb6ea918645e4a9ee024be334bb916a5f6fe964dc788d8a32cdb87000000006a47304402201d252ee0e20f82e0c13ee3294ad6cfd0ab6021c276724c41fec1a80f56ee1a1802202f65a9c4108c10fe5f481b63f1e464dfd14f5041fe617c37777c6ecf975e147d0121029b740120a2a4af5d751d5e3d67f6d2aa9f92792af6bc0df37d30584b6d65bb54ffffffff9c95cfb65acb6ea918645e4a9ee024be334bb916a5f6fe964dc788d8a32cdb87010000006a47304402200d6aabb1ff6eef5aa2b01fcd0e00bdbd0c1dda63a6842c0ad4b06b1e91ba46450220428bee9060569844b9d19206576ee648924817c1f2d42a379f22f7f2115122de0121029b740120a2a4af5d751d5e3d67f6d2aa9f92792af6bc0df37d30584b6d65bb54ffffffff01d4530f0000000000225120682aa8a7746b1ef87d5cc5b78664e48bbc8e6587b77404bb3f97f8a8ef7076a200000000")
		// txn detail: https://mempool.space/signet/tx/932bce660d7336aea21213a4d7a25c7fa82ae66a0e5cc9e1e8a82a88b7c40622
	}

}
