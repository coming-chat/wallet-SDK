package sui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	address := "0x2ecb102385afd954bf06f2a3a4ac648eb7a536e0"

	chain := DevnetChain()
	b, err := chain.BalanceOfAddress(address)
	assert.Nil(t, err)

	t.Log(b)
}

func TestTokenBalance(t *testing.T) {
	chain := DevnetChain()
	token, err := NewToken(chain, "0x2d79a3c70aa3f3a3feabbf54b7b520f956c4ef8d::AAA::AAA")
	if err != nil {
		t.Fatal(err)
	}
	balance, err := token.BalanceOfAddress("0x2ecb102385afd954bf06f2a3a4ac648eb7a536e0")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)
}
