package cosmos

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/cosmos/cosmos-sdk/types/tx"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	tendermintHttp "github.com/tendermint/tendermint/rpc/client/http"
	tendermintTypes "github.com/tendermint/tendermint/rpc/core/types"
)

type Chain struct {
	RpcUrl  string
	RestUrl string

	client *tendermintHttp.HTTP
}

func NewChainWithRpc(rpcUrl string, restUrl string) *Chain {
	return &Chain{
		RpcUrl:  rpcUrl,
		RestUrl: restUrl,
	}
}

// MARK - Implement the protocol Chain

// Warning Cosmos not supported main token
func (c *Chain) MainToken() base.Token {
	return nil
}

func (c *Chain) DenomToken(prefix, denom string) *Token {
	return &Token{chain: c, Denom: denom, Prefix: prefix}
}

func (c *Chain) BalanceOfAddress(address string) (*base.Balance, error) {
	return c.BalanceOfAddressAndDenom(address, "")
}

// Warning: Unable to use public key to query balance
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return base.EmptyBalance(), errors.New("Unable to use public key to query balance")
}

func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (string, error) {
	client, err := c.GetClient()
	if err != nil {
		return "", err
	}

	txBytes, err := types.HexDecodeString(signedTx)
	if err != nil {
		return "", err
	}
	commit, err := client.BroadcastTxSync(context.Background(), txBytes)
	if err != nil {
		return "", err
	}
	return commit.Hash.String(), nil
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	result, err := c.fetchTxResult(hash)
	if err != nil {
		return
	}

	detail = &base.TransactionDetail{}
	detail.HashString = hash
	if result.TxResult.Data == nil {
		detail.Status = base.TransactionStatusFailure
		detail.FailureMessage = result.TxResult.Log
	} else {
		detail.Status = base.TransactionStatusSuccess
	}

	originTx := &tx.Tx{}
	err = originTx.XXX_Unmarshal(result.Tx)
	if err != nil {
		return
	}
	detail.EstimateFees = originTx.AuthInfo.Fee.Amount[0].Amount.String()

	msgSend := &bankTypes.MsgSend{}
	err = msgSend.XXX_Unmarshal(originTx.Body.Messages[0].Value)
	if err != nil {
		return
	}
	detail.FromAddress = msgSend.FromAddress
	detail.ToAddress = msgSend.ToAddress
	detail.Amount = msgSend.Amount[0].Amount.String()

	client, err := c.GetClient()
	if err != nil {
		return
	}
	blockHeight := result.Height
	blockInfo, err := client.Block(context.Background(), &blockHeight)
	if err != nil {
		return
	}
	detail.FinishTimestamp = blockInfo.Block.Time.Unix()

	return detail, nil
}

// Fetch transaction status through transaction hash
func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	result, err := c.fetchTxResult(hash)
	if err != nil {
		return base.TransactionStatusNone
	}
	if result.TxResult.Data == nil {
		return base.TransactionStatusFailure
	} else {
		return base.TransactionStatusSuccess
	}
}

func (c *Chain) fetchTxResult(hash string) (*tendermintTypes.ResultTx, error) {
	client, err := c.GetClient()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("tx.hash='%s'", strings.TrimPrefix(hash, "0x"))
	txSearch, err := client.TxSearch(context.Background(), query, true, nil, nil, "asc")
	if err != nil {
		return nil, err
	}
	if len(txSearch.Txs) <= 0 {
		return nil, errors.New("Transaction not found: " + hash)
	}
	return txSearch.Txs[0], nil
}

// Batch fetch the transaction status, the hash list and the return value,
// which can only be passed as strings separated by ","
// @param hashListString The hash of the transactions to be queried in batches, a string concatenated with ",": "hash1,hash2,hash3"
// @return Batch transaction status, its order is consistent with hashListString: "status1,status2,status3"
func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

// MARK - Client

func (c *Chain) GetClient() (*tendermintHttp.HTTP, error) {
	if c.client != nil {
		return c.client, nil
	}

	client, err := tendermintHttp.New(c.RpcUrl, "/websocket")
	if err == nil {
		c.client = client
	}
	return client, base.MapAnyToBasicError(err)
}
