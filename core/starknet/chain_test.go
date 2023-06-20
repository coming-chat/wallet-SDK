package starknet

import (
	"context"
	"math/big"
	"testing"

	"github.com/NethermindEth/juno/utils"
	"github.com/dontpanicdao/caigo/types"
	"github.com/stretchr/testify/require"
)

func MainnetChain() *Chain {
	return NewChainWithRpc(BaseRpcUrlMainnet)
}
func GoerliChain() *Chain {
	return NewChainWithRpc(BaseRpcUrlGoerli)
}

func TestBalance(t *testing.T) {
	owner := "0x0023C4475F2f2355580f5994294997d3A18237ef62223D20C41876556327A05E"
	chain := GoerliChain()

	balance, err := chain.BalanceOf(owner, ETHTokenContractAddressGoerli)
	require.Nil(t, err)
	t.Log(balance.Total)
}

func TestDeployAccount(t *testing.T) {
	acc := M1Account(t)
	chain := GoerliChain()

	pubX := acc.PublicKeyHex()

	txn, err := deployAccountTxnForArgentX(pubX)
	require.Nil(t, err)

	txhash, err := TransactionHash(txn, utils.Network(utils.GOERLI))
	require.Nil(t, err)

	s1, s2, err := acc.SignHash(txhash.BigInt(&big.Int{}))
	require.Nil(t, err)

	txnReq := *parseDeployAccountTransaction(txn)
	txnReq.Signature = types.Signature{s1, s2}

	t.Log(s1.String())
	t.Log(s2.String())

	res, err := chain.gw.DeployAccount(context.Background(), txnReq)
	require.Nil(t, err)
	t.Log(res.ContractAddress)
	t.Log(res.TransactionHash)
}
