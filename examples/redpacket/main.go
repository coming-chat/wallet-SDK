package main

import (
	"os"

	"github.com/coming-chat/wallet-SDK/core/aptos"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const (
	testNetUrl = "https://fullnode.devnet.aptoslabs.com"
)

func main() {
	chain := aptos.NewChainWithRestUrl(testNetUrl)
	account, err := aptos.NewAccountWithMnemonic(os.Getenv("mnemonic"))
	if err != nil {
		panic(err)
	}
	contract := aptos.NewRedPacketContract(chain, os.Getenv("red_packet"))
	action, err := base.NewRedPacketActionCreate("", 3, "1000")
	if err != nil {
		panic(err)
	}
	txHash, err := action.Do(chain, account, contract, "")
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
