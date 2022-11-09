package aptos

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	hexType "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/go-aptos/aptostypes"
	txbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type SignMessagePayload struct {
	// Should we include the address of the account in the message
	Address bool `json:"address,omitempty"`
	// Should we include the domain of the dApp
	Application bool `json:"application,omitempty"`
	// Should we include the current chain id the wallet is connected to
	ChainId bool `json:"chainId,omitempty"`
	// The message to be signed and displayed to the user
	Message string `json:"message"`
	// A nonce the dApp should generate
	Nonce int64 `json:"nonce"`
}

type SignMessageResponse struct {
	Address     string `json:"address,omitempty"`
	Application string `json:"application,omitempty"`
	ChainId     int64  `json:"chainId,omitempty"`
	Message     string `json:"message"` // The message passed in by the user
	Nonce       int64  `json:"nonce"`
	Prefix      string `json:"prefix"`           // Should always be APTOS
	FullMessage string `json:"fullMessage"`      // The message that was generated to sign
	Signature   string `json:"signature"`        // The signed full message
	Bitmap      []byte `json:"bitmap,omitempty"` // a 4-byte (32 bits) bit-vector of length N
}

func (j *SignMessagePayload) JsonString() string {
	bytes, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (j *SignMessageResponse) JsonString() string {
	bytes, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (c *Chain) GenerateTransaction(senderPublicKey string, payload aptostypes.Payload) (txn *aptostypes.Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	payload.Type = aptostypes.EntryFunctionPayload

	sender, err := EncodePublicKeyToAddress(senderPublicKey)
	if err != nil {
		return
	}
	client, err := c.client()
	if err != nil {
		return
	}
	remoteBuilder, err := txbuilder.NewTransactionBuilderRemoteABIWithFunc(payload.Function, client)
	if err != nil {
		return
	}
	bcsPayload, err := remoteBuilder.BuildTransactionPayload(payload.Function, payload.TypeArguments, payload.Arguments)
	if err != nil {
		return
	}
	accountData, err := client.GetAccount(sender)
	if err != nil {
		return
	}
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		return
	}
	gasPrice, err := client.EstimateGasPrice()
	if err != nil {
		return
	}

	senderAddress, _ := txbuilder.NewAccountAddressFromHex(sender)
	rawTxn := txbuilder.RawTransaction{
		Sender:                  *senderAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            MaxGasAmount,
		GasUnitPrice:            gasPrice,
		Payload:                 bcsPayload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // 10 minutes timeout
		ChainId:                 uint8(ledgerInfo.ChainId),
	}
	edPublicKey, _ := hexType.HexDecodeString(senderPublicKey)
	signedTxn, err := txbuilder.GenerateBCSSimulation(edPublicKey, &rawTxn)
	if err != nil {
		return
	}
	txns, err := client.SimulateSignedBCSTransaction(signedTxn)
	if err != nil {
		return
	}
	if len(txns) <= 0 {
		return nil, errors.New("Generate transaction failed.")
	}
	txn = txns[0]
	if !txn.Success {
		return nil, errors.New(txn.VmStatus)
	}
	gasUsed := txn.GasUsed/10*15 + 14 // ceil(gasUsed * 1.5)
	txn.MaxGasAmount = gasUsed

	txn.Hash = ""
	txn.Signature = nil // clean simulate signature.
	return txn, nil
}

func (c *Chain) SignTransaction(account base.Account, transaction aptostypes.Transaction) (txn *aptostypes.Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	client, err := c.client()
	if err != nil {
		return
	}
	signingMessage, err := client.CreateTransactionSigningMessage(&transaction)
	if err != nil {
		return
	}

	signatureData, err := account.Sign(signingMessage, "")
	signatureHex := "0x" + hex.EncodeToString(signatureData)
	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: account.PublicKeyHex(),
		Signature: signatureHex,
	}
	return &transaction, nil
}

func (c *Chain) SubmitTransaction(signedTxn aptostypes.Transaction) (txhash string, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if signedTxn.Signature == nil {
		return "", errors.New("Submit failed: Your transaction has not been signed.")
	}
	client, err := c.client()
	if err != nil {
		return
	}
	subTxn, err := client.SubmitTransaction(&signedTxn)
	if err != nil {
		return
	}
	return subTxn.Hash, nil
}

func (c *Chain) SignMessage(account base.Account, payload *SignMessagePayload) (*SignMessageResponse, error) {
	chainId := 0
	client, err := c.client()
	if err != nil {
		err = nil
	} else {
		chainId = client.ChainId()
	}

	resp := &SignMessageResponse{
		Address:     account.Address(),
		Application: "",
		ChainId:     int64(chainId),
		Message:     payload.Message,
		Nonce:       payload.Nonce,
		Prefix:      "APTOS",
		Bitmap:      nil,
	}
	msg := resp.Prefix
	if payload != nil {
		if payload.Address {
			msg += "\naddress: " + resp.Address
		}
		if payload.Application {
			msg += "\napplication: " + resp.Application
		}
		if payload.ChainId {
			msg += "\nchainId: " + fmt.Sprintf("%v", resp.ChainId)
		}
		if payload.Message != "" {
			msg += "\nmessage: " + resp.Message
		}
		msg += "\nnonce: " + fmt.Sprintf("%v", resp.Nonce)
	}

	bytes := []byte(resp.FullMessage)
	signature, err := account.Sign(bytes, "")
	if err != nil {
		return nil, err
	}
	resp.FullMessage = msg
	resp.Signature = "0x" + hex.EncodeToString(signature)
	return resp, nil
}

// function for mobile client

func (c *Chain) GenerateTransactionJson(senderPublicKey string, payload string) (*base.OptionalString, error) {
	bytes := []byte(payload)
	var payloadObj aptostypes.Payload
	err := json.Unmarshal(bytes, &payloadObj)
	if err != nil {
		return nil, err
	}
	res, err := c.GenerateTransaction(senderPublicKey, payloadObj)
	if err != nil {
		return nil, err
	}
	bytes, err = json.Marshal(res)
	if err != nil {
		return nil, err
	}
	resJson := string(bytes)
	return &base.OptionalString{Value: resJson}, nil
}

func (c *Chain) SignTransactionJson(account base.Account, transaction string) (*base.OptionalString, error) {
	bytes := []byte(transaction)
	var transactionObj aptostypes.Transaction
	err := json.Unmarshal(bytes, &transactionObj)
	if err != nil {
		return nil, err
	}
	res, err := c.SignTransaction(account, transactionObj)
	if err != nil {
		return nil, err
	}
	bytes, err = json.Marshal(res)
	if err != nil {
		return nil, err
	}
	resJson := string(bytes)
	return &base.OptionalString{Value: resJson}, nil
}

func (c *Chain) SubmitTransactionJson(signedTxn string) (*base.OptionalString, error) {
	bytes := []byte(signedTxn)
	var transactionObj aptostypes.Transaction
	err := json.Unmarshal(bytes, &transactionObj)
	if err != nil {
		return nil, err
	}
	res, err := c.SubmitTransaction(transactionObj)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: res}, nil
}

func (c *Chain) SignMessageJson(account base.Account, payload string) (*SignMessageResponse, error) {
	bytes := []byte(payload)
	var payloadObj SignMessagePayload
	err := json.Unmarshal(bytes, &payloadObj)
	if err != nil {
		return nil, err
	}
	return c.SignMessage(account, &payloadObj)
}
