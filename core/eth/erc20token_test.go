package eth

import "testing"

func TestErc20(t *testing.T) {
	e, _ := NewEthChain().CreateRemote(binanceTestRpcUrl)
	contractAddress := "0x6cd2bf22b3ceadff6b8c226487265d81164396c5"

	tokenName := ""
	err := e.CallContractConstant(
		&tokenName,
		contractAddress,
		Erc20AbiStr,
		"name",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tokenName)

}

func TestBatch(t *testing.T) {
	c := NewChainWithRpc(binanceTestRpcUrl)
	array := []string{
		// "0xE7e312dfC08e060cda1AF38C234AEAcc7A982143", // 报错
		// "0x4B53739D798EF0BEa5607c254336b40a93c75b52", // 报错
		// "0x935CC842f220CF3A7D10DA1c99F01B1A6894F7C5", // 报错
		// "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56", // 报错
		"0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee", // 可请求余额
	}

	address := "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"
	result, err := c.BatchErc20TokenBalance(array, address)
	if err != nil {
		t.Fatal(err)
	}

	for i, x := range result {
		t.Log(i, x)
	}
}

func TestSdkBatch(t *testing.T) {
	e, _ := NewEthChain().CreateRemote(binanceTestRpcUrl)
	contract := "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"
	address := "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee"
	result, err := e.SdkBatchTokenBalance(contract, address)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(result)
}
