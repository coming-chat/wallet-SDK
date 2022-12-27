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
	SuiName     = "Sui"
	SuiSymbol   = "Sui"
	SuiDecimal  = 0
	SuiCoinType = "0x2::coin::Coin<0x2::sui::SUI>"
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
	var pickedCoin *PickedCoins
	if t.coinType() != SuiCoinType {
		pickedCoin, err = pickupTransferCoin(coins, amount, 0)
		suiToken := NewTokenMain(t.chain)
		suiCoins, err := suiToken.getCoins(account.Address())
		if err != nil {
			return nil, err
		}
		pickedGasCoin, err := pickupTransferCoin(suiCoins, "0", MaxGasForTransfer)
		if err != nil {
			return nil, err
		}
		pickedCoin.Coins = append(pickedCoin.Coins, pickedGasCoin.Coins...)
	} else {
		pickedCoin, err = pickupTransferCoin(coins, amount, MaxGasForPay)
	}

	if err != nil {
		return
	}

	cli, err := t.chain.client()
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
		signedTxn := txn.SignWith(account.account.PrivateKey)
		response, err2 := cli.ExecuteTransaction(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
		if err2 != nil {
			return nil, err2
		}
		cert := response.EffectsCert
		if cert == nil || cert.ConfirmedLocalExecution == false {
			return nil, fmt.Errorf("Merge coins failed.")
		}
		status := cert.Effects.Effects.Status
		if status.Status != types.TransactionStatusSuccess {
			return nil, fmt.Errorf(`Merge coins failed: %v`, status.Error)
		}
	}

	// send sui coin
	firstCoin := pickedCoin.Coins[0]
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
