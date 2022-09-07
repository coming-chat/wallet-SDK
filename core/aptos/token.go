package aptos

import (
	"strconv"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	txbuilder "github.com/coming-chat/go-aptos/transaction_builder"
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
	payload, err := t.buildTransferPayload(receiverAddress, amount)
	if err != nil {
		return nil, err
	}
	transaction, err := t.chain.createTransactionFromPayloadBCS(account, payload)
	if err != nil {
		return nil, err
	}
	signedTx, err := txbuilder.GenerateBCSTransaction(account.account, transaction)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: types.HexEncodeToString(signedTx)}, nil
}

func (t *Token) EstimateFees(account *Account, receiverAddress, amount string) (f *base.OptionalString, err error) {
	f = &base.OptionalString{Value: "2000"}

	payload, err := t.buildTransferPayload(receiverAddress, amount)
	if err != nil {
		return
	}
	transaction, err := t.chain.createTransactionFromPayloadBCS(account, payload)
	if err != nil {
		return
	}
	gasFee := transaction.MaxGasAmount * transaction.GasUnitPrice
	gasString := strconv.FormatUint(gasFee, 10)
	return &base.OptionalString{Value: gasString}, nil
}

func (t *Token) buildTransferPayload(receiverAddress, amount string) (p txbuilder.TransactionPayload, err error) {
	moduleName, err := txbuilder.NewModuleIdFromString("0x1::coin")
	if err != nil {
		return
	}
	token, err := txbuilder.NewTypeTagStructFromString("0x1::aptos_coin::AptosCoin")
	if err != nil {
		return
	}
	toAddr, err := txbuilder.NewAccountAddressFromHex(receiverAddress)
	if err != nil {
		return
	}
	amountInt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return
	}
	amountBytes := txbuilder.BCSSerializeBasicValue(amountInt)
	return txbuilder.TransactionPayloadEntryFunction{
		ModuleName:   *moduleName,
		FunctionName: "transfer",
		TyArgs:       []txbuilder.TypeTag{*token},
		Args: [][]byte{
			toAddr[:], amountBytes,
		},
	}, nil
}
