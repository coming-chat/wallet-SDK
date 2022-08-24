package aptos

import (
	"encoding/json"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/go-aptos/aptostypes"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const (
	AptosName    = "Aptos"
	AptosSymbol  = "Aptos"
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
	payload := t.buildTransferPayload(receiverAddress, amount)
	transaction, err := t.chain.createTransactionFromPayload(account, payload)
	if err != nil {
		return nil, err
	}
	signedTransaction, err := t.chain.signTransaction(account, transaction)
	if err != nil {
		return nil, err
	}
	signedTransactionData, err := json.Marshal(signedTransaction)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: types.HexEncodeToString(signedTransactionData)}, nil
}

func (t *Token) EstimateFees(account *Account, receiverAddress, amount string) (f *base.OptionalString, err error) {
	f = &base.OptionalString{Value: "2000"}

	payload := t.buildTransferPayload(receiverAddress, amount)
	transaction, err := t.chain.createTransactionFromPayload(account, payload)
	if err != nil {
		return
	}
	return t.chain.EstimateGasFee(account, transaction)
}

func (t *Token) buildTransferPayload(receiverAddress, amount string) *aptostypes.Payload {
	return &aptostypes.Payload{
		Type:          aptostypes.EntryFunctionPayload,
		Function:      "0x1::account::transfer",
		TypeArguments: []string{},
		Arguments: []interface{}{
			receiverAddress, amount,
		},
	}
}
