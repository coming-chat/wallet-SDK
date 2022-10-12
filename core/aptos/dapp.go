package aptos

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coming-chat/go-aptos/aptostypes"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type SignMessagePayload struct {
	// Should we include the address of the account in the message
	Address bool `json:"address"`
	// Should we include the domain of the dApp
	Application bool `json:"application"`
	// Should we include the current chain id the wallet is connected to
	ChainId bool `json:"chainId"`
	// The message to be signed and displayed to the user
	Message string `json:"message"`
	// A nonce the dApp should generate
	Nonce string `json:"nonce"`
}

type SignMessageResponse struct {
	Address     string `json:"address"`
	Application string `json:"application"`
	ChainId     int64  `json:"chainId"`
	Message     string `json:"message"` // The message passed in by the user
	Nonce       string `json:"nonce"`
	Prefix      string `json:"prefix"`      // Should always be APTOS
	FullMessage string `json:"fullMessage"` // The message that was generated to sign
	Signature   string `json:"signature"`   // The signed full message
	Bitmap      []byte `json:"bitmap"`      // a 4-byte (32 bits) bit-vector of length N
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

	txn = &aptostypes.Transaction{
		Sender:                  sender,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            MaxGasAmount,
		GasUnitPrice:            gasPrice,
		Payload:                 &payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // 10 minutes timeout
	}

	simTxn, err := client.SimulateTransaction(txn, senderPublicKey)
	if err != nil {
		return nil, err
	}
	if len(simTxn) > 0 {
		if !simTxn[0].Success {
			return nil, errors.New(simTxn[0].VmStatus)
		} else {
			txn.MaxGasAmount = simTxn[0].GasUsed
		}
	}

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
	resp := &SignMessageResponse{
		Address:     account.Address(),
		Application: "",
		ChainId:     0,
		Message:     payload.Message,
		Nonce:       payload.Nonce,
		Prefix:      "APTOS",
		Bitmap:      nil,
	}
	msg := fmt.Sprintf(`APTOS\naddress: %v\napplication: %v\nchainId: %v\nmessage: %v\nnonce: %v`,
		resp.Address, resp.Application, resp.ChainId, resp.Message, resp.Nonce)
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

func (c *Chain) GenerateTransactionJson(sender string, payload string) (*base.OptionalString, error) {
	bytes := []byte(payload)
	var payloadObj aptostypes.Payload
	err := json.Unmarshal(bytes, &payloadObj)
	if err != nil {
		return nil, err
	}
	res, err := c.GenerateTransaction(sender, payloadObj)
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
