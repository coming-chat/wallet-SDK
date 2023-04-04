package sui

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/coming-chat/go-sui/sui_types"
	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const (
	SuiName    = "Sui"
	SuiSymbol  = "Sui"
	SuiDecimal = 0
)

type Token struct {
	chain *Chain

	rType types.ResourceType
}

func NewTokenMain(chain *Chain) *Token {
	token, _ := NewToken(chain, "0x2::sui::SUI")
	return token
}

// @param tag format `address::module_name::name`, e.g. "0x2::sui::SUI"
func NewToken(chain *Chain, tag string) (*Token, error) {
	token, err := types.NewResourceType(tag)
	if err != nil {
		return nil, err
	}
	return &Token{chain, *token}, nil
}

func (t *Token) coinType() string {
	return fmt.Sprintf("0x2::coin::Coin<%v>", t.rType.ShortString())
}

// MARK - Implement the protocol Token

func (t *Token) Chain() base.Chain {
	return t.chain
}

func (t *Token) TokenInfo() (*base.TokenInfo, error) {
	metadata, err := t.getTokenMetadata(t.rType.ShortString())
	if err != nil {
		return nil, err
	}
	return &base.TokenInfo{
		Name:    metadata.Name,
		Symbol:  metadata.Symbol,
		Decimal: int16(metadata.Decimals),
	}, nil
}

func (t *Token) BalanceOfAddress(address string) (b *base.Balance, err error) {
	coins, err := t.getCoins(address)
	if err != nil {
		return nil, err
	}
	total := coins.TotalBalance().String()
	return &base.Balance{
		Total:  total,
		Usable: total,
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

func (t *Token) BuildTransferTxWithAccount(account *Account, receiverAddress, amount string) (s *base.OptionalString, err error) {
	txn, err := t.BuildTransferTransaction(account, receiverAddress, amount)
	if err != nil {
		return
	}
	return txn.SignWithAccount(account)
}

func (t *Token) BuildTransferTransaction(account *Account, receiverAddress, amount string) (s *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	recipient, err := types.NewAddressFromHex(receiverAddress)
	if err != nil {
		return
	}
	amountInt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return
	}

	coins, err := t.getCoins(account.Address())
	if err != nil {
		return nil, errors.New("Failed to get coins information.")
	}
	pickedCoin, err := pickupTransferCoin(coins, amount)
	if err != nil {
		return
	}

	cli, err := t.chain.Client()
	if err != nil {
		return
	}

	signer, _ := types.NewAddressFromHex(account.Address())
	if len(pickedCoin.Coins) >= 2 {
		// firstly, we should merge all coin's balance to firstCoin
		txn, err2 := cli.PayAllSui(context.Background(), *signer, *signer, pickedCoin.CoinIds(), pickedCoin.EstimateMergeGas())
		if err != nil {
			return nil, err2
		}
		signature, err2 := account.account.SignSecureWithoutEncode(txn.TxBytes, sui_types.DefaultIntent())
		if err != nil {
			return nil, err
		}
		response, err2 := cli.ExecuteTransactionBlock(context.Background(), txn.TxBytes, []any{signature}, &types.SuiTransactionBlockResponseOptions{ShowEffects: true}, types.TxnRequestTypeWaitForLocalExecution)
		if err2 != nil {
			return nil, err2
		}
		effects := response.Effects
		if *response.ConfirmedLocalExecution == false {
			return nil, fmt.Errorf("Merge coins failed.")
		}
		if !effects.IsSuccess() {
			return nil, fmt.Errorf(`Merge coins failed: %v`, effects.Status.Error)
		}
	}

	// send sui coin
	firstCoin := pickedCoin.Coins[0]
	txnBytes, err := cli.TransferSui(context.Background(), *signer, *recipient, firstCoin.CoinObjectId, amountInt, MaxGasForTransfer)
	if err != nil {
		return
	}
	return &Transaction{
		Txn:          *txnBytes,
		MaxGasBudget: MaxGasForTransfer,
	}, nil
}

func (t *Token) EstimateFees(account *Account, receiverAddress, amount string) (f *base.OptionalString, err error) {
	txn, err := t.BuildTransferTransaction(account, receiverAddress, amount)
	if err != nil {
		return
	}
	return t.chain.EstimateGasFee(txn)
}
