package starknet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/artifacts"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
	"github.com/stretchr/testify/require"
)

const (
	env = "testnet"
	// compiledOZAccount     = "./contracts/account/OZAccount_compiled.json"
	// compiledERC20Contract = "./contracts/erc20/erc20_custom_compiled.json"
	// predeployedContract = "0x0024e9f35c5d6a14dcbb3f08be5fb7703e76611767266068c521fe8cba27983c"
	maxPoll      = 15
	pollInterval = 5
)

func M1Account(t *testing.T) *Account {
	mnemonic := testcase.M1
	acc, err := NewAccountWithMnemonic(mnemonic)
	require.Nil(t, err)
	return acc
}

func TestEncodeAddress(t *testing.T) {
	{
		// ArgentX
		// see https://testnet.starkscan.co/tx/0x04001e02d0397ebd0821f7c6865f0914e293955417a33505f99e0e1ec1182ea3
		pub := "0x28081ae2bc3668241b1303df98a61e229ee760eb554f9c7fb21cd968a1b74b1"
		param := deployParamForArgentX(*mustFelt(pub))
		addr, err := param.ComputeContractAddress()
		require.Nil(t, err)
		require.Equal(t, addr.String(), "0x7384b9770dce88ee83a62a8a0ab0fac476e513a9e4b611b80fa08e844ce1f2")
	}
	{
		// Braavos
		// see https://testnet.starkscan.co/tx/0x2d72531b049bcf72dbaa4730161e082798e10fa849763f12b3788f7c275b682
		pub := "0x28081ae2bc3668241b1303df98a61e229ee760eb554f9c7fb21cd968a1b74b1"
		param := deployParamForBraavos(*mustFelt(pub))
		addr, err := param.ComputeContractAddress()
		require.Nil(t, err)
		require.Equal(t, addr.String(), "0x8debaf4740ac184b2e879d4d3fd773f2c7f5d453b795212d4098899a73fc19")
	}
}

func TestAccount(t *testing.T) {
	mnemonic := testcase.M1
	account, err := NewAccountWithMnemonic(mnemonic)
	require.Nil(t, err)

	prikey, err := account.PrivateKeyHex()
	require.Nil(t, err)

	account2, err := AccountWithPrivateKey(prikey)
	require.Nil(t, err)
	require.Equal(t, account.PublicKey(), account2.PublicKey())
	require.Equal(t, account.Address(), account2.Address())

	t.Log(prikey)
	t.Log(account.PublicKeyHex())
	t.Log(account.Address())

	require.Equal(t, account.Address(), "0x8debaf4740ac184b2e879d4d3fd773f2c7f5d453b795212d4098899a73fc19")
}

func TestAccount_ImportPrivateKey(t *testing.T) {
	priHex := "0x1234567890"
	priDecimal := "78187493520"
	require.Equal(t, types.HexToBN(priHex), types.StrToBig(priDecimal))

	accountHex, err := AccountWithPrivateKey(priHex)
	require.Nil(t, err)
	accountDecimal, err := AccountWithPrivateKey(priDecimal)
	require.Nil(t, err)

	require.Equal(t, accountHex.PublicKey(), accountDecimal.PublicKey())
	require.Equal(t, accountHex.Address(), accountDecimal.Address())
	require.Equal(t, accountHex.Address(), "0x7d090c124f2cac618e5b53ad97cdb204debc61e9fc63f94d63f4f75a183ceef")
}

func TestGrindKey(t *testing.T) {
	prikey := "86F3E7293141F20A8BAFF320E8EE4ACCB9D4A4BF2B4D295E8CEE784DB46E0519"
	seed, ok := big.NewInt(0).SetString(prikey, 16)
	require.True(t, ok)
	res, err := grindKey(seed.Bytes())
	require.Nil(t, err)
	require.Equal(t, res.Text(16), "5c8c8683596c732541a59e03007b2d30dbbbb873556fe65b5fb63c16688f941")
}

