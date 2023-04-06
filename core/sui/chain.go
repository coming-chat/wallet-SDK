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
	MaxGasBudget   = 10000000
	MaxGasForMerge = 10000000

	MaxGasForPay      = 10000000
	MaxGasForTransfer = 10000000

	TestNetFaucetUrl = client.TestNetFaucetUrl
	DevNetFaucetUrl  = client.DevNetFaucetUrl
)

var (
	_ IChain = &Chain{} // check implement
)

type IChain interface {
	base.Chain
	GetClient() (*client.Client, error)
	EstimateGasFee(transaction *Transaction) (fee *base.OptionalString, err error)
}

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
	response, err := cli.ExecuteTransactionBlock(context.Background(), *signedTxn.TxBytes, []any{signedTxn.Signature}, &options, types.TxnRequestTypeWaitForEffectsCert)
	if err != nil {
		return
	}
	hash = response.Digest
	if !response.Effects.IsSuccess() {
		return hash, errors.New(response.Effects.Status.Error)
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
	})
	if err != nil {
		return nil, err
	}

	var notCoinTransferErr = errors.New("Invalid coin transfer transaction.")
	var firstRecipient string
	var total string
	var data = resp.Transaction.Data
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

	detail = &base.TransactionDetail{
		HashString:      hash,
		FromAddress:     data.Sender.ShortString(),
		ToAddress:       firstRecipient,
		Amount:          total,
		EstimateFees:    strconv.FormatUint(resp.Effects.GasFee(), 10),
		FinishTimestamp: int64(*resp.TimestampMs / 1000),
	}
	status := resp.Effects.Status
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
	tx, err := client.TransferObject(context.Background(), *senderAddress, *receiverAddress, *nftObject, gas, uint64(gasBudget))
	if err != nil {
		return
	}
	return &Transaction{
		Txn:          *tx,
		MaxGasBudget: gasBudget,
	}, nil
}

func (c *Chain) GasPrice() (*base.OptionalString, error) {
	cli, err := c.Client()
	if err != nil {
		return nil, err
	}
	price, err := cli.GetReferenceGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: price.String()}, nil
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

	gasFee := effects.Effects.GasFee()
	if gasFee == 0 {
		gasFee = MaxGasBudget
	} else {
		gasFee = gasFee/10*15 + 14 // >= ceil(fee * 1.5)
	}
	transaction.EstimateGasFee = int64(gasFee)
	gasString := strconv.FormatUint(gasFee, 10)
	return &base.OptionalString{Value: gasString}, nil
}

// MARK - Implement the protocol IChain
func (c *Chain) GetClient() (*client.Client, error) {
	return c.rpcClient, nil
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
