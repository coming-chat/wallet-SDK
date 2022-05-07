package btc

import (
	"github.com/coming-chat/wallet-SDK/core/base"
)

type Chain struct {
	*Util
}

func NewChainWithChainnet(chainnet string) (*Chain, error) {
	util, err := NewUtilWithChainnet(chainnet)
	if err != nil {
		return nil, err
	}

	return &Chain{Util: util}, nil
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return c
}

func (c *Chain) BalanceOfAddress(address string) (*base.Balance, error) {
	b, err := queryBalance(address, c.Chainnet)
	if err != nil {
		return nil, err
	}
	return &base.Balance{
		Total:  b,
		Usable: b,
	}, nil
}
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	b, err := queryBalancePubkey(publicKey, c.Chainnet)
	if err != nil {
		return nil, err
	}
	return &base.Balance{
		Total:  b,
		Usable: b,
	}, nil
}
func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfPublicKey(account.PublicKeyHex())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (string, error) {
	return sendRawTransaction(signedTx, c.Chainnet)
}

// Fetch transaction details through transaction hash
// Note: The input parsing of bitcoin is very complex and the network cost is relatively high,
// So only the status and timestamp can be queried.
func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	return fetchTransactionDetail(hash, c.Chainnet)
}

func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	return fetchTransactionStatus(hash, c.Chainnet)
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	return sdkBatchTransactionStatus(hashListString, c.Chainnet)
}