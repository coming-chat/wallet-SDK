package aptos

import (
	"encoding/json"
	"strconv"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/go-aptos/aptostypes"
	"github.com/coming-chat/wallet-SDK/core/base"
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
		Name:    "Aptos",
		Symbol:  "Aptos",
		Decimal: 0,
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
	fromAddress := account.Address()

	client, err := t.chain.client()
	if err != nil {
		return nil, err
	}

	accountData, err := client.GetAccount(fromAddress)
	if err != nil {
		return nil, err
	}
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		return nil, err
	}

	payload := &aptostypes.Payload{
		Type:          "script_function_payload",
		Function:      "0x1::coin::transfer",
		TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
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

func (t *Token) EstimateFees(receiverAddress, amount string) (*base.OptionalString, error) {
	gas := GasPrice * MaxGasAmount
	gasString := strconv.FormatInt(int64(gas), 10)
	return &base.OptionalString{Value: gasString}, nil
}
