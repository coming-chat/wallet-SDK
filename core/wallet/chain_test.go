package wallet

import (
	"testing"
	"time"
)

const (
	rpcChainxProd  = "https://mainnet.chainx.org/rpc"
	rpcChainxTest  = "https://testnet3.chainx.org/rpc"
	rpcMinixProd   = "https://minichain-mainnet.coming.chat/rpc"
	rpcMinixTest   = "https://rpc-minichain.coming.chat"
	rpcSherpaxProd = "https://mainnet.sherpax.io/rpc"
	rpcSherpaxTest = "https://sherpax-testnet.chainx.org/rpc"
)

func TestQueryBalance(t *testing.T) {
	// rpcUrl := "wss://testnet3.chainx.org"
	// address := "5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnH" // no balance
	// address := "5UXKnBuqVdoBgRDxrCgZErJojebb1pYivt4ei9D8NYQkbg9U" // have balance
	// address := "5UXKnBuqVdoBgRDxrCg" // error address

	rpcUrl := "https://mainnet.sherpax.io/rpc"
	address := "5QNvL6E6qfKBhV2VnvdLbdv2ou4VmU7FDFJ43XvcnuKgzpUp"

	chain := NewPolkaChain(rpcUrl, "x")

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

func TestMultiSignBalanceMINI(t *testing.T) {
	address := "5SeET3a9GgDRennxH89ixezysM5VGAUSArGAv5rJCUqyDpvH"

	chain := NewPolkaChain(rpcMinixProd, "")
	balance, err := chain.QueryBalance(address)
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
	// rpcUrl := "wss://kusama.api.onfinality.io/public-ws"
	// address := "JHnAxcdUANjiszJVpfDCQyn6T8swMKbvrCAgMbyinAQM2Aj"

	rpcUrl := "https://mainnet.sherpax.io/rpc"
	address := "5R1kzaPnLasMiNgigKdjdyddYd2jQg6QqQtbFcZm4RLUuKQY"

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
	rpcUrl = "https://mainnet.chainx.org/rpc"
	rpcUrl = "https://mainnet.sherpax.io/rpc"
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

func TestChainTimeCost(t *testing.T) {
	rpcUrl := "wss://testnet3.chainx.org"

	t.Log("测试.............", "kjnjk")

	start := time.Now()
	c1, _ := newPolkaClient(rpcUrl, "")
	metadata, err := c1.MetadataString()

	cost := time.Since(start)
	t.Log("metadata 是否获取成功：", err == nil)
	t.Log("不使用 metadata 初始化的耗时：", cost)

	t.Log("---------------------------")

	start = time.Now()
	c2, _ := newPolkaClient(rpcUrl, metadata)
	_, err = c2.MetadataString()

	cost = time.Since(start)
	t.Log("metadata 是否获取成功：", err == nil)
	t.Log("使用 metadata 初始化的耗时：", cost)

	t.Log("--------------------------- 再测一次，不使用 metadata 初始化")

	start = time.Now()
	c3, _ := newPolkaClient(rpcUrl, "")
	_, err = c3.MetadataString()

	cost = time.Since(start)
	t.Log("metadata 是否获取成功：", err == nil)
	t.Log("第二次不使用 metadata 初始化的耗时：", cost)
}
