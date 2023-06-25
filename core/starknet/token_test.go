package starknet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenInfo(t *testing.T) {
	chain := GoerliChain()
	token := chain.MainToken()

	info, err := token.TokenInfo()
	require.Nil(t, err)
	t.Log(info)
}
