package sui

import (
	"context"
	"math/big"
	"strconv"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

// "0x2::sui::SUI"
const SUI_COIN_TYPE = types.SUI_COIN_TYPE

// = 256-1
const MAX_INPUT_COUNT_MERGE = types.MAX_INPUT_COUNT_MERGE

// = 512-1
const MAX_INPUT_COUNT_STAKE = types.MAX_INPUT_COUNT_STAKE

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

func (t *Token) CoinType() string {
	return t.rType.ShortString()
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
	txn, err := t.BuildTransfer(account.Address(), receiverAddress, amount)
	return txn.(*Transaction), err
}

func (t *Token) EstimateFees(account *Account, receiverAddress, amount string) (f *base.OptionalString, err error) {
	txn, err := t.BuildTransferTransaction(account, receiverAddress, amount)
	if err != nil {
		return
	}
	gasString := strconv.FormatInt(txn.EstimateGasFee, 10)
	return &base.OptionalString{Value: gasString}, nil
}

func (t *Token) BuildTransfer(sender, receiver, amount string) (txn base.Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := types.NewAddressFromHex(sender)
	if err != nil {
		return
	}
	recipient, err := types.NewAddressFromHex(receiver)
	if err != nil {
		return
	}
	amountInt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return
	}
	cli, err := t.chain.Client()
	if err != nil {
		return
	}

	coinType := t.CoinType()
	coins, err := cli.GetCoins(context.Background(), *signer, &coinType, nil, MAX_INPUT_COUNT_MERGE)
	if err != nil {
		return
	}
	targetAmount := amountInt
	if t.IsSUI() {
		targetAmount = amountInt + MaxGasForTransfer // We will use PaySui if coin is SUI, the amount need plus gas
	}
	pickedCoin, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(targetAmount), MAX_INPUT_COUNT_MERGE, 0)
	if err != nil {
		return
	}

	maxGasBudget := maxGasBudget(pickedCoin, MaxGasForTransfer)
	return t.chain.EstimateTransactionFeeAndRebuildTransaction(maxGasBudget, func(gasBudget uint64) (*Transaction, error) {
		gasInt := types.NewSafeSuiBigInt(gasBudget)
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
				gasInt)
		} else {
			txnBytes, err = cli.Pay(context.Background(), *signer,
				pickedCoin.CoinIds(),
				[]types.Address{*recipient},
				[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amountInt)},
				nil, gasInt)
		}
		if err != nil {
			return nil, err
		}

		return &Transaction{Txn: *txnBytes}, nil
	})
}

func (t *Token) CanTransferAll() bool {
	return true
}

// Before invoking this method, it is best to check `CanTransferAll()`
func (t *Token) BuildTransferAll(sender, receiver string) (txn base.Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := types.NewAddressFromHex(sender)
	if err != nil {
		return
	}
	recipient, err := types.NewAddressFromHex(receiver)
	if err != nil {
		return
	}
	cli, err := t.chain.Client()
	if err != nil {
		return
	}

	coinType := t.CoinType()
	coins, err := cli.GetCoins(context.Background(), *signer, &coinType, nil, MAX_INPUT_COUNT_MERGE)
	if err != nil {
		return
	}
	if len(coins.Data) <= 0 {
		return nil, ErrNoCoinsFound
	}
	if coins.HasNextPage {
		return nil, ErrNeedMergeCoin
	}
	totalAmount := big.NewInt(0)
	coinIds := make([]types.ObjectId, len(coins.Data))
	for idx, coin := range coins.Data {
		coinIds[idx] = coin.CoinObjectId
		totalAmount.Add(totalAmount, big.NewInt(0).SetUint64(coin.Balance.Uint64()))
	}

	return t.chain.EstimateTransactionFeeAndRebuildTransaction(MaxGasForTransfer, func(gasBudget uint64) (*Transaction, error) {
		gasInt := types.NewSafeSuiBigInt(gasBudget)
		var txnBytes *types.TransactionBytes
		if t.IsSUI() {
			txnBytes, err = cli.PayAllSui(context.Background(), *signer,
				*recipient,
				coinIds,
				gasInt)
		} else {
			txnBytes, err = cli.Pay(context.Background(), *signer,
				coinIds,
				[]types.Address{*recipient},
				[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(totalAmount.Uint64())},
				nil, gasInt)
		}
		if err != nil {
			return nil, err
		}

		return &Transaction{Txn: *txnBytes}, nil
	})
}
