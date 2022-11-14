package wallet

import "testing"

func TestChainType(t *testing.T) {
	address := "12eV7FtPbXBgDG6mX4zwPaJdQKgigVtnYofSpS8mgEQbX625"

	chains := ChainTypeFrom(address)
	t.Log(chains.String())
	t.Log(chains.Count())
}
