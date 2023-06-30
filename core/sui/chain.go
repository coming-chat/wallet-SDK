package sui

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/fardream/go-bcs/bcs"
)

const (
	MaxGasBudget = 90000000
	MinGasBudget = 1100000
	// MaxGasForMerge = 10000000

	MaxGasForPay      = 10000000
	MaxGasForTransfer = 10000000

	TestNetFaucetUrl = client.TestNetFaucetUrl
	DevNetFaucetUrl  = client.DevNetFaucetUrl
)

type Chain struct {
	rpcClient *client.Client
	RpcUrl    string

	gasPrice uint64
}

func NewChainWithRpcUrl(rpcUrl string) *Chain {
	return &Chain{RpcUrl: rpcUrl}
}

func (c *Chain) Client() (*client.Client, error) {
	if c.rpcClient != nil {
		return c.rpcClient, nil
	}
	var err error
	c.rpcClient, err = client.Dial(c.RpcUrl)
	return c.rpcClient, err
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return NewTokenMain(c)
}

func (c *Chain) BalanceOfAddress(address string) (b *base.Balance, err error) {
	return c.MainToken().BalanceOfAddress(address)
}
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return c.MainToken().BalanceOfPublicKey(publicKey)
}
func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (hash string, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	bytes, err := lib.NewBase64Data(signedTx)
	if err != nil {
		return
	}
	signedTxn := SignedTransaction{}
	err = json.Unmarshal(bytes.Data(), &signedTxn)
	if err != nil {
		return
	}
	cli, err := c.Client()
	if err != nil {
		return
	}
	options := types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	}
	response, err := cli.ExecuteTransactionBlock(context.Background(), signedTxn.TxBytes.Data(), []any{signedTxn.Signature}, &options, types.TxnRequestTypeWaitForLocalExecution)
	if err != nil {
		return
	}
	hash = response.Digest.String()
	if !response.Effects.Data.IsSuccess() {
		return hash, errors.New(response.Effects.Data.V1.Status.Error)
	}
	return hash, nil
}

func (c *Chain) SendSignedTransaction(signedTxn base.SignedTransaction) (hash *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	digest, err := sui_types.NewDigest(hash)
	if err != nil {
		return
	}
	cli, err := c.Client()
	if err != nil {
		return
	}
	resp, err := cli.GetTransactionBlock(context.Background(), *digest, types.SuiTransactionBlockResponseOptions{
		ShowInput:   true,
		ShowEffects: true,
		ShowEvents:  true,
	})
	if err != nil {
		return nil, err
	}

	var notCoinTransferErr = errors.New("Invalid coin transfer transaction.")
	var firstRecipient string
	var total string
	if resp.Transaction == nil || resp.Transaction.Data.Data.V1 == nil {
		return nil, errors.New("failed to retrieve data")
	}
	var data = resp.Transaction.Data.Data.V1
	if data.Transaction.Data.ProgrammableTransaction == nil {
		return nil, notCoinTransferErr
	}
	dataBytes, err := json.Marshal(data.Transaction.Data.ProgrammableTransaction)
	if err != nil {
		return nil, err
	}
	var transaction struct {
		Inputs []struct {
			Type      string `json:"type"`
			ValueType string `json:"valueType"`
			Value     string `json:"value"`
		} `json:"inputs"`
	}
	err = json.Unmarshal(dataBytes, &transaction)
	if err != nil {
		return nil, err
	}
	if len(transaction.Inputs) == 2 {
		for _, input := range transaction.Inputs {
			if input.Type != "pure" {
				return nil, notCoinTransferErr
			}
			switch input.ValueType {
			case "address":
				firstRecipient = input.Value
			case "u64":
				total = input.Value
			default:
				return nil, notCoinTransferErr
			}
		}
	} else {
		firstRecipient = ""
		total = ""
	}

	var effects types.SuiTransactionBlockEffects
	if resp.Effects != nil {
		effects = resp.Effects.Data
	}
	detail = &base.TransactionDetail{
		HashString:      hash,
		FromAddress:     data.Sender.ShortString(),
		ToAddress:       firstRecipient,
		Amount:          total,
		EstimateFees:    strconv.FormatInt(effects.GasFee(), 10),
		FinishTimestamp: resp.TimestampMs.Int64() / 1000,
	}
	status := effects.V1.Status
	if status.Status == types.ExecutionStatusSuccess {
		detail.Status = base.TransactionStatusSuccess
	} else {
		detail.Status = base.TransactionStatusFailure
		detail.FailureMessage = status.Error
	}
	return
}

