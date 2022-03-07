package main

import (
	"fmt"
	"os"

	coin "github.com/coming-chat/wallet-SDK/core/eth"
)

var (
	// 测试网rpc 地址，主网切换 https://bsc-dataseed.binance.org
	rpcUrl = "https://data-seed-prebsc-1-s1.binance.org:8545"

	walletAddress = "0xB553803EE21b486BB86f2A63Bd682529Aa7FCE8D"

	// 转账钱包私钥地址
	privateKey = os.Getenv("privateKeyHex")

	// bsc 测试网 busd 合约地址
	busdContractAddress = "0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee"
)

func main() {
	wallet := coin.NewEthChain()
	wallet.InitRemote(rpcUrl)

	// 获取主网代币 BNB 余额
	balance, _ := wallet.Balance(walletAddress)
	fmt.Printf("bnb balance: %v\n", balance)

	// 获取 Erc20代币 余额
	busdBalance, _ := wallet.TokenBalance(busdContractAddress, walletAddress)

	tokenDecimal, err := wallet.TokenDecimal(busdContractAddress, walletAddress)

	fmt.Printf("busdBalance balance: %v, decimal: %v \n", busdBalance, tokenDecimal)

	if err != nil {
		fmt.Printf("get busdt balance error: %v\n", err)
		return
	}
	nonce, _ := wallet.Nonce(walletAddress)

	// 构造多笔交易则nonce + 1
	callMethodOpts := &coin.CallMethodOpts{
		Nonce: nonce,
	}

	// erc20 代币转账
	buildTxResult, err := wallet.BuildCallMethodTx(
		privateKey,
		busdContractAddress,
		coin.Erc20AbiStr,
		// 调用的合约方法名
		"transfer",
		callMethodOpts,

		// 转账目标地址
		"{\"toAddress\":\"0x178a8AB44b71858b38Cc68f349A06f397A73bFf5\", \"amount\":\"10000000\", \"method\":\"transfer\"}")

	if err != nil {
		fmt.Printf("build call method tx error: %v\n", err)
		return
	}
	// 发送交易
	sendTxResult, err := wallet.SendRawTransaction(buildTxResult.TxHex)
	if err != nil {
		fmt.Printf("send raw transaction error: %v\n", err)
	}
	// 打印交易hash
	fmt.Printf("sendTxResult: %v\n", sendTxResult)

}
