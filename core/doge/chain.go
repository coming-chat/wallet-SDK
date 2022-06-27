package doge

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

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
	// return sendRawTransaction(signedTx, c.Chainnet)
	return "", nil
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

	sdklist := &SDKUTXOList{
		Txids:      utxos,
		FastestFee: 1,
	}
	data, err := json.Marshal(sdklist)
	if err != nil {
		return nil, err
	}

	return &base.OptionalString{Value: string(data)}, nil
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