func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	digest, err := sui_types.NewDigest(hash)
	if err != nil {
		return base.TransactionStatusNone
	}
	cli, err := c.Client()
	if err != nil {
		return base.TransactionStatusNone
	}
	resp, err := cli.GetTransactionBlock(context.Background(), *digest, types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	})
	if err != nil {
		return base.TransactionStatusNone
	}
	if resp.Effects == nil {
		return base.TransactionStatusNone
	}
	if resp.Effects.Data.IsSuccess() {
		return base.TransactionStatusSuccess
	} else {
		return base.TransactionStatusFailure
	}
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

func (c *Chain) TransferObject(sender, receiver, objectId string, gasBudget int64) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := sui_types.NewAddressFromHex(sender)
	if err != nil {
		return
	}
	receiverAddress, err := sui_types.NewAddressFromHex(receiver)
	if err != nil {
		return
	}
	nftObject, err := sui_types.NewObjectIdFromHex(objectId)
	if err != nil {
		return nil, err
	}
	client, err := c.Client()
	if err != nil {
		return
	}
	object, err := client.GetObject(context.Background(), *nftObject, nil)
	if err != nil {
		return
	}
	gasUint := base.Max(uint64(gasBudget), MinGasBudget)
	pickedGasCoins, err := c.PickGasCoins(*signer, gasUint)
	if err != nil {
		return
	}

	maxGasBudget := maxGasBudget(pickedGasCoins, gasUint)
	gasPrice, _ := c.CachedGasPrice()
	return c.EstimateTransactionFeeAndRebuildTransactionBCS(maxGasBudget, func(gasBudget uint64) (*Transaction, error) {
		ptb := sui_types.NewProgrammableTransactionBuilder()
		reference := object.Data.Reference()
		err = ptb.TransferObject(*receiverAddress, []*sui_types.ObjectRef{&reference})
		if err != nil {
			return nil, err
		}

		pt := ptb.Finish()
		tx := sui_types.NewProgrammable(*signer, pickedGasCoins.CoinRefs(), pt, gasBudget, gasPrice)
		txBytes, err := bcs.Marshal(tx)
		if err != nil {
			return nil, err
		}
		return &Transaction{TxnBytes: txBytes}, nil
	})
}

func (c *Chain) GasPrice() (gasprice *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := c.Client()
	if err != nil {
		return
	}
	price, err := cli.GetReferenceGasPrice(context.Background())
	if err != nil {
		return
	}
	str := strconv.FormatUint(price.Uint64(), 10)
	return &base.OptionalString{Value: str}, nil
}

func (c *Chain) CachedGasPrice() (price uint64, err error) {
	if c.gasPrice > 0 {
		return c.gasPrice, nil
	}

	price = 1000
	defer base.CatchPanicAndMapToBasicError(&err)
	cli, err := c.Client()
	if err != nil {
		return
	}
	gasPrice, err := cli.GetReferenceGasPrice(context.Background())
	if err != nil {
		return
	}
	c.gasPrice = gasPrice.Uint64()
	return c.gasPrice, nil
}

func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	txn, ok := transaction.(*Transaction)
	if !ok {
		return nil, base.ErrInvalidTransactionType
	}

	defer base.CatchPanicAndMapToBasicError(&err)
	fee = &base.OptionalString{Value: strconv.FormatInt(MaxGasBudget, 10)}

	cli, err := c.Client()
	if err != nil {
		return
	}
	effects, err := cli.DryRunTransaction(context.Background(), txn.TransactionBytes())
	if err != nil {
		return
	}
	if !effects.Effects.Data.IsSuccess() {
		return nil, errors.New(effects.Effects.Data.V1.Status.Error)
	}

	gasFee := effects.Effects.Data.GasFee()
	if gasFee == 0 {
		gasFee = MaxGasBudget
	} else {
		gasFee = gasFee/10*11 + 10 // >= ceil(fee * 1.1)
	}
	txn.EstimateGasFee = gasFee
	gasString := strconv.FormatInt(gasFee, 10)
	return &base.OptionalString{Value: gasString}, nil
}
func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction base.Transaction, pubkey string) (fee *base.OptionalString, err error) {
	return c.EstimateTransactionFee(transaction)
}

/**
 * @param address Hex-encoded 16 bytes Sui account address wich mints tokens
 * @param faucetUrl default https://faucet.testnet.sui.io/gas
 * @return digest of the faucet transfer transaction.
 */
func FaucetFundAccount(address string, faucetUrl string) (h *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if faucetUrl == "" {
		faucetUrl = TestNetFaucetUrl
	}
	hash, err := client.FaucetFundAccount(address, faucetUrl)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: hash}, nil
}

