package wallet

import (
	"testing"
)

func TestChainType(t *testing.T) {
	address := "😁"

	chains := ChainTypeFrom(address)
	t.Log(chains.String())
	t.Log(chains.Count())
}
