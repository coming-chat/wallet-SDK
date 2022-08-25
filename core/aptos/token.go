package aptos

import (
	"encoding/json"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const (
	AptosName    = "APT"
	AptosSymbol  = "APT"
	AptosDecimal = 0
)

type Token struct {
	chain *Chain
}

func NewToken(chain *Chain) *Token {
	return &Token{chain}
}

// MARK - Implement the protocol Token

func (t *Token) Chain() base.Chain {
	return t.chain
}

func (t *Token) TokenInfo() (*base.TokenInfo, error) {
	return &base.TokenInfo{
		Name:    AptosName,
		Symbol:  AptosSymbol,
		Decimal: AptosDecimal,
	}, nil
}

func (t *Token) BalanceOfAddress(address string) (*base.Balance, error) {
	return t.chain.BalanceOfAddress(address)
}

func (t *Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return t.chain.BalanceOfPublicKey(publicKey)
}

func (t *Token) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.chain.BalanceOfAddress(account.Address())
}

// MARK - token

func (t *Token) BuildTransferTx(privateKey, receiverAddress, amount string) (*base.OptionalString, error) {
	account, err := AccountWithPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return t.BuildTransferTxWithAccount(account, receiverAddress, amount)
}

func (t *Token) BuildTransferTxWithAccount(account *Account, receiverAddress, amount string) (*base.OptionalString, error) {
	client, err := t.chain.client()
	if err != nil {
		return nil, err
	}

	transaction, err := t.buildSigningTransaction(client, account, receiverAddress, amount)
	if err != nil {
		return nil, err
	}

	signingMessage, err := client.CreateTransactionSigningMessage(transaction)
	if err != nil {
		return nil, err
	}
	signatureData, _ := account.Sign(signingMessage, "")
	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: account.PublicKeyHex(),
		Signature: types.HexEncodeToString(signatureData),
	}

	signedTransactionData, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	return &base.OptionalString{Value: types.HexEncodeToString(signedTransactionData)}, nil
}

func (t *Token) EstimateFees(account *Account, receiverAddress, amount string) (f *base.OptionalString, err error) {
	f = &base.OptionalString{Value: "2000"}

	client, err := t.chain.client()
	if err != nil {
		return
	}

	transaction, err := t.buildSigningTransaction(client, account, receiverAddress, amount)
	if err != nil {
		return
	}
	return t.chain.EstimateGasFee(account, transaction)
}

func (t *Token) buildSigningTransaction(client *aptosclient.RestClient, account *Account, receiverAddress, amount string) (tx *aptostypes.Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	fromAddress := account.Address()
	accountData, err := client.GetAccount(fromAddress)
	if err != nil {
		return
	}
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		return
	}
	payload := &aptostypes.Payload{
		Type:          aptostypes.EntryFunctionPayload,
		Function:      "0x1::account::transfer",
		TypeArguments: []string{},
		Arguments: []interface{}{
			receiverAddress, amount,
		},
	}
	transaction := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            MaxGasAmount,
		GasUnitPrice:            GasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // timeout 10 mins
	}

	return transaction, nil
}
