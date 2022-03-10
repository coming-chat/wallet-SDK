package eth

import (
	"testing"
)

const (
	transferFromAddress = "0x8de5ff2eded4d897da535ab0f379ec1b9257ebab"
	transferToAddress   = "0x6cd2bf22b3ceadff6b8c226487265d81164396c5"
)

func TestEstimateGasLimit(t *testing.T) {
	gasprice := "10"
	amount := "1"
	gasLimit, err := ethChain.EstimateGasLimit(transferFromAddress, transferToAddress, gasprice, amount)
	if err != nil {
		t.Fatal("gas---" + err.Error())
	}

	t.Log("TestEstimateGasLimit success", gasLimit)
}
