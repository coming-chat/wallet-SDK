package sui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	address := "0xbb8f7e72ae99d371020a1ccfe703bfb64a8a430f"

	chain := DevnetChain()
	b, err := chain.BalanceOfAddress(address)
	assert.Nil(t, err)

	t.Log(b)
}
