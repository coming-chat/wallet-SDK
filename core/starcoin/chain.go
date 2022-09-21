package starcoin

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
	"strings"

	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/bcs"
	"github.com/starcoinorg/starcoin-go/client"
	"github.com/starcoinorg/starcoin-go/types"
)

const (
	MaxGasAmount = client.DEFAULT_MAX_GAS_AMOUNT

	txnSuccessStatus = "Executed"
)

type Chain struct {
	RpcUrl string
	client client.StarcoinClient
}

func NewChainWithRpc(rpcUrl string) *Chain {
	return &Chain{
		RpcUrl: rpcUrl,
		client: client.NewStarcoinClient(rpcUrl),
	}
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return NewMainToken(c)
}

func (c *Chain) BalanceOfAddress(address string) (b *base.Balance, err error) {
	return c.MainToken().BalanceOfAddress(address)
}
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return c.MainToken().BalanceOfPublicKey(publicKey)
}
func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.MainToken().BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (hash string, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	txnBytes, err := hexTypes.HexDecodeString(signedTx)
	if err != nil {
		return
	}
	return c.client.SubmitSignedTransactionBytes(context.Background(), txnBytes)
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	txn, err := c.client.GetTransactionByHash(context.Background(), hash)
	if err != nil {
		return
	}
	if txn.BlockHash == "" {
		return nil, errors.New("The transaction is still pending ")
	}
	rawTxn := txn.UserTransaction.RawTransaction
	payloadBytes, err := hexTypes.HexDecodeString(rawTxn.Payload)
	if err != nil {
		return
	}
	payload, err := types.BcsDeserializeTransactionPayload(payloadBytes)
	if err != nil {
		return
	}
	txnPayload, ok := payload.(*types.TransactionPayload__ScriptFunction)
	val := txnPayload.Value
	if !ok || val.Module.Name != "TransferScripts" || val.Function != "peer_to_peer_v2" || len(val.Args) != 2 {
		return nil, errors.New("Invalid transfer transaction.")
	}
	amount, err := bcs.NewDeserializer(txnPayload.Value.Args[1]).DeserializeU128()
	if err != nil {
		return
	}
	detail = &base.TransactionDetail{
		HashString:      hash,
		FromAddress:     rawTxn.Sender,
		ToAddress:       "0x" + hex.EncodeToString(val.Args[0]),
		Amount:          client.U128ToBigInt(&amount).String(),
		EstimateFees:    "0",
		FinishTimestamp: 0,
		Status:          base.TransactionStatusNone,
		FailureMessage:  "",
	}

	gasPrice, ok := big.NewInt(0).SetString(rawTxn.GasUnitPrice, 10)
	if !ok {
		gasPrice = big.NewInt(1)
	}
	maxGas, ok := big.NewInt(0).SetString(rawTxn.MaxGasAmount, 10)
	if !ok {
		maxGas = big.NewInt(0)
	}
	detail.EstimateFees = big.NewInt(0).Mul(gasPrice, maxGas).String()

	txnInfo, err := c.client.GetTransactionInfoByHash(context.Background(), hash)
	if err != nil {
		return detail, nil
	}
	status := ""
	err = json.Unmarshal(txnInfo.Status, &status)
	if err != nil {
		detail.Status = base.TransactionStatusNone
	} else if status == txnSuccessStatus {
		detail.Status = base.TransactionStatusSuccess
	} else {
		detail.Status = base.TransactionStatusFailure
		detail.FailureMessage = status
	}
	gasUsed, ok := big.NewInt(0).SetString(txnInfo.GasUsed, 10)
	if !ok {
		gasUsed = maxGas
	}
	detail.EstimateFees = big.NewInt(0).Mul(gasPrice, gasUsed).String()

	blockInfo, err := c.client.GetBlockByHash(context.Background(), txnInfo.BlockHash)
	if err != nil {
		return detail, nil
	}
	timeInt, err := strconv.ParseInt(blockInfo.BlockHeader.Timestamp, 10, 64)
	if err != nil {
		timeInt = 0
	}
	detail.FinishTimestamp = timeInt / 1000
	return detail, nil
}

func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	txnInfo, err := c.client.GetTransactionInfoByHash(context.Background(), hash)
	if err != nil {
		return base.TransactionStatusNone
	}
	status := ""
	err = json.Unmarshal(txnInfo.Status, &status)
	if err != nil {
		return base.TransactionStatusNone
	} else if status == txnSuccessStatus {
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

func (c *Chain) GasPrice() (*base.OptionalString, error) {
	price, err := c.client.GetGasUnitPrice(context.Background())
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: strconv.FormatInt(int64(price), 10)}, nil
}

func (c *Chain) BuildRawUserTransaction(from *Account, payload types.TransactionPayload) (txn *types.RawUserTransaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	ctx := context.Background()
	price, err := c.client.GetGasUnitPrice(ctx)
	if err != nil {
		return
	}
	state, err := c.client.GetState(ctx, from.Address())
	if err != nil {
		return
	}
	rawTxn, err := c.client.BuildRawUserTransaction(ctx, from.AccountAddress(), payload, price, MaxGasAmount, state.SequenceNumber)
	if err != nil {
		return
	}
	gasLimit, err := c.client.EstimateGasByDryRunRaw(ctx, *rawTxn, from.PublicKey())
	if err != nil {
		return
	}
	rawTxn.MaxGasAmount = big.NewInt(0).Div(big.NewInt(0).Mul(gasLimit, big.NewInt(3)), big.NewInt(2)).Uint64() // gaslimit * 1.5
	return rawTxn, nil
}
