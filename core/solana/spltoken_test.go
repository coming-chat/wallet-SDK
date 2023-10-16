package solana

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const aTokenAddress = "38GJtbJkJJQKkR6CTZwwok4hAXoKnByU2MFAZjgTe114"

func TestSPLToken_GetBalance(t *testing.T) {
	chain := DevnetChain()
	token, err := NewSPLToken(chain, "38GJtbJkJJQKkR6CTZwwok4hAXoKnByU2MFAZjgTe114")
	require.Nil(t, err)

	balance, err := token.BalanceOfAddress("4MPScMzmKwQfpzQ4MtkSaqKQbTEzGsWqovUMweNz7nFo")
	require.Nil(t, err)
	t.Log(balance)
}

func TestSPLToken_CreateTokenAccount(t *testing.T) {
	signer := M1Account(t)
	owner := "9HquD8jfJNm1wKuVEbeYWCr792Yt9D9zj17i1rgzC7t4"

	chain := DevnetChain()
	token, err := NewSPLToken(chain, "5PtchJqBwDiAJnvoHKYRKhXYmgurAJygj4WSL9C9xghJ")
	require.Nil(t, err)

	txn, err := token.CreateTokenAccount(owner, signer.Address())
	require.Nil(t, err)

	fee, err := chain.EstimateTransactionFee(txn)
	require.Nil(t, err)
	t.Log(fee.Value)

	signedTxn, err := txn.SignedTransactionWithAccount(signer)
	require.Nil(t, err)

	if false {
		txhash, err := chain.SendSignedTransaction(signedTxn)
		require.Nil(t, err)
		t.Log("create token account for other user success, hash = ", txhash)
	}
}

func TestSPLToken_BuildTransfer(t *testing.T) {
	account := M1Account(t)
	chain := DevnetChain()

	mintAddr := "38GJtbJkJJQKkR6CTZwwok4hAXoKnByU2MFAZjgTe114"
	token, err := NewSPLToken(chain, mintAddr)
	require.Nil(t, err)

	receiverAddr := "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqPY"
	txn, err := token.BuildTransfer(account.Address(), receiverAddr, "100")
	// txn, err := token.BuildTransferAll(account.Address(), receiverAddr)
	require.Nil(t, err)

	fee, err := chain.EstimateTransactionFee(txn)
	require.Nil(t, err)
	t.Log("estimate transaction fee = ", fee.Value)

	signedTxn, err := txn.SignedTransactionWithAccount(account)
	require.Nil(t, err)
	if false {
		txhash, err := chain.SendSignedTransaction(signedTxn)
		require.Nil(t, err)
		t.Log("transfer success hash = ", txhash)
	}
}

func TestSPLToken_TokenInfo(t *testing.T) {
	chain := MainnetChain()
	token, err := NewSPLToken(chain, "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB") // USDT
	require.Nil(t, err)
	info, err := token.TokenInfo()
	require.Nil(t, err)
	t.Log(info)

	chain = DevnetChain()
	token, err = NewSPLToken(chain, "38GJtbJkJJQKkR6CTZwwok4hAXoKnByU2MFAZjgTe114")
	require.Nil(t, err)
	info, err = token.TokenInfo()
	require.Nil(t, err)
	t.Log(info)
}
