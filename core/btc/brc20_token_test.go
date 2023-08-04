package btc

import (
	"testing"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/stretchr/testify/require"
)

func TestBrc20Token_TokenInfo(t *testing.T) {
	token := NewBrc20Token("bruh")
	info, err := token.TokenInfo()
	require.Nil(t, err)
	require.Equal(t, info, &base.TokenInfo{
		Name:    "BRUH",
		Symbol:  "BRUH",
		Decimal: 18,
	})
}

func TestBrc20Token_FullTokenInfo(t *testing.T) {
	token := NewBrc20Token("meme")
	info, err := token.FullTokenInfo()
	require.Nil(t, err)
	require.Equal(t, info.Decimal, int16(18))
	require.Equal(t, info.Max, "99999")

	timeStart := time.Now().UnixMilli()
	info222, err := token.FullTokenInfo()
	timeSpent := time.Now().UnixMilli() - timeStart
	require.Nil(t, err)
	require.True(t, timeSpent < 10) // The second use of the cache should be very fast
	require.Equal(t, info, info222)
}

func TestBrc20TokenBalances(t *testing.T) {
	owner := "bc1qdgflzu306s75lgskkgssmz3vscpvuawvafv3xjshyc6t73x3zzvquvtafp"

	chain, err := NewChainWithChainnet(ChainMainnet)
	// chain, err := NewChainWithChainnet(ChainTestnet)
	require.Nil(t, err)
	balancePage, err := chain.FetchBrc20TokenBalance(owner, "0", 10)
	require.Nil(t, err)
	t.Log(balancePage.ItemArray().Values...)
}
