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
		Nonce:       "2132",
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
