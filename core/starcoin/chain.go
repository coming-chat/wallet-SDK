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

func (c *Chain) SendSignedTransaction(signedTxn base.SignedTransaction) (hash *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	txn := AsSignedTransaction(signedTxn)
	if txn == nil {
		return nil, base.ErrInvalidTransactionType
	}
	hashString, err := c.client.SubmitSignedTransaction(context.Background(), txn.Txn)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: hashString}, nil
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

func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	return nil, base.ErrEstimateGasNeedPublicKey
}
func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction base.Transaction, pubkey string) (fee *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	pubData, err := hexTypes.HexDecodeString(pubkey)
	if err != nil {
		return nil, base.ErrInvalidPublicKey
	}
	txn := transaction.(*Transaction)
	if txn == nil {
		return nil, base.ErrInvalidTransactionType
	}
	gasFee, err := c.client.EstimateGasByDryRunRaw(context.Background(), *txn.Txn, types.Ed25519PublicKey(pubData))
	if err != nil {
		return nil, err
	}
	gasFee = big.NewInt(0).Div(big.NewInt(0).Mul(gasFee, big.NewInt(3)), big.NewInt(2)) // gasFee * 3 / 2
	txn.Txn.MaxGasAmount = gasFee.Uint64()
	return &base.OptionalString{Value: gasFee.String()}, nil
}

func (c *Chain) GasPrice() (*base.OptionalString, error) {
	price, err := c.client.GetGasUnitPrice(context.Background())
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: strconv.FormatInt(int64(price), 10)}, nil
}

func (c *Chain) GetState(context context.Context, address string) (*types.AccountResource, error) {
	state, err := c.client.GetState(context, address)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Bcs Deserialize AccountResource failed") {
			return nil, base.ErrInsufficientBalance
		}
		return nil, err
	}
	return state, nil
}
