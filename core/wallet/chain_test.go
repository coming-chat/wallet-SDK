package wallet

import (
	"testing"
)

func TestQueryBalance(t *testing.T) {
	rpcUrl := "wss://testnet3.chainx.org"
	// address := "5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnH" // no balance
	address := "5UXKnBuqVdoBgRDxrCgZErJojebb1pYivt4ei9D8NYQkbg9U" // have balance
	// address := "5UXKnBuqVdoBgRDxrCg" // error address

	chain := NewPolkaChain(rpcUrl, "")

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

	chain := NewPolkaChain(rpcUrl, "")

	balance, err := chain.QueryBalancePubkey(pubkey)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(balance)
}

func TestXBTCBalance(t *testing.T) {
	rpcUrl := "wss://minichain-testnet.coming.chat" // 非 chainx 会抛出 error
	rpcUrl = "wss://testnet3.chainx.org"            // chainx 才可以正常查询
	// address := "5QUEnWNMDFqsbUGpvvtgWGUgiiojnEpLf7581ELLAQyQ1xnT"
	// address := "5PjZ58jF72pCz6Y3FkB3jtyWbhhEbWxBz8CkDD7NG3yjL6s1"
	address := "5TrfBkZz213mgFL59pxqjGThzFCw2VvgkaKDMZi1pv9yNYCY"

	chain := NewPolkaChain(rpcUrl, "")

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

	chain, err := NewPolkaChainWithRpc(rpcUrl, "", meta)
	if err != nil {
		t.Fatal(err)
	}

	meta2, err := chain.GetMetadataString()
	if meta != meta2 {
		t.Fatal("metadata restore failed")
	}

	t.Log(chain)
}

func TestTransactionDetail(t *testing.T) {
	// rpcUrl := "https://mainnet.chainx.org/rpc"
	// scanUrl := "https://multiscan-api.coming.chat/chainx/extrinsics"
	// hashString := "0xb6dbc48dd686cd52897cc8f4871b406a2c64bf9f1d6f08903400f809d3f1ff75"

	rpcUrl := "https://rpc-minichain.coming.chat"
	scanUrl := "https://multiscan-api-pre.coming.chat/minix/extrinsics"
	hashString := "0x4f9ea1cf8337b27a335ef21f1d2806a0d25cc12a87f843a64102fbecfea77cb3"

	chain := NewPolkaChain(rpcUrl, scanUrl)

	res, err := chain.FetchTransactionDetail(hashString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestKusamaBalance(t *testing.T) {
	rpcUrl := "wss://kusama.api.onfinality.io/public-ws"
	address := "JHnAxcdUANjiszJVpfDCQyn6T8swMKbvrCAgMbyinAQM2Aj"

	chain := NewPolkaChain(rpcUrl, "")

	balance, err := chain.QueryBalance(address)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(balance)
}

func TestMINIScriptHash(t *testing.T) {
	rpcUrl := "https://mainnet.chainx.org/rpc"   // 非 minix 会抛出 error
	rpcUrl = "https://rpc-minichain.coming.chat" // 需要 minix 链才可以正常获取
	address := "5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnH"
	amount := "10000"

	chain := NewPolkaChain(rpcUrl, "")

	scriptHash, err := chain.FetchScriptHashForMiniX(address, amount)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(scriptHash)
}

func TestEstimatedFee(t *testing.T) {
	rpcUrl := "https://rpc-minichain.coming.chat"
	// rpcUrl := "https://mainnet.chainx.org/rpc"
	// rpcUrl := "https://mainnet.sherpax.io/rpc"
	address := "5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnH"
	amount := "10000"

	chain := NewPolkaChain(rpcUrl, "")
	metadata, _ := chain.GetMetadataString()
	tx, _ := NewTx(metadata)
	transaction, _ := tx.NewBalanceTransferTx(address, amount)

	fee, err := chain.EstimateFeeForTransaction(transaction)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fee)
}
