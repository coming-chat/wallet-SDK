package eth

import (
	"testing"
)

func TestConnect(t *testing.T) {
	chain, err := NewEthChain().CreateRemote(rpcs.ethereumProd.url)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(chain)
	// for i := 0; i < 1; i++ {
	// 	time.Sleep(1 * time.Second)
	// 	address := "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"
	// 	balance, err := chain.TokenBalance(address, address)
	// 	if err != nil {
	// 		t.Log("...... catched err", err)
	// 	} else {
	// 		t.Log("...... balance", balance)
	// 	}
	// }

	// t.Log("should successd connect", chain)
}

func TestBatchStatus(t *testing.T) {
	var ethChain, _ = NewEthChain().CreateRemote(rpcs.binanceProd.url)

	var hashStrings = "0x9ef27c4983b18fd25a149e737feefb952889253aa6e2cddb62c6cf80a23887c3,0x39fa91a5e34d50f373339b2a5e9102ffc2c321f497a49841c62fd213e433290d,0x2cbf78965bbddecf86d2d0fb17069fa760fa652d81ee79d9a99f0add92b05364"

	var statuses = ethChain.SdkBatchTransactionStatus(hashStrings)

	t.Log(statuses)
}

const (
	transferFromAddress = "0x8de5ff2eded4d897da535ab0f379ec1b9257ebab"
	transferToAddress   = "0x6cd2bf22b3ceadff6b8c226487265d81164396c5"
)

func TestEstimateGasLimit(t *testing.T) {
	var ethChain, _ = NewEthChain().CreateRemote(rpcs.binanceTest.url)
	gasprice := "10"
	amount := "1"
	gasLimit, err := ethChain.EstimateGasLimit(transferFromAddress, transferToAddress, gasprice, amount)
	if err != nil {
		t.Fatal("gas---" + err.Error())
	}

	t.Log("TestEstimateGasLimit success", gasLimit)
}
