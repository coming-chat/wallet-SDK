package wallet

import "testing"

func TestQueryBalance(t *testing.T) {
	rpcUrl := "wss://testnet3.chainx.org"
	// address := "5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnH" // no balance
	address := "5UXKnBuqVdoBgRDxrCgZErJojebb1pYivt4ei9D8NYQkbg9U" // have balance
	// address := "5UXKnBuqVdoBgRDxrCg" // error address

	chain := NewPolkaChain(rpcUrl)

	balance, err := chain.QueryBalance(address)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(balance)
}

func TestQueryBalancePubkey(t *testing.T) {
	rpcUrl := "wss://testnet3.chainx.org"
	pubkey := "0xd8110ae501b7f7b12d7bb3a0596097828b2b9398e4ff6dfa4bc8bdb1dc0e505d"
	// pubkey := ""
	// pubkey := "0xd8110ae501b7f7b12d7bb3a0596"

	chain := NewPolkaChain(rpcUrl)

	balance, err := chain.QueryBalancePubkey(pubkey)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(balance)
}

func TestXBTCBalance(t *testing.T) {
	rpcUrl := "wss://minichain-testnet.coming.chat" // 非 chainx 会抛出 error
	rpcUrl = "wss://testnet3.chainx.org"            // chainx 才可以正常查询
	address := "5QUEnWNMDFqsbUGpvvtgWGUgiiojnEpLf7581ELLAQyQ1xnT"

	chain := NewPolkaChain(rpcUrl)

	balance, err := chain.QueryBalanceXBTC(address)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(balance)
}

func TestMetadataString(t *testing.T) {
	rpcUrl := "wss://testnet3.chainx.org" // chainx 才可以正常查询

	client, err := newPolkaClient(rpcUrl, "")
	if err != nil {
		t.Fatal(err)
	}

	meta, err := client.MetadataString()
	if err != nil {
		t.Fatal(err)
	}

	chain, err := NewPolkaChainWithRpc(rpcUrl, meta)
	if err != nil {
		t.Fatal(err)
	}

	meta2, err := chain.GetMetadataString()
	if meta != meta2 {
		t.Fatal("metadata restore failed")
	}

	t.Log(chain)
}
