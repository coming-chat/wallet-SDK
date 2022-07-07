package solana

import (
	"context"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/portto/solana-go-sdk/rpc"
)

func newChainAndAccount() (*Chain, *Account) {
	// chain := NewChainWithRpc(rpc.MainnetRPCEndpoint)
	chain := NewChainWithRpc(rpc.DevnetRPCEndpoint)
	// c := client.NewClient(rpc.LocalnetRPCEndpoint)
	account, _ := NewAccountWithMnemonic(testcase.M1)
	return chain, account
}

func TestAirdrop(t *testing.T) {
	chain, account := newChainAndAccount()
	txhash, err := chain.client().RequestAirdrop(context.Background(), account.Address(), 1e9)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(txhash)
}

func TestBalance(t *testing.T) {
	chain, acc := newChainAndAccount()

	balance, err := chain.client().GetBalance(context.Background(), acc.Address())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)
}

func TestBuildtxAndSendTransaction(t *testing.T) {
	chain, acc := newChainAndAccount()
	token := &Token{chain: chain}

	receiver := "9B5XszUGdMaxCZ7uSQhPzdks5ZQSmWxrmzCSvtJ6Ns6g"
	amount := "10000000"
	// amount := "1879985000"

	signedTx, err := token.BuildTransferTxWithAccount(acc, receiver, amount)
	if err != nil {
		t.Fatal(err)
	}

	txHash, err := chain.SendRawTransaction(signedTx.Value)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}

func TestEstimateFee(t *testing.T) {
	chain := NewChainWithRpc(rpc.DevnetRPCEndpoint)
	token := &Token{chain: chain}

	receiver := "9B5XszUGdMaxCZ7uSQhPzdks5ZQSmWxrmzCSvtJ6Ns6g"
	amount := "10000000"
	fee, err := token.EstimateFees(receiver, amount)
	t.Log(fee, err)
}

func TestTransactionDetail(t *testing.T) {
	chain := NewChainWithRpc(rpc.DevnetRPCEndpoint)
	txhash := "2CLc8VSHsS67JT6a4UQZMMwFzcgnriHL4d1rpwu9WBwAqzucj6XPBcL2AYJy7n6xvmrnXTgGRKoThHizN8E8NTFN" // 我的普通的
	// txhash := "2LkRwB9QyttUYmCjfCx39QyMHWzcWrDsM2nkjbK4EN2AGmTeZW6TBNfe5rQ9pjMgCsVecGz1a9vPmXNsSy7gTYJQ" // 错误的
	// txhash := "33yLja4FF9ZZQWMwdb72KroV4qQrmjTiKHcrR4KNRT3rK75mY2TPVQxMs9YGiQPhK5vxxqS5d5sjfJXCM2E8urFB" // 有两个人扣钱的
	// txhash := "3eETnjk4jSTudYe3tyZR7VKd9E5gy3r9c149f78drfJoq52yo4bnJGmyU5NNvpxF3JmQNYHF2SA8RxmNYySiVkgN" // 有两人扣钱，两人收钱
	response, err := chain.client().GetTransaction(context.Background(), txhash)
	if err != nil {
		t.Fatal(err)
	}

	detail := &base.TransactionDetail{HashString: txhash}
	err = decodeTransaction(response, detail)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(detail)
}
