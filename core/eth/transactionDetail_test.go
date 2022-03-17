package eth

import (
	"testing"
)

func TestTransactionDetail(t *testing.T) {
	var ethChain, _ = NewEthChain().CreateRemote(ethRpcUrl)

	// 手续费 为 0 的
	// hashString := "0xac90c27075e3843ef43f06bdbca6426830fa547e71bc2ab024e13f3eaef57a49"

	// 转账失败的
	// hashString := "0x4be8dd93de963634ec3170d85859a0072dbf160efe878ee9f025395e619362c4"

	// Out of gas
	// hashString := "0x9ef27c4983b18fd25a149e737feefb952889253aa6e2cddb62c6cf80a23887c3"
	hashString := "0x7571b8a1bd77d426395365601dd10d051dfb61914e2ed6c37b2fe7045cf96d47"

	detail, err := ethChain.FetchTransactionDetail(hashString)
	if err != nil {
		t.Fatal("fetch detail failure:", err)
	}

	t.Log(detail.JsonString())
}

func TestTransaction(t *testing.T) {
	rpcUrl := "https://mainnet.sherpax.io/rpc"
	hashString := "0xee634a4a4152e018fbc9af27dd4a6791a4e74e852fc769173aba8bb3339fb089"

	// rpcUrl := ethRpcUrl

	var ethChain, _ = NewEthChain().CreateRemote(rpcUrl)
	detail, err := ethChain.FetchTransactionDetail(hashString)
	if err != nil {
		t.Fatal("fetch detail failure:", err)
	}

	t.Log("success")
	t.Log(detail.JsonString())
}
