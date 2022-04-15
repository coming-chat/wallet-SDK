package eth

import "testing"

func TestUtilQueryBalance(t *testing.T) {
	util := NewCoinUtilWithRpc(binanceTestRpcUrl, "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee", "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee")
	// util := NewCoinUtilWithRpc(binanceTestRpcUrl, "", "0xed24fc36d5ee211ea25a80239fb8c4cfd80f12ee")

	balance, err := util.QueryBalance()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)

	info, err := util.CoinInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}
