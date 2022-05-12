package cosmos

import (
	"context"
	"strconv"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type chainInfo struct {
	rpc  string
	rest string
	scan string
}

type chainConfig struct {
	cosmosProd chainInfo
	cosmosTest chainInfo
}

var rpcs = chainConfig{
	cosmosProd: chainInfo{
		rpc:  "https://cosmos-mainnet-rpc.allthatnode.com:26657",
		rest: "https://cosmos-mainnet-rpc.allthatnode.com:1317",
		scan: "https://www.mintscan.io/cosmos",
	},
	cosmosTest: chainInfo{
		rpc:  "https://cosmos-testnet-rpc.allthatnode.com:26657",
		rest: "https://cosmos-testnet-rpc.allthatnode.com:1317",
		scan: "https://cosmoshub-testnet.mintscan.io/cosmoshub-testnet",
	},
}

// $request cosmos1unek4dqvkwxv6sfrakk4903m0gmxkfyeprcqtg  theta

func TestQueryBalance(t *testing.T) {
	// req := banktypes.QueryBalanceRequest{
	// 	Address: accountCase1.address,
	// 	Denom:   "",
	// }
	// cli := banktypes.NewQueryClient(grpc.ClientConn)
	// res, err := cli.Balance(context.Background(), &req)
	// t.Log(res, err)

	c := NewChainWithRpc(rpcs.cosmosTest.rpc, rpcs.cosmosTest.rest)
	acc, err := c.AccountOf(accountCase1.address)
	t.Log(acc, err)
}

func TestTransssss(t *testing.T) {
	rpcinfo := rpcs.cosmosTest

	from := accountCase1.mnemonic
	account, _ := NewCosmosAccountWithMnemonic(from)

	toAddress := accountCase2.address
	gasPrice := GasPriceLow
	gasLimit := GasLimitDefault
	amount := "1000"

	chain := NewChainWithRpc(rpcinfo.rpc, rpcinfo.rest)
	token := chain.DenomToken("cosmos", "uatom")

	signedTx, err := token.BuildTransferTxWithAccount(account, toAddress, gasPrice, gasLimit, amount)
	if err != nil {
		t.Fatal(err)
	}

	txHash, err := chain.SendRawTransaction(signedTx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash, "at2", rpcinfo.scan)
}

func TestTransfer(t *testing.T) {
	rpcinfo := rpcs.cosmosTest

	from := accountCase1.mnemonic
	account, _ := NewCosmosAccountWithMnemonic(from)
	privKey := account.privKey

	toAddress := accountCase2.address
	toAccAddress, err := AccAddressFromBech32(toAddress, "cosmos")
	if err != nil {
		t.Fatal(err)
	}

	chain := NewChainWithRpc(rpcinfo.rpc, rpcinfo.rest)
	client, err := chain.GetClient()
	if err != nil {
		t.Fatal(err)
	}

	accountInfo, err := chain.AccountOf(account.Address())
	if err != nil {
		t.Fatal(err)
	}
	accountNumber, err := strconv.ParseUint(accountInfo.AccountNumber, 10, 64)
	if err != nil {
		t.Fatal(accountNumber)
	}
	sequence, err := strconv.ParseUint(accountInfo.Sequence, 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	blockInfo, err := client.Block(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	chainId := blockInfo.Block.ChainID
	latestHeight := blockInfo.Block.Height

	transferCoin := sdk.NewInt64Coin("uatom", 1000)
	feeCoin := sdk.NewInt64Coin("uatom", 100)

	encCfg := simapp.MakeTestEncodingConfig()

	// Create a new TxBuilder.
	txBuilder := encCfg.TxConfig.NewTxBuilder()
	msg1 := banktypes.NewMsgSend(privKey.PubKey().Address().Bytes(), toAccAddress.Bytes(), sdk.NewCoins(transferCoin))
	err = txBuilder.SetMsgs(msg1)
	if err != nil {
		t.Fatal(err)
	}

	txBuilder.SetGasLimit(110000)
	txBuilder.SetFeeAmount(sdk.NewCoins(feeCoin))
	// txBuilder.SetMemo("bridge")
	txBuilder.SetTimeoutHeight(uint64(latestHeight) + 1000)

	sigV2 := signing.SignatureV2{
		PubKey: privKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  encCfg.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: sequence,
	}

	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		t.Fatal(err)
	}

	signerData := xauthsigning.SignerData{
		ChainID:       chainId,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
	sigV2, err = clientTx.SignWithPrivKey(
		encCfg.TxConfig.SignModeHandler().DefaultMode(), signerData,
		txBuilder, privKey, encCfg.TxConfig, sequence)
	if err != nil {
		t.Fatal(err)
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		t.Fatal(err)
	}

	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		t.Fatal(err)
	}

	signedTx := types.HexEncodeToString(txBytes)

	hashString, err := chain.SendRawTransaction(signedTx)
	if err != nil {
		return
	}
	t.Log(hashString, "at2", rpcinfo.scan)
}
