package aptos

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/assert"
)

const (
	btcTag  = "0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::BTC"
	usdtTag = "0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::coins::USDT"
)

func TestTokenInfo(t *testing.T) {
	chain := NewChainWithRestUrl(testnetRestUrl)

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
	chain := NewChainWithRestUrl(testnetRestUrl)

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

func TestTransafer(t *testing.T) {
	account, err := NewAccountWithMnemonic(testcase.M1)
	assert.Nil(t, err)
	toAddress := "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"
	amount := "100"

	chain := NewChainWithRestUrl(testnetRestUrl)
	token := NewMainToken(chain)

	signedTx, err := token.BuildTransferTxWithAccount(account, toAddress, amount)
	assert.Nil(t, err)
	t.Log(signedTx.Value)

	txHash, err := chain.SendRawTransaction(signedTx.Value)
	assert.Nil(t, err)
	t.Log(txHash)
}

func TestEstimateFee(t *testing.T) {
	account, err := NewAccountWithMnemonic(testcase.M1)
	assert.Nil(t, err)
	toAddress := "0x559c26e61a74a1c40244212e768ab282a2cbe2ed679ad8421f7d5ebfb2b79fb5"
	amount := "100"

	chain := NewChainWithRestUrl(testnetRestUrl)
	token := NewMainToken(chain)

	fee, err := token.EstimateFees(account, toAddress, amount)
	assert.Nil(t, err)
	t.Log(fee)
}

func TestTokenRegister(t *testing.T) {
	chain := NewChainWithRestUrl(testnetRestUrl)
	account, err := NewAccountWithMnemonic(testcase.M1)
	assert.Nil(t, err)

	token, err := NewToken(chain, btcTag)
	assert.Nil(t, err)

	_, err = token.EnsureOwnerRegistedToken(account)
	assert.Nil(t, err)

	// test duplicate registration
	hash, err := token.RegisterTokenForOwner(account)
	t.Log(hash, err)
}