func TestAccount22(t *testing.T) {
	gw := gateway.NewClient(gateway.WithChain(env))

	// privateKey, err := caigo.Curve.GetRandomPrivateKey()
	// if err != nil {
	// 	fmt.Println("can't get random private key:", err)
	// 	os.Exit(1)
	// }
	// privateKey := types.SNValToBN("2522809042406563759994430227158961551351879727850354505141300412654704193860") // argent x
	privateKey := types.SNValToBN("0x072ea243466a8e5e3eea63900c7a0a2d7afcd82aa6a6402b1bdd409000020667") // braavos
	pubX, _, err := caigo.Curve.PrivateToPoint(privateKey)
	if err != nil {
		fmt.Println("can't generate public key:", err)
		os.Exit(1)
	}

	contractClass := types.ContractClass{}
	err = json.Unmarshal(artifacts.AccountCompiled, &contractClass)
	if err != nil {
		fmt.Println("could not log file", err)
		os.Exit(1)
	}

	txxx, err := gw.Declare(context.Background(), contractClass, gateway.DeclareRequest{})
	if err != nil {
		t.Fatalf("could not declare contract: %v\n", err)
	}

	fmt.Println("Deploying account to testnet. It may take a while.")
	accountResponse, err := gw.DeployAccount(context.Background(), types.DeployAccountRequest{
		Type:                gateway.DEPLOY,
		ContractAddressSalt: types.BigToHex(pubX),
		ConstructorCalldata: []string{pubX.String()},
		MaxFee:              big.NewInt(1000000000000),
		ClassHash:           txxx.ClassHash,
		Version:             1,
	})
	// accountResponse, err := gw.Deploy(context.Background(), contractClass, types.DeployRequest{
	// 	Type:                gateway.DEPLOY,
	// 	ContractAddressSalt: types.BigToHex(pubX),     // salt to hex
	// 	ConstructorCalldata: []string{pubX.String()}}) // public key
	if err != nil {
		fmt.Println("can't deploy account:", err)
		os.Exit(1)
	}

	// if err := waitForTransaction(gw, accountResponse.TransactionHash); err != nil {
	// 	fmt.Println("Account deployement transaction failure:", err)
	// 	os.Exit(1)
	// }

	tx, err := gw.Transaction(context.Background(), gateway.TransactionOptions{TransactionHash: accountResponse.TransactionHash})
	if err != nil {
		fmt.Println("can't fetch transaction data:", err)
		os.Exit(1)
	}

	t.Log(tx.Transaction.ContractAddress)

	account, err := caigo.NewGatewayAccount(privateKey.String(), tx.Transaction.ContractAddress, &gateway.GatewayProvider{Gateway: *gw})
	if err != nil {
		fmt.Println("can't create account:", err)
		os.Exit(1)
	}
	t.Log(account.AccountAddress)

	// ba, err := balanceOf(gw, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7", "0x0023C4475F2f2355580f5994294997d3A18237ef62223D20C41876556327A05E")
	// require.Nil(t, err)
	// t.Log(ba)

	// fmt.Println("Account deployed. Contract address: ", account.Address)
	// if err := savePrivateKey(types.BigToHex(privateKey)); err != nil {
	// 	fmt.Println("can't save private key:", err)
	// 	os.Exit(1)
	// }
}

// Utils function to wait for transaction to be accepted on L2 and print tx status.
func waitForTransaction(gw *gateway.Gateway, transactionHash string) error {
	acceptedOnL2 := false
	var receipt *gateway.TransactionReceipt
	var err error
	fmt.Println("Polling until transaction is accepted on L2...")
	for !acceptedOnL2 {
		_, receipt, err = gw.WaitForTransaction(context.Background(), transactionHash, pollInterval, maxPoll)
		if err != nil {
			fmt.Println(receipt.Status)
			return fmt.Errorf("Transaction Failure (%s): can't poll to desired status: %s", transactionHash, err.Error())
		}
		fmt.Println("Current status : ", receipt.Status)
		if receipt.Status == "ACCEPTED_ON_L2" { // types.ACCEPTED_ON_L2.String() {
			acceptedOnL2 = true
		}
	}
	return nil
}

// balanceOf returns the balance of the account at the accountAddress address.
func balanceOf(gw *gateway.Gateway, erc20address, accountAddress string) (string, error) {
	res, err := gw.Call(context.Background(), types.FunctionCall{
		ContractAddress:    types.HexToHash(erc20address),
		EntryPointSelector: "balanceOf",
		Calldata: []string{
			types.HexToBN(accountAddress).String(),
		},
	}, "")
	if err != nil {
		return "", fmt.Errorf("can't call erc20: %s. Error: %w", accountAddress, err)
	}
	low := types.StrToFelt(res[0])
	hi := types.StrToFelt(res[1])

	balance, err := types.NewUint256(low, hi)
	if err != nil {
		return "", nil
	}
	return balance.String(), nil
}
