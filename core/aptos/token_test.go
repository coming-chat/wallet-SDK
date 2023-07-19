package aptos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	btcTag  = "0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::BTC"
	usdtTag = "0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::USDT"
)

func TestTokenInfo(t *testing.T) {
	chain := NewChainWithRestUrl(devnetRestUrl)

	token := NewMainToken(chain)
	info, err := token.TokenInfo()
	t.Log(info, err)

	btcToken, err := NewToken(chain, btcTag)
	info, err = btcToken.TokenInfo()
	t.Log(info, err)

	usdtToken, err := NewToken(chain, usdtTag)
	info, err = usdtToken.TokenInfo()
	t.Log(info, err)
}

func TestTokenBalance(t *testing.T) {
	address := "0xe1c1deec04ed6d7f92f867875c5c9733b64e376ca5a7f5da5b6bdaf3dd28eb9c"
	chain := NewChainWithRestUrl(devnetRestUrl)

	aptToken := NewMainToken(chain)
	aptBalance, err := aptToken.BalanceOfAddress(address)
	assert.Nil(t, err)
	t.Logf("APT Balance = %v", aptBalance.Total)

	btcToken, err := NewToken(chain, btcTag)
	assert.Nil(t, err)
	btcBalance, err := btcToken.BalanceOfAddress(address)
	assert.Nil(t, err)
	t.Logf("BTC Balance = %v", btcBalance.Total)

	usdtToken, err := NewToken(chain, usdtTag)
	assert.Nil(t, err)
	usdtBalance, err := usdtToken.BalanceOfAddress(address)
	assert.Nil(t, err)
	t.Logf("USDT Balance = %v", usdtBalance.Total)
}

func TestBuildTransferTxWithAccount(t *testing.T) {
	account := M1Account(t)
	toAddress := "0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5"
	amount := "100"

	chain := NewChainWithRestUrl(testnetRestUrl)
	token := NewMainToken(chain)

	signedTxn, err := token.BuildTransferTxWithAccount(account, toAddress, amount)
	require.Nil(t, err)

	if false {
		txHash, err := chain.SendRawTransaction(signedTxn.Value)
		require.Nil(t, err)
		t.Log(txHash)
	}
}

func TestBuildTransafer(t *testing.T) {
	account := M1Account(t)
	toAddress := "0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5"
	amount := "100"

	chain := NewChainWithRestUrl(testnetRestUrl)
	token := NewMainToken(chain)

	txn, err := token.BuildTransfer(account.Address(), toAddress, amount)
	require.Nil(t, err)

	estimateFee, err := chain.EstimateTransactionFeeUsePublicKey(txn, account.PublicKeyHex())
	require.Nil(t, err)
	t.Log(estimateFee)

	signedTxn, err := txn.SignWithAccount(account)
	require.Nil(t, err)

	if false {
		txHash, err := chain.SendRawTransaction(signedTxn.Value)
		require.Nil(t, err)
		t.Log(txHash)
	}
}

func TestEstimateFee(t *testing.T) {
	account := M1Account(t)
	toAddress := "0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5"
	amount := "100"

	chain := NewChainWithRestUrl(devnetRestUrl)
	token := NewMainToken(chain)

	fee, err := token.EstimateFees(account, toAddress, amount)
	assert.Nil(t, err)
	t.Log(fee)
}

func TestTokenRegister(t *testing.T) {
	chain := NewChainWithRestUrl(devnetRestUrl)
	account := M1Account(t)

	token, err := NewToken(chain, btcTag)
	assert.Nil(t, err)

	_, err = token.EnsureOwnerRegistedToken(account)
	assert.Nil(t, err)

	// test duplicate registration
	hash, err := token.RegisterTokenForOwner(account)
	t.Log(hash, err)
}

func TestToken_BuildTransfer_SignedTransaction(t *testing.T) {
	account := M1Account(t)
	chain := ChainTestnet()
	token := chain.MainToken()

	balance, err := token.BalanceOfAddress(account.Address())
	require.Nil(t, err)
	t.Log("sender address = ", account.Address())
	t.Log("balance = ", balance.Usable)

	txn, err := token.BuildTransfer(account.Address(), account.Address(), "100")
	require.Nil(t, err)

	gasfee, err := chain.EstimateTransactionFeeUsePublicKey(txn, account.PublicKeyHex())
	require.Nil(t, err)
	t.Log("Estimate fee = ", gasfee.Value)

	signedTxn, err := txn.SignedTransactionWithAccount(account)
	require.Nil(t, err)

	if false {
		hash, err := chain.SendSignedTransaction(signedTxn)
		require.Nil(t, err)
		t.Log("Transaction hash = ", hash.Value)
	}
}
