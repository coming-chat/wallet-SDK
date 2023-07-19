package starcoin

import (
	"context"
	"fmt"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/starcoinorg/starcoin-go/client"
	"github.com/starcoinorg/starcoin-go/types"
)

type Token struct {
	chain    *Chain
	tokenTag types.StructTag
}

func NewMainToken(chain *Chain) *Token {
	token, _ := newTokenWithTag(chain, client.GAS_TOKEN_CODE)
	return token
}

// @param tag format `address::module_name::name`, e.g. "0x1::STC::STC"
func newTokenWithTag(chain *Chain, tag string) (*Token, error) {
	tokenTag, err := NewStructTag(tag)
	if err != nil {
		return nil, err
	}
	return &Token{
		chain:    chain,
		tokenTag: *tokenTag,
	}, nil
}

// MARK - Implement the protocol Token

func (t *Token) Chain() base.Chain {
	return t.chain
}

func (t *Token) TokenInfo() (*base.TokenInfo, error) {
	return &base.TokenInfo{
		Name:    "STC",
		Symbol:  "STC",
		Decimal: 9,
	}, nil
}

func (t *Token) BalanceOfAddress(address string) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	b = base.EmptyBalance()

	ls, err := t.chain.client.ListResource(context.Background(), address)
	if err != nil {
		return
	}
	balances, err := ls.GetBalances()
	if err != nil {
		return
	}

	identifier := fmt.Sprintf("0x00000000000000000000000000000001::Account::Balance<%v>", StructTagToString(t.tokenTag))
	for key, balance := range balances {
		if key == identifier {
			b.Total = balance.String()
			b.Usable = balance.String()
			return
		}
	}
	return
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

func (t *Token) BuildTransferTxWithAccount(account *Account, receiverAddress, amount string) (s *base.OptionalString, err error) {
	txn, err := t.BuildTransfer(account.Address(), receiverAddress, amount)
	if err != nil {
		return
	}
	signedTxn, err := txn.SignedTransactionWithAccount(account)
	if err != nil {
		return
	}
	return signedTxn.HexString()
}

func (t *Token) EstimateFees(account *Account, receiverAddress, amount string) (f *base.OptionalString, err error) {
	txn, err := t.BuildTransfer(account.Address(), receiverAddress, amount)
	if err != nil {
		return
	}
	return t.chain.EstimateTransactionFeeUsePublicKey(txn, account.PublicKeyHex())
}

func (t *Token) BuildTransferPayload(receiverAddress, amount string) (p types.TransactionPayload, err error) {
	receiver, err := NewAccountAddressFromHex(receiverAddress)
	if err != nil {
		return
	}
	amountInt, err := NewU128FromString(amount)
	if err != nil {
		return nil, fmt.Errorf("Invalid transfer amount %v", amount)
	}
	p = client.Encode_peer_to_peer_v2_script_function(&types.TypeTag__Struct{Value: t.tokenTag}, *receiver, *amountInt)
	return
}

func (t *Token) BuildTransfer(sender, receiver, amount string) (txn base.Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	senderAddress, err := NewAccountAddressFromHex(sender)
	if err != nil {
		return
	}
	payload, err := t.BuildTransferPayload(receiver, amount)
	if err != nil {
		return
	}

	ctx := context.Background()
	cli := t.chain.client
	price, err := cli.GetGasUnitPrice(ctx)
	if err != nil {
		return
	}
	state, err := t.chain.GetState(ctx, sender)
	if err != nil {
		return
	}
	rawTxn, err := cli.BuildRawUserTransaction(ctx, *senderAddress, payload, price, MaxGasAmount, state.SequenceNumber)
	if err != nil {
		return
	}
	return &Transaction{
		Txn: rawTxn,
	}, nil
}
func (t *Token) CanTransferAll() bool {
	return false
}
func (t *Token) BuildTransferAll(sender, receiver string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}
