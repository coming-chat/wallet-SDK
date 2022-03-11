package eth

import (
	"testing"
)

func TestTransferResult(t *testing.T) {
	var ethChain, _ = NewEthChain().CreateRemote(ethRpcUrl)
	// 手续费 为 0 的
	hashString := "0xac90c27075e3843ef43f06bdbca6426830fa547e71bc2ab024e13f3eaef57a49"
	// 转账失败的
	// hashString := "0x4be8dd93de963634ec3170d85859a0072dbf160efe878ee9f025395e619362c4"
	// hash := common.HexToHash(hashString)
	// receipt, err := ethChain.RemoteRpcClient.TransactionReceipt(context.Background(), hash)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	detail, err := ethChain.FetchTransferDetail(hashString)
	if err != nil {
		t.Fatal("fetch detail failure:", err)
	}

	t.Log(detail.JsonString())
}
