package sui

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBalance(t *testing.T) {
	address := "0x7e875ea78ee09f08d72e2676cf84e0f1c8ac61d94fa339cc8e37cace85bebc6e"

	chain := DevnetChain()
	b, err := chain.BalanceOfAddress(address)
	assert.Nil(t, err)

	t.Log(b)
}

func TestTokenBalance(t *testing.T) {
	chain := DevnetChain()
	token, err := NewToken(chain, "0x2d79a3c70aa3f3a3feabbf54b7b520f956c4ef8d::AAA::AAA")
	require.NoError(t, err)

	balance, err := token.BalanceOfAddress("0x2ecb102385afd954bf06f2a3a4ac648eb7a536e0")
	require.NoError(t, err)
	require.Equal(t, "0", balance.Total) // invalid address
}

func TestTokenInfo(t *testing.T) {
	chain := DevnetChain()
	token, err := NewToken(chain, "0x2d79a3c70aa3f3a3feabbf54b7b520f956c4ef8d::AAA::AAA")
	require.NoError(t, err)

	info, err := token.TokenInfo()
	require.Error(t, err) // token not found
	t.Log(info)

	mainToken := NewTokenMain(chain)
	tokenInfo, err := mainToken.TokenInfo()
	require.NoError(t, err)

	t.Log(tokenInfo)
}
