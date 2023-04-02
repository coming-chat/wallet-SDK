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

	// toAddress := "0x0c61c2622b77e2a9a3c953690e915ab82d6370d9"
	// amount := "8000000"
	toAddress := M2Account(t).Address()
	amount := strconv.FormatUint(uint64(0.01e9), 10)

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
	// digest := "4nMHqXi60PLxj/DxLCWwkiO3L41kIz89qMDEpStRdP8="
	// digest := "hPOfmwiRRsxleD0JGA67bWFBur+z1BdbLo6vYxzB+9w=" // normal coin transfer
	// digest := "uJYpq7vh/3dI4tzmc5wsecUGTMzYiae4829C1VBuQHM=" // testnet nft transfer
	digest := "4SjJh84sJ2KC56ULKv42Ftq6r9RHSuBqq64vraGvosmJ" // devnet sui transfer
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
	signature, err := account.account.SignSecure(txn.TxBytes.Data(), sui_types.DefaultIntent())
	require.Nil(t, err)

	detail, err := client.ExecuteTransactionBlock(context.Background(), txn.TxBytes, []any{signature}, &types.SuiTransactionBlockResponseOptions{ShowEffects: true}, types.TxnRequestTypeWaitForEffectsCert)
	require.Nil(t, err)
	t.Log(detail)
}

func TestFaucet(t *testing.T) {
	address := "0x6c5d2cd6e62734f61b4e318e58cbfd1c4b99dfaf"
	// address = "0x30d903963ceb4a5c74de9f87498fb467cae72008"
	digest, err := FaucetFundAccount(address, "")
	if err != nil {
		t.Logf("error = %v", err)
	} else {
		t.Logf("digest = %v", digest)
	}
}
