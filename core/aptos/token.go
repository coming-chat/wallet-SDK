package aptos

import (
	"encoding/json"
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
	mainTokenTag = "0x1::aptos_coin::AptosCoin"
)

type Token struct {
	chain *Chain

	token txbuilder.TypeTagStruct
}

func NewMainToken(chain *Chain) *Token {
	token, _ := NewToken(chain, mainTokenTag)
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
	contractAddress := t.token.Address.ToShortString()
	tag := "0x1::coin::CoinInfo<" + t.token.ShortFunctionName() + ">"
	client, err := t.chain.client()
	if err != nil {
		return nil, err
	}
	res, err := client.GetAccountResource(contractAddress, tag, 0)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(res.Data)
	if err != nil {
		return nil, err
	}
	info := struct {
		Decimals int16  `json:"decimals"`
		Name     string `json:"name"`
		Symbol   string `json:"symbol"`
	}{}
	err = json.Unmarshal(jsonData, &info)
	if err != nil {
		return nil, err
	}
	return &base.TokenInfo{
		Name:    info.Name,
		Symbol:  info.Symbol,
		Decimal: info.Decimals,
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
	f = &base.OptionalString{Value: "200000"}

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
	toAddr, err := txbuilder.NewAccountAddressFromHex(receiverAddress)
	if err != nil {
		return
	}
	amountInt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return
	}
	amountBytes := txbuilder.BCSSerializeBasicValue(amountInt)

	payloadBuilder := func(moduleName string, args []txbuilder.TypeTag) (txbuilder.TransactionPayload, error) {
		module, err := txbuilder.NewModuleIdFromString(moduleName)
		if err != nil {
			return nil, err
		}
		return txbuilder.TransactionPayloadEntryFunction{
			ModuleName:   *module,
			FunctionName: "transfer",
			TyArgs:       args,
			Args: [][]byte{
				toAddr[:], amountBytes,
			},
		}, nil
	}

	if t.token.ShortFunctionName() == mainTokenTag {
		hasRegisted, e := t.HasRegisted(receiverAddress)
		if e != nil || !hasRegisted.Value {
			// call "0x1::aptos_account::transfer"
			return payloadBuilder("0x1::aptos_account", []txbuilder.TypeTag{})
		}
	}
	// call "0x1::coin::transfer"
	return payloadBuilder("0x1::coin", []txbuilder.TypeTag{t.token})
}

func (t *Token) HasRegisted(ownerAddress string) (*base.OptionalBool, error) {
	tag := "0x1::coin::CoinStore<" + t.token.ShortFunctionName() + ">"
	client, err := t.chain.client()
	if err != nil {
		return nil, err
	}
	registed, err := client.IsAccountHasResource(ownerAddress, tag, 0)
	if err != nil {
		return nil, err
	}
	return &base.OptionalBool{Value: registed}, nil
}

func (t *Token) EnsureOwnerRegistedToken(owner *Account) (*base.OptionalString, error) {
	registed, err := t.HasRegisted(owner.Address())
	if err != nil {
		return nil, err
	}
	if registed.Value {
		return &base.OptionalString{}, nil
	}
	return t.RegisterTokenForOwner(owner)
}

// @return transaction hash if register token succeed.
func (t *Token) RegisterTokenForOwner(owner *Account) (*base.OptionalString, error) {
	moduleName, err := txbuilder.NewModuleIdFromString("0x1::managed_coin")
	if err != nil {
		return nil, err
	}
	payload := txbuilder.TransactionPayloadEntryFunction{
		ModuleName:   *moduleName,
		FunctionName: "register",
		TyArgs:       []txbuilder.TypeTag{t.token},
	}
	transaction, err := t.chain.createTransactionFromPayloadBCS(owner, payload)
	if err != nil {
		return nil, err
	}
	signedTx, err := txbuilder.GenerateBCSTransaction(owner.account, transaction)
	if err != nil {
		return nil, err
	}
	txString := types.HexEncodeToString(signedTx)
	hash, err := t.chain.SendRawTransaction(txString)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: hash}, nil
}
