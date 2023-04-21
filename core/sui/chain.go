package sui

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/coming-chat/go-sui/client"
	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const (
	MaxGasBudget = 90000000
	// MaxGasForMerge = 10000000

	MaxGasForPay      = 10000000
	MaxGasForTransfer = 10000000

	TestNetFaucetUrl = client.TestNetFaucetUrl
	DevNetFaucetUrl  = client.DevNetFaucetUrl
)

type Chain struct {
	rpcClient *client.Client
	RpcUrl    string
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

	bytes, err := types.NewBase64Data(signedTx)
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
	response, err := cli.ExecuteTransactionBlock(context.Background(), *signedTxn.TxBytes, []any{signedTxn.Signature}, &options, types.TxnRequestTypeWaitForLocalExecution)
	if err != nil {
		return
	}
	hash = response.Digest
	if !response.Effects.Data.IsSuccess() {
		return hash, errors.New(response.Effects.Data.V1.Status.Error)
	}
	return hash, nil
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := c.Client()
	if err != nil {
		return
	}
	resp, err := cli.GetTransactionBlock(context.Background(), hash, types.SuiTransactionBlockResponseOptions{
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
	if len(transaction.Inputs) != 2 {
		return nil, notCoinTransferErr
	}
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
	detail, err := c.FetchTransactionDetail(hash)
	if err != nil {
		return base.TransactionStatusNone
	}
	return detail.Status
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

// @param gasId gas object to be used in this transaction, the gateway will pick one from the signer's possession if not provided
func (c *Chain) TransferObject(sender, receiver, objectId, gasId string, gasBudget int64) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	senderAddress, err := types.NewAddressFromHex(sender)
	if err != nil {
		return
	}
	receiverAddress, err := types.NewAddressFromHex(receiver)
	if err != nil {
		return
	}
	nftObject, err := types.NewHexData(objectId)
	if err != nil {
		return nil, err
	}
	var gas *types.ObjectId = nil
	if gasId != "" {
		gas, err = types.NewHexData(gasId)
		if err != nil {
			return nil, errors.New("Invalid gas object id")
		}
	}
	client, err := c.Client()
	if err != nil {
		return
	}
	gasInt := types.NewSafeSuiBigInt(uint64(gasBudget))
	tx, err := client.TransferObject(context.Background(), *senderAddress, *receiverAddress, *nftObject, gas, gasInt)
	if err != nil {
		return
	}
	return &Transaction{
		Txn:          *tx,
		MaxGasBudget: gasBudget,
	}, nil
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

func (c *Chain) EstimateGasFee(transaction *Transaction) (fee *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	fee = &base.OptionalString{Value: strconv.FormatInt(MaxGasBudget, 10)}

	cli, err := c.Client()
	if err != nil {
		return
	}
	effects, err := cli.DryRunTransaction(context.Background(), &transaction.Txn)
	if err != nil {
		return
	}

	gasFee := effects.Effects.Data.GasFee()
	if gasFee == 0 {
		gasFee = MaxGasBudget
	} else {
		gasFee = gasFee/10*15 + 14 // >= ceil(fee * 1.5)
	}
	transaction.EstimateGasFee = gasFee
	gasString := strconv.FormatInt(gasFee, 10)
	return &base.OptionalString{Value: gasString}, nil
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
