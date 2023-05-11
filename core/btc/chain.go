package btc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
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

func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}
func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction base.Transaction, pubkey string) (fee *base.OptionalString, err error) {
	return c.EstimateTransactionFee(transaction)
}

type FeeRate struct {
	Low     int64
	Average int64
	High    int64
}

func SuggestFeeRate() (*FeeRate, error) {
	url := "https://mempool-mainnet.coming.chat/api/v1/fees/recommended"

	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, base.MapAnyToBasicError(err)
	}

	if response.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}
	respDict := make(map[string]interface{})
	err = json.Unmarshal(response.Body, &respDict)
	if err != nil {
		return nil, err
	}

	var low, avg, high float64
	var ok bool
	if low, ok = respDict["minimumFee"].(float64); !ok {
		low = 1
	}
	if avg, ok = respDict["halfHourFee"].(float64); !ok {
		avg = low
	}
	if high, ok = respDict["fastestFee"].(float64); !ok {
		high = avg
	}
	return &FeeRate{
		Low:     int64(low),
		Average: int64(avg),
		High:    int64(high),
	}, nil
}
