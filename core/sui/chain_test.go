package sui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const DevnetRpcUrl = "https://gateway.devnet.sui.io:443"

func DevnetChain() *Chain {
	return NewChainWithRpcUrl(DevnetRpcUrl)
}

func TestTransfer(t *testing.T) {
	account := ChromeAccount()
	chain := DevnetChain()
	token := NewTokenMain(chain)

	// toAddress := "0xc4173a804406a365e69dfb297d4eaaf002546ebd"
	// amount := "1"
	toAddress := M1Account().Address()
	amount := "10000"

	signedTxn, err := token.BuildTransferTxWithAccount(account, toAddress, amount)
	assert.Nil(t, err)

	hash, err := chain.SendRawTransaction(signedTxn.Value)
	assert.Nil(t, err)

	t.Log(hash)
}

func TestFetchTransactionDetail(t *testing.T) {
	// digest := "4nMHqXi60PLxj/DxLCWwkiO3L41kIz89qMDEpStRdP8="
	digest := "RiP1hhhaNQKwJaEl+KixLtrkW1Z8WT8jtrzv8LLasA0="
	chain := DevnetChain()

	detail, err := chain.FetchTransactionDetail(digest)
	assert.Nil(t, err)

	t.Log(detail)
}
