package aptos

import (
	"encoding/json"
	"testing"

	"github.com/coming-chat/go-aptos/aptostypes"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestGenerateTransaction(t *testing.T) {
	payload := payloadDemo()
	bytes, err := json.Marshal(payload)
	require.Nil(t, err)
	payloadJson := string(bytes)

	account, err := NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)
	chain := NewChainWithRestUrl(devnetRestUrl) // devnet

	// GenerateTransaction
	txn1, err := chain.GenerateTransaction(account.PublicKeyHex(), payload)
	require.Nil(t, err)
	txn2json, err := chain.GenerateTransactionJson(account.PublicKeyHex(), payloadJson)
	require.Nil(t, err)

	txn2 := aptostypes.Transaction{}
	err = json.Unmarshal([]byte(txn2json.Value), &txn2)
	require.Nil(t, err)

	txn1.ExpirationTimestampSecs = txn2.ExpirationTimestampSecs
	require.Equal(t, *txn1, txn2)
	t.Log(txn1)
}

func TestSignTransaction(t *testing.T) {
	account, err := NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)
	chain := NewChainWithRestUrl(devnetRestUrl) // devnet

	payload := payloadDemo()
	txn, err := chain.GenerateTransaction(account.PublicKeyHex(), payload)
	require.Nil(t, err)

	signedTxn, err := chain.SignTransaction(account, *txn)
	require.Nil(t, err)
	require.NotNil(t, signedTxn.Signature)
	t.Log(signedTxn)
}

func TestSignMessageJson(t *testing.T) {
	payload := SignMessagePayload{
		Application: false,
		Address:     true,
		ChainId:     true,
		Message:     "ComingChat message",
		Nonce:       2132,
	}
	bytes, err := json.Marshal(payload)
	require.Nil(t, err)
	payloadJson := string(bytes)

	account, err := NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)
	chain := NewChainWithRestUrl(testnetRestUrl)

	// SignMessage
	resp1, err := chain.SignMessage(account, &payload)
	require.Nil(t, err)
	resp2, err := chain.SignMessageJson(account, payloadJson)
	require.Nil(t, err)

	require.Equal(t, resp1, resp2)
	t.Log(resp1)
	t.Log(resp1.JsonString())
}

func payloadDemo() aptostypes.Payload {
	return aptostypes.Payload{
		Function:      "0x1::coin::transfer",
		TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
		Arguments: []any{
			"0x6ed6f83f1891e02c00c58bf8172e3311c982b1c4fbb1be2d85a55562d4085fb1", "100",
		},
	}
}

func TestDapp_GenerateTransactionJson(t *testing.T) {
	payloadJson := `{"function":"0x89576037b3cc0b89645ea393a47787bb348272c76d6941c574b053672b848039::aggregator::three_step_route","type_arguments":["0x1::aptos_coin::AptosCoin","0x1000000fa32d122c18a6a31c009ce5e71674f22d06a581bb0a15575e6addadcc::usda::USDA","0x84d7aeef42d38a5ffc3ccef853e1b82e4958659d16a7de736a29c55fbbeb0114::staked_aptos_coin::StakedAptosCoin","0x5e156f1207d0ebfa19a9eeff00d62a282278fb8719f4fab3a586a0a2c0fffbea::coin::T","u8","0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12::curves::Uncorrelated","u8"],"arguments":[9,"0",true,3,"1",true,8,"4",false,"10000000","694840"]}`

	chain := NewChainWithRestUrl(mainnetRestUrl)
	account, err := NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)

	txn, err := chain.GenerateTransactionJson(account.PublicKeyHex(), payloadJson)
	require.Nil(t, err)
	signedTxn, err := chain.SignTransactionJson(account, txn.Value)
	require.Nil(t, err)
	t.Log(signedTxn)
}
