package main

import (
	"os"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/eth"
)

func main() {
	chain := eth.NewChainWithRpc(os.Getenv("rpc"))
	account, err := eth.NewAccountWithMnemonic(os.Getenv(""))
	if err != nil {
		panic(err)
	}
	contract := eth.NewRedPacketContract(chain, os.Getenv("red_packet"))
	action, err := base.NewRedPacketActionCreate("", 5, "100000")
	if err != nil {
		panic(err)
	}
	txHash, err := contract.SendTransaction(account, action)
	if err != nil {
		panic(err)
	}
	txDetail, err := chain.FetchTransactionDetail(txHash)
	if err != nil {
		panic(err)
	}
	println(txHash)
	println(txDetail.Status)
}
