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

func TestContractGasLimit(t *testing.T) {
	// rpcUrl := rpcs.binanceTest.url
	// contractAddress := "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	// walletAddress := "0x1988EbF818FF475f680AF72cf44BBEe1A7CEA666"
	// toAddress := "0x1988EbF818FF475f680AF72cf44BBEe1A7CEA666"
	// amount := "111000000"
	id := 2
	println(id)

	rpcInfo := rpcs.ethereumProd
	// rpcUrl := rpcs.binanceProd.url
	// walletAddress := "0x46D608080FF930D847185Ea6811CC0652457E76c"
	walletAddress := "0x1F05e1419D511C5f1Df9a624FC31Afe24170b4A2"
	toAddress := "0x1F05e1419D511C5f1Df9a624FC31Afe24170b4A2"
	amount := "10000"

	u := &CoinUtil{
		RpcUrl:          rpcInfo.url,
		ContractAddress: rpcInfo.contracts.USDT,
		WalletAddress:   walletAddress,
	}
	gasLimit, err := u.EstimateGasLimit(toAddress, "34891", amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("gas  limit = ", gasLimit)
}
