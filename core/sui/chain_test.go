package sui

import (
	"context"
	"strconv"
	"testing"

	"github.com/coming-chat/go-sui/sui_types"
	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

const DevnetRpcUrl = "https://fullnode.devnet.sui.io"
const TestnetRpcUrl = "https://fullnode.testnet.sui.io"

func DevnetChain() *Chain {
	return NewChainWithRpcUrl(DevnetRpcUrl)
}

func TestnetChain() *Chain {
	return NewChainWithRpcUrl(TestnetRpcUrl)
}

func TestTransfer(t *testing.T) {
	account := M1Account(t)
	chain := DevnetChain()
	token := NewTokenMain(chain)

	toAddress := M2Account(t).Address()
	amount := strconv.FormatUint(uint64(0.01e9), 10)
	// toAddress := account.Address()
	// amount := strconv.FormatUint(4e9, 10) // test big amount transfer

	signedTxn, err := token.BuildTransferTxWithAccount(account, toAddress, amount)
	require.Nil(t, err)

	hash, err := chain.SendRawTransaction(signedTxn.Value)
	require.Nil(t, err)

	t.Log(hash)
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
	digest := "29MYmpk3kzcmB6e7FMwe6mD7x5pqDCeRoRvhJDFnXvAX"
	chain := DevnetChain()

	detail, err := chain.FetchTransactionDetail(digest)
	require.Nil(t, err)

	t.Log(detail)
}

func TestSplit(t *testing.T) {
	account := M1Account(t)
	chain := DevnetChain()

	client, err := chain.Client()
	require.Nil(t, err)

	signer, _ := types.NewAddressFromHex(account.Address())
	coin := "0x0a1248b37b452627eaa588166464cb84718e9c032da3d986c8a9f7f99c1eb6d8"
	coinID, err := types.NewHexData(coin)
	require.Nil(t, err)

	txn, err := client.SplitCoinEqual(context.Background(), *signer, *coinID, 2, nil, 2000)
	require.Nil(t, err)
	signature, err := account.account.SignSecureWithoutEncode(txn.TxBytes, sui_types.DefaultIntent())
	require.Nil(t, err)

	detail, err := client.ExecuteTransactionBlock(context.Background(), txn.TxBytes, []any{signature}, &types.SuiTransactionBlockResponseOptions{ShowEffects: true}, types.TxnRequestTypeWaitForEffectsCert)
	require.Nil(t, err)
	t.Log(detail)
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
