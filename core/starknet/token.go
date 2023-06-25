package starknet

import (
	"context"
	"errors"
	"math/big"

	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/dontpanicdao/caigo/types"
)

const (
	ETHTokenAddress = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
)

type Token struct {
	chain        *Chain
	TokenAddress string
}

func NewToken(chain *Chain, tokenAddress string) (*Token, error) {
	_, err := hexTypes.HexDecodeString(tokenAddress)
	if err != nil {
		return nil, err
	}
	return &Token{
		chain:        chain,
		TokenAddress: tokenAddress,
	}, nil
}

// MARK - Implement the protocol Token

func (t *Token) Chain() base.Chain {
	return t.chain
}

// Warning: Main token does not support
func (t *Token) TokenInfo() (info *base.TokenInfo, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if t == nil || t.chain == nil || t.chain.gw == nil {
		return nil, errors.New("nil params")
	}
	ctx := context.Background()

	nameHex, err := t.callContract(ctx, "name")
	if err != nil {
		return
	}
	nameBytes, _ := types.HexToBytes(nameHex)

	symbolHex, err := t.callContract(ctx, "symbol")
	if err != nil {
		return
	}
	symbolBytes, _ := types.HexToBytes(symbolHex)

	decimalHex, err := t.callContract(ctx, "decimals")
	if err != nil {
		return
	}
	decimal := types.HexToBN(decimalHex).Int64()

	return &base.TokenInfo{
		Name:    string(nameBytes),
		Symbol:  string(symbolBytes),
		Decimal: int16(decimal),
	}, nil
}

func (t *Token) callContract(ctx context.Context, funcName string) (string, error) {
	res, err := t.chain.gw.Call(ctx, types.FunctionCall{
		ContractAddress:    types.HexToHash(t.TokenAddress),
		EntryPointSelector: funcName,
	}, "")
	if err != nil || len(res) <= 0 {
		return "", err
	}
	return res[0], nil
}

func (t *Token) BalanceOfAddress(address string) (*base.Balance, error) {
	return t.chain.BalanceOf(address, t.TokenAddress)
}

func (t *Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	address, err := encodePublicKeyToAddressArgentX(publicKey)
	if err != nil {
		return nil, err
	}
	return t.chain.BalanceOf(address, t.TokenAddress)
}

func (t *Token) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.chain.BalanceOf(account.Address(), t.TokenAddress)
}

func (t *Token) BuildTransfer(sender, receiver, amount string) (txn base.Transaction, err error) {
	if _, err = hexTypes.HexDecodeString(sender); err != nil {
		return nil, base.ErrInvalidAccountAddress
	}
	if _, err = hexTypes.HexDecodeString(receiver); err != nil {
		return nil, base.ErrInvalidAccountAddress
	}
	if _, ok := big.NewInt(0).SetString(amount, 10); !ok {
		return nil, base.ErrInvalidAmount
	}
	// Transaction that will be executed by the account contract.
	tx := []types.FunctionCall{
		{
			ContractAddress:    types.HexToHash(t.TokenAddress),
			EntryPointSelector: "transfer",
			Calldata: []string{
				// sender,
				receiver,
				amount, // amount to transfer
				"0",    // UInt256 additional parameter
			},
		},
	}
	return &Transaction{
		calls:   tx,
		details: types.ExecuteDetails{},
	}, nil
}

// Before invoking this method, it is best to check `CanTransferAll()`
func (t *Token) CanTransferAll() bool {
	return false
}
func (t *Token) BuildTransferAll(sender, receiver string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}
