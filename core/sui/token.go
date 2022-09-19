package sui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

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
	token, _ := newToken(chain, "0x2::sui::SUI")
	return token
}

// @param tag format `address::module_name::name`, e.g. "0x2::sui::SUI"
func newToken(chain *Chain, tag string) (*Token, error) {
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
	return &base.TokenInfo{
		Name:    SuiName,
		Symbol:  SuiSymbol,
		Decimal: SuiDecimal,
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

	cli, err := t.chain.client()
	if err != nil {
		return
	}

	signer, _ := types.NewAddressFromHex(account.Address())
	firstCoin := pickedCoin.Coins[0]
	if len(pickedCoin.Coins) >= 2 {
		var txn *types.TransactionBytes
		var response *types.TransactionResponse
		// firstly, we should merge all coin's balance to firstCoin
		for i := 1; i < len(pickedCoin.Coins); i++ {
			coin := pickedCoin.Coins[i]
			txn, err = cli.MergeCoins(context.Background(), *signer, firstCoin.Reference.ObjectId, coin.Reference.ObjectId, nil, MaxGasForMerge)
			if err != nil {
				return
			}
			signedTxn := txn.SignWith(account.account.PrivateKey)
			response, err = cli.ExecuteTransaction(context.Background(), *signedTxn)
			if err != nil {
				return
			}
			if response.Effects.Status.Status != types.TransactionStatusSuccess {
				return nil, fmt.Errorf(`Merge coin failed: %v`, response.Effects.Status.Error)
			}
		}
	}

	// send sui coin
	txn, err := cli.TransferSui(context.Background(), *signer, *recipient, firstCoin.Reference.ObjectId, amountInt, MaxGasForTransfer)
	if err != nil {
		return
	}
	signedTxn := txn.SignWith(account.account.PrivateKey)
	bytes, err := json.Marshal(signedTxn)
	if err != nil {
		return
	}
	txnString := types.Bytes(bytes).GetBase64Data().String()

	return &base.OptionalString{Value: txnString}, nil
}

func (t *Token) EstimateFees(account *Account, receiverAddress, amount string) (f *base.OptionalString, err error) {
	f = &base.OptionalString{Value: strconv.FormatInt(MaxGasBudget, 10)}
	return f, nil
}
