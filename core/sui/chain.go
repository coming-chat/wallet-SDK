package sui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/coming-chat/go-sui/client"
	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const (
	MaxGasBudget = 2000

	MaxGasForMerge    = 500
	MaxGasForTransfer = 100
)

type Chain struct {
	rpcClient *client.Client
	RpcUrl    string
}

func NewChainWithRpcUrl(rpcUrl string) *Chain {
	return &Chain{RpcUrl: rpcUrl}
}

func (c *Chain) client() (*client.Client, error) {
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
	signedTxn := types.SignedTransaction{}
	err = json.Unmarshal(bytes.Data(), &signedTxn)
	if err != nil {
		return
	}
	cli, err := c.client()
	if err != nil {
		return
	}
	response, err := cli.ExecuteTransaction(context.Background(), signedTxn)
	if err != nil {
		return
	}
	if response.Effects.Status.Status != types.TransactionStatusSuccess {
		return "", fmt.Errorf(`chain error: %v`, response.Effects.Status.Error)
	}

	hash = response.Certificate.TransactionDigest.String()
	return
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	digest, err := types.NewBase64Data(hash)
	if err != nil {
		return
	}
	cli, err := c.client()
	if err != nil {
		return
	}
	resp, err := cli.GetTransaction(context.Background(), *digest)
	if err != nil {
		return nil, err
	}

	var firstRecipient *types.HexData
	var total uint64
	for _, txn := range resp.Certificate.Data.Transactions {
		tsui := txn.TransferSui
		if tsui == nil {
			continue
		}
		if firstRecipient == nil {
			firstRecipient = &tsui.Recipient
			total = tsui.Amount
		} else if bytes.Compare(firstRecipient.Data(), tsui.Recipient.Data()) == 0 {
			total = total + tsui.Amount
		}
	}
	if firstRecipient == nil {
		return nil, errors.New("Invalid coin transfer transaction.")
	}

	detail = &base.TransactionDetail{
		HashString:      hash,
		FromAddress:     resp.Certificate.Data.Sender.String(),
		ToAddress:       firstRecipient.String(),
		Amount:          strconv.FormatUint(total, 10),
		EstimateFees:    strconv.FormatUint(resp.Effects.GasFee(), 10),
		FinishTimestamp: int64(resp.TimestampMs / 1000),
	}
	status := resp.Effects.Status
	if status.Status == types.TransactionStatusSuccess {
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
