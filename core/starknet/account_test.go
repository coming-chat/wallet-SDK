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
