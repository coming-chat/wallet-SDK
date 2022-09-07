package aptos

import (
	"errors"
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

	token txbuilder.TypeTagStruct
}

func NewMainToken(chain *Chain) *Token {
	token, _ := NewToken(chain, "0x1::aptos_coin::AptosCoin")
	return token
}

// @param tag format `address::module_name::name`, e.g. "0x1::aptos_coin::AptosCoin"
func NewToken(chain *Chain, tag string) (*Token, error) {
	token, err := txbuilder.NewTypeTagStructFromString(tag)
	if err != nil {
		return nil, err
	}
	return &Token{chain, *token}, nil
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

func (t *Token) BalanceOfAddress(address string) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	client, err := t.chain.client()
	if err != nil {
		return
	}
	balance, err := client.BalanceOf(address, t.token.ShortFunctionName())
	if err != nil {
		return
	}

	return &base.Balance{
		Total:  balance.String(),
		Usable: balance.String(),
	}, nil
}

func (t *Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	address, err := EncodePublicKeyToAddress(publicKey)
	if err != nil {
		return nil, err
	}
	return t.BalanceOfAddress(address)
}

func (t *Token) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.BalanceOfAddress(account.Address())
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
	if t.token.Address.ToShortString() == "0x" {
		return nil, errors.New("Invalid token tag: " + t.token.ShortFunctionName())
	}
	moduleName, err := txbuilder.NewModuleIdFromString("0x1::coin")
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
		TyArgs:       []txbuilder.TypeTag{t.token},
		Args: [][]byte{
			toAddr[:], amountBytes,
		},
	}, nil
}

func (t *Token) EnsureOwnerRegistedToken(ownerAddress string, from *Account) error {
	tag := "0x1::coin::CoinStore<" + t.token.ShortFunctionName() + ">"
	client, err := t.chain.client()
	if err != nil {
		return err
	}
	registed, err := client.IsAccountHasResource(ownerAddress, tag, 0)
	if err != nil {
		return err
	}
	if registed {
		return nil
	}
	_, err = t.RegisterTokenForOwner(ownerAddress, from)
	return err
}

// @return transaction hash if register token succeed.
func (t *Token) RegisterTokenForOwner(ownerAddress string, from *Account) (string, error) {
	moduleName, err := txbuilder.NewModuleIdFromString("0x1::managed_coin")
	if err != nil {
		return "", err
	}
	payload := txbuilder.TransactionPayloadEntryFunction{
		ModuleName:   *moduleName,
		FunctionName: "register",
		TyArgs:       []txbuilder.TypeTag{t.token},
	}
	transaction, err := t.chain.createTransactionFromPayloadBCS(from, payload)
	if err != nil {
		return "", err
	}
	signedTx, err := txbuilder.GenerateBCSTransaction(from.account, transaction)
	if err != nil {
		return "", err
	}
	txString := types.HexEncodeToString(signedTx)
	return t.chain.SendRawTransaction(txString)
}