func (c *Chain) PickGasCoins(owner sui_types.SuiAddress, maxGasBudget uint64) (picked *types.PickedCoins, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	cli, err := c.Client()
	if err != nil {
		return
	}
	coins, err := cli.GetCoins(context.Background(), owner, nil, nil, 30)
	if err != nil {
		return
	}
	return types.PickupCoins(coins, *big.NewInt(0), maxGasBudget, 30, 0)
}

func (c *Chain) EstimateTransactionFeeAndRebuildTransactionBCS(maxGasBudget uint64, buildTransaction func(gasBudget uint64) (*Transaction, error)) (*Transaction, error) {
	if maxGasBudget < MinGasBudget {
		maxGasBudget = MinGasBudget
	}
	txn, err := buildTransaction(maxGasBudget)
	if err != nil {
		return nil, err
	}
	_, err = c.EstimateTransactionFee(txn)
	if err != nil {
		return nil, err
	}
	if txn.EstimateGasFee < MinGasBudget {
		txn.EstimateGasFee = MinGasBudget
	}

	// second call the builder
	newTxn, err := buildTransaction(uint64(txn.EstimateGasFee))
	if err != nil {
		return nil, err
	}
	newTxn.EstimateGasFee = txn.EstimateGasFee
	return newTxn, nil
}

// @param maxGasBudget: the firstly build required gas
// @param builer: the builder should build a transaction, it maybe will invoking twice, the firstly build gas pass the maxGasBudget, the second build will pass the estimate gas.
func (c *Chain) EstimateTransactionFeeAndRebuildTransaction(maxGasBudget uint64, buildTransaction func(gasBudget uint64) (*Transaction, error)) (*Transaction, error) {
	if maxGasBudget < MinGasBudget {
		maxGasBudget = MinGasBudget
	}
	var (
		txn *Transaction = nil
		err error        = nil

		estimateFeeUint = uint64(0)
	)
	isLowGasError := func(err error) bool {
		regInsufficientGas := regexp.MustCompile(`Insufficient.*Gas|GasBudgetTooLow|.*is less than the reference gas price.*`)
		match := regInsufficientGas.FindAllStringSubmatch(err.Error(), -1)
		return len(match) > 0
	}
	isCannotFindGasError := func(err error) bool {
		reg := regexp.MustCompile(`Cannot find gas coin for signer.*[0-9a-f]+.*amount sufficient.*required gas amount.*\d+.*`)
		match := reg.FindAllStringSubmatch(err.Error(), -1)
		return len(match) > 0
	}
	count := 0
	for {
		if count >= 20 {
			return nil, errors.New("build transaction failed")
		}
		count++
		txn, err = buildTransaction(maxGasBudget)
		if err != nil {
			if isLowGasError(err) {
				maxGasBudget = nextTryingGas(maxGasBudget)
				continue
			}
			if isCannotFindGasError(err) {
				return nil, ErrNeedSplitGasCoin
			}
			return nil, err
		}
		_, err = c.EstimateTransactionFee(txn)
		if err != nil {
			if isLowGasError(err) {
				maxGasBudget = nextTryingGas(maxGasBudget)
				continue
			}
			if isCannotFindGasError(err) {
				return nil, ErrNeedSplitGasCoin
			}
			return nil, err
		}
		if txn.EstimateGasFee < MinGasBudget {
			estimateFeeUint = MinGasBudget
		} else {
			estimateFeeUint = uint64(txn.EstimateGasFee)
		}
		break
	}
	if estimateFeeUint/5*6 > maxGasBudget && (estimateFeeUint > maxGasBudget || maxGasBudget-estimateFeeUint < 1000000) {
		// estimate*1.2 > max && max-estimate < 0.001SUI
		// The estimated transaction fee is not much different from the build transaction.
		txn.EstimateGasFee = int64(maxGasBudget)
		return txn, nil
	}

	// second call the builder
	newTxn, err := buildTransaction(estimateFeeUint)
	if err != nil {
		return nil, err
	}
	newTxn.EstimateGasFee = txn.EstimateGasFee
	return newTxn, nil
}

func nextTryingGas(currentGas uint64) uint64 {
	if currentGas < MinGasBudget {
		return MinGasBudget
	} else if currentGas < 10e6 {
		return 10e6
	} else if currentGas < 30e6 {
		return 30e6
	} else if currentGas < 60e6 {
		return 60e6
	} else if currentGas < MaxGasBudget { // 90e6
		return MaxGasBudget
	} else {
		return currentGas / 2 * 3 // gas * 1.5
	}
}

func maxGasBudget(pickedCoins *types.PickedCoins, maxGasBudget uint64) uint64 {
	suggestGas := pickedCoins.SuggestMaxGasBudget()
	return base.Min(suggestGas, maxGasBudget)
}
