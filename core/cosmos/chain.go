package cosmos

import (
	"context"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	tendermintHttp "github.com/tendermint/tendermint/rpc/client/http"
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

// TODO
func (c *Chain) BalanceOfAddress(address string) (*base.Balance, error) {
	return nil, nil
}

// TODO
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return c.BalanceOfAddress(publicKey)
}

func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
// TODO
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
// TODO
func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	return nil, nil
}

// Fetch transaction status through transaction hash
func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	return 0
}

// Batch fetch the transaction status, the hash list and the return value,
// which can only be passed as strings separated by ","
// @param hashListString The hash of the transactions to be queried in batches, a string concatenated with ",": "hash1,hash2,hash3"
// @return Batch transaction status, its order is consistent with hashListString: "status1,status2,status3"
// TODO
func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	return ""
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
