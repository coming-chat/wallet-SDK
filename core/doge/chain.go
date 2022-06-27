package doge

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

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
	return queryBalance(address, c.Chainnet)
}
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	address, err := EncodePublicKeyToAddress(publicKey, c.Chainnet)
	if err != nil {
		return nil, err
	}
	return c.BalanceOfAddress(address)
}
func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (string, error) {
	transaction, err := sendRawTransaction(signedTx, c.Chainnet)
	if err != nil {
		return "", err
	}
	return transaction.Hash, nil
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	d, err := fetchTransactionDetail(hash, c.Chainnet)
	if err != nil {
		return nil, err
	} else {
		return d.SdkDetail(), nil
	}
}

func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	d, err := fetchTransactionDetail(hash, c.Chainnet)
	if err != nil {
		return base.TransactionStatusFailure
	} else {
		return d.Status()
	}
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

// @param limit Specify how many the latest utxos to fetch, The minimum value of the limit is 100.
func (c *Chain) FetchUtxos(address string, limit int) (*base.OptionalString, error) {
	if limit < 100 {
		limit = 100
	}
	res, err := fetchUtxos(address, c.Chainnet, limit)
	if err != nil {
		return nil, err
	}

	utxos := res.Utxos
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Value.Cmp(utxos[j].Value) == 1
	})

	feeRate, err := c.SuggestFeeRate()
	if err != nil {
		return nil, err
	}

	sdklist := &SDKUTXOList{
		Txids:      utxos,
		FastestFee: int(feeRate.Average),
	}
	data, err := json.Marshal(sdklist)
	if err != nil {
		return nil, err
	}

	return &base.OptionalString{Value: string(data)}, nil
}

type FeeRate struct {
	Low     int64 `json:"low_fee_per_kb"`
	Average int64 `json:"medium_fee_per_kb"`
	High    int64 `json:"high_fee_per_kb"`
}

func (c *Chain) SuggestFeeRate() (*FeeRate, error) {
	return suggestFeeRate(c.Chainnet)
}
