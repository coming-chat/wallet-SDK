package sui

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const SUI_COIN_TYPE = "0x2::sui::SUI"

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
	token, _ := NewToken(chain, SUI_COIN_TYPE)
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

func (t *Token) IsSUI() bool {
	return t.rType.ShortString() == SUI_COIN_TYPE
}

// MARK - Implement the protocol Token

func (t *Token) Chain() base.Chain {
	return t.chain
}

func (t *Token) TokenInfo() (info *base.TokenInfo, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := t.chain.Client()
	if err != nil {
		return
	}
	metadata, err := cli.GetCoinMetadata(context.Background(), t.rType.ShortString())
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
	defer base.CatchPanicAndMapToBasicError(&err)

	owner, err := types.NewAddressFromHex(address)
	if err != nil {
		return nil, err
	}
	cli, err := t.chain.Client()
	if err != nil {
		return nil, err
	}
	balance, err := cli.GetBalance(context.Background(), *owner, t.rType.ShortString())
	if err != nil {
		return nil, err
	}
	balanceStr := balance.TotalBalance.String()
	return &base.Balance{
		Total:  balanceStr,
		Usable: balanceStr,
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

	coins, err := t.getCoins(account.Address(), 0)
	if err != nil {
		return nil, errors.New("Failed to get coins information.")
	}
	pickedCoin, err := pickupTransferCoin(coins, amountInt, t.IsSUI())
	if err != nil {
		return
	}

	cli, err := t.chain.Client()
	if err != nil {
		return
	}

	signer, _ := types.NewAddressFromHex(account.Address())
	gasBudget := types.NewSafeSuiBigInt[uint64](MaxGasForTransfer)
	var txnBytes *types.TransactionBytes
	// TODO: we can transfer object now, but we cannot parse it's to a coin transfer event.
	// if pickedCoin.CanUseTransferObject {
	// 	txnBytes, err = cli.TransferObject(context.Background(), *signer, *recipient,
	// 		pickedCoin.Coins[0].CoinObjectId,
	// 		nil, gasBudget)
	// } else {
	// }
	if t.IsSUI() {
		txnBytes, err = cli.PaySui(context.Background(), *signer,
			pickedCoin.CoinIds(),
			[]types.Address{*recipient},
			[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountInt)},
			gasBudget)
	} else {
		txnBytes, err = cli.Pay(context.Background(), *signer,
			pickedCoin.CoinIds(),
			[]types.Address{*recipient},
			[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountInt)},
			nil, gasBudget)
	}
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
