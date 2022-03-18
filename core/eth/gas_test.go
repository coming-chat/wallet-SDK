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
	u := &CoinUtil{
		RpcUrl:          binanceTestRpcUrl,
		ContractAddress: "0xdAC17F958D2ee523a2206206994597C13D831ec7",
		WalletAddress:   "0x1988EbF818FF475f680AF72cf44BBEe1A7CEA666",
	}
	address := "0x1988EbF818FF475f680AF72cf44BBEe1A7CEA666"
	amount := "111000000"

	gasLimit, err := u.EstimateGasLimit(address, "", amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("gas  limit = ", gasLimit)
}
