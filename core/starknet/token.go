package starknet

import (
	"context"
	"errors"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
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
	felt         *felt.Felt
}

func NewToken(chain *Chain, tokenAddress string) (*Token, error) {
	felt, err := utils.HexToFelt(tokenAddress)
	if err != nil {
		return nil, err
	}
	return &Token{
		chain:        chain,
		TokenAddress: tokenAddress,
		felt:         felt,
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

	nameFelt, err := t.callContract(ctx, "name")
	if err != nil {
		return
	}
	nameBytes := utils.FeltToBigInt(nameFelt).Bytes()

	symbolFelt, err := t.callContract(ctx, "symbol")
	if err != nil {
		return
	}
	symbolBytes := utils.FeltToBigInt(symbolFelt).Bytes()

	decimalFelt, err := t.callContract(ctx, "decimals")
	if err != nil {
		return
	}
	decimal := utils.FeltToBigInt(decimalFelt).Int64()

	return &base.TokenInfo{
		Name:    string(nameBytes),
		Symbol:  string(symbolBytes),
		Decimal: int16(decimal),
	}, nil
}

func (t *Token) callContract(ctx context.Context, funcName string) (*felt.Felt, error) {
	res, err := t.chain.rpc.Call(ctx, rpc.FunctionCall{
		ContractAddress:    t.felt,
		EntryPointSelector: utils.GetSelectorFromNameFelt(funcName),
	}, latestBlockId)
	if err != nil || len(res) <= 0 {
		return nil, err
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
