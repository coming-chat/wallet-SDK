package eth

import (
	"testing"
)

const (
	transferFromAddress = "0x8de5ff2eded4d897da535ab0f379ec1b9257ebab"
	transferToAddress   = "0x6cd2bf22b3ceadff6b8c226487265d81164396c5"
)

func TestEstimateGasLimit(t *testing.T) {
	var ethChain, _ = NewEthChain().CreateRemote(ethRpcUrl)
	gasprice := "10"
	amount := "1"
	gasLimit, err := ethChain.EstimateGasLimit(transferFromAddress, transferToAddress, gasprice, amount)
	if err != nil {
		t.Fatal("gas---" + err.Error())
	}

	t.Log("TestEstimateGasLimit success", gasLimit)
}

func TestContractGasLimit(t *testing.T) {
	// rpcUrl := binanceTestRpcUrl
	// contractAddress := "0xdAC17F958D2ee523a2206206994597C13D831ec7"
	// walletAddress := "0x1988EbF818FF475f680AF72cf44BBEe1A7CEA666"
	// toAddress := "0x1988EbF818FF475f680AF72cf44BBEe1A7CEA666"
	// amount := "111000000"
	id := 2
	println(id)

	rpcUrl := ethMainProdRpcUrl
	// rpcUrl := binanceProdRpcUrl
	contractAddress := contractUSDT
	// walletAddress := "0x46D608080FF930D847185Ea6811CC0652457E76c"
	walletAddress := "0x1F05e1419D511C5f1Df9a624FC31Afe24170b4A2"
	toAddress := "0x1F05e1419D511C5f1Df9a624FC31Afe24170b4A2"
	amount := "10000"

	u := &CoinUtil{
		RpcUrl:          rpcUrl,
		ContractAddress: contractAddress,
		WalletAddress:   walletAddress,
	}
	gasLimit, err := u.EstimateGasLimit(toAddress, "", amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("gas  limit = ", gasLimit)
}
