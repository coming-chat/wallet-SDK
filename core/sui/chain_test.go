package sui

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/coming-chat/go-sui/account"
	"github.com/coming-chat/go-sui/sui_types"
	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

const DevnetRpcUrl = "https://fullnode.devnet.sui.io"

const TestnetRpcUrl = "https://fullnode.testnet.sui.io"

// const TestnetRpcUrl = "https://sui-testnet-wave3.coming.chat"

const MainnetRpcUrl = "https://sui-mainnet.coming.chat"

func DevnetChain() *Chain {
	return NewChainWithRpcUrl(DevnetRpcUrl)
}

func TestnetChain() *Chain {
	return NewChainWithRpcUrl(TestnetRpcUrl)
}

func MainnetChain() *Chain {
	return NewChainWithRpcUrl(MainnetRpcUrl)
}

func TestEstimateGas(t *testing.T) {
	account := M1Account(t)
	chain := DevnetChain()
	token := NewTokenMain(chain)

	toAddress := M2Account(t).Address()
	amount := strconv.FormatUint(uint64(0.01e9), 10)

	txn, err := token.BuildTransferTransaction(account, toAddress, amount)
	require.Nil(t, err)

	fee, err := chain.EstimateGasFee(txn)
	require.Nil(t, err)

	t.Log("gas fee = ", fee.Value)
}

func TestFetchTransactionDetail(t *testing.T) {
	// digest := "3aFbrGBfi9A5ZSjv9jcEwx8TQjm1XC8NqWvSkzKJEbVE" // normal transfer
	// digest := "C9grwYWbJyBypSbgXEMaQ47LJ2uy3bToQLtqA9cVee2z" // not coin transfer
	// digest := "29MYmpk3kzcmB6e7FMwe6mD7x5pqDCeRoRvhJDFnXvAX"
	// digest := "FD4onoYMKTNC4f7UFS4UmeaeDKsqt73eaRciDm7UcEdZ"
	digest := "GH87s7pc8EWhnuq96tGe34he12hox6TuK5JPzqGJvm8S" // transfer object
	// digest := "5PLq48GYcKwsA3P1rDpUbWrNLX63tp2iamYJDhTHhskC" // pay sui
	chain := TestnetChain()

	detail, err := chain.FetchTransactionDetail(digest)
	require.Nil(t, err)

	t.Log(detail)
}

func TestSplit(t *testing.T) {
	account := M1Account(t)
	chain := TestnetChain()

	client, err := chain.Client()
	require.Nil(t, err)

	signer, err := types.NewAddressFromHex(account.Address())
	require.Nil(t, err)
	coins, err := client.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(coins.Data), 1)

	coin := coins.Data[0]

	gasBudget := SUI(0.01).Uint64()
	txn, err := client.SplitCoinEqual(context.Background(), *signer, coin.CoinObjectId, types.NewSafeSuiBigInt[uint64](2), nil, types.NewSafeSuiBigInt(gasBudget))
	require.Nil(t, err)

	simulateCheck(t, chain, txn, false)
}

func TestFaucet(t *testing.T) {
	address := "0x7e875ea78ee09f08d72e2676cf84e0f1c8ac61d94fa339cc8e37cace85bebc6e"
	digest, err := FaucetFundAccount(address, "")
	if err != nil {
		t.Logf("error = %v", err)
	} else {
		t.Logf("digest = %v", digest)
	}
}

func Test_EstimateFee_And_RebuildTxn(t *testing.T) {
	account := M1Account(t)
	chain := TestnetChain()
	token := NewTokenMain(chain)

	txn, err := chain.EstimateTransactionFeeAndRebuildTransaction(MinGasBudget, func(gasBudget uint64) (*Transaction, error) {
		return token.BuildTransferTransaction(account, account.Address(), SUI(1).String())
	})
	require.Nil(t, err)

	t.Log(txn)
}

func Test_TryToFindTheMiniumGasBudget(t *testing.T) {
	chain := MainnetChain()

	addressHasManySUIObject := "0xab73ad38c63f83eda02182422b545395be1d3caeb54b5869159a9f70b678cd56"
	totalBigAmount := SUI(266000000).String()

	req, err := chain.BuildMergeCoinRequest(addressHasManySUIObject, "", totalBigAmount)
	require.Nil(t, err)
	_, err = chain.BuildMergeCoinPreview(req)
	require.Nil(t, err)
	// 我们需要修改下面方法的第一次构建的 maxGasBudget 的值，保证它在第一次模拟交易能通过
	// chain.EstimateTransactionFeeAndRebuildTransaction()
}

func Test_TryBuild(t *testing.T) {
	chain := MainnetChain()

	address := "0xab73ad38c63f83eda02182422b545395be1d3caeb54b5869159a9f70b678cd56"

	_, err := chain.AddDelegation(address, SUI(2).String(), ComingChatValidatorMainnet)
	require.Nil(t, err)
}

func simulateCheck(t *testing.T, chain *Chain, txn *types.TransactionBytes, showJson bool) *types.DryRunTransactionBlockResponse {
	cli, err := chain.Client()
	require.Nil(t, err)
	resp, err := cli.DryRunTransaction(context.Background(), txn)
	require.Nil(t, err)
	require.Equal(t, resp.Effects.Data.V1.Status.Error, "")
	require.True(t, resp.Effects.Data.IsSuccess())
	if showJson {
		data, err := json.Marshal(resp)
		require.Nil(t, err)
		respStr := string(data)
		t.Log("simulate run resp: ", respStr)
	}
	return resp
}

func executeTransaction(t *testing.T, chain *Chain, txn *types.TransactionBytes, acc *account.Account) *types.SuiTransactionBlockResponse {
	// firstly we best ensure the transaction simulate call can be success.
	simulateCheck(t, chain, txn, false)

	// execute
	cli, err := chain.Client()
	require.NoError(t, err)
	signature, err := acc.SignSecureWithoutEncode(txn.TxBytes, sui_types.DefaultIntent())
	require.NoError(t, err)
	options := types.SuiTransactionBlockResponseOptions{
		ShowEffects:        true,
		ShowBalanceChanges: true,
		ShowObjectChanges:  true,
		ShowInput:          true,
		ShowEvents:         true,
	}
	resp, err := cli.ExecuteTransactionBlock(
		context.TODO(), txn.TxBytes, []any{signature}, &options,
		types.TxnRequestTypeWaitForLocalExecution,
	)
	require.NoError(t, err)
	return resp
}
