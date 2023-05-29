package btc

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/base/inter"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type StringMap struct {
	inter.AnyMap[string]
}

func NewStringMap() *StringMap {
	return &StringMap{inter.AnyMap[string]{Values: map[string]string{}}}
}

// BatchQueryBalance
// @return If any address is successfully queried, it will return normally, and the amount of failed request is 0
// @throw error if all address query balance failed
func BatchQueryBalance(addresses *base.StringArray, chainnet string) (*StringMap, error) {
	var (
		res     = NewStringMap()
		success = false
		anyErr  error
	)
	base.MapListConcurrentStringToString(addresses.Values, func(address string) (string, error) {
		balance, err := queryBalance(address, chainnet)
		if err != nil {
			anyErr = err
			res.SetValue("0", address)
		} else {
			success = true
			res.SetValue(balance, address)
		}
		return "", nil
	})
	if success {
		return res, nil
	} else {
		return nil, anyErr
	}
}

// queryBalance
// query the balance according to the address.
func queryBalance(address, chainnet string) (string, error) {
	host, err := scanHostOf(chainnet)
	if err != nil {
		return "0", err
	}
	url := host + "/address/" + address

	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return "0", base.MapAnyToBasicError(err)
	}

	return parseBalanceResponse(response)
}

// queryBalancePubkey
// query the balance according to the public key.
func queryBalancePubkey(pubkey, chainnet string) (string, error) {
	pubkey = strings.TrimPrefix(pubkey, "0x")
	host, err := scanHostOf(chainnet)
	if err != nil {
		return "0", err
	}
	url := host + "/pubkey/" + pubkey

	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return "0", base.MapAnyToBasicError(err)
	}

	return parseBalanceResponse(response)
}

func parseBalanceResponse(response *httpUtil.Res) (string, error) {
	if response.Code != http.StatusOK {
		return "0", fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}
	respDict := make(map[string]interface{})
	err := json.Unmarshal(response.Body, &respDict)
	if err != nil {
		return "0", err
	}

	chainStats, ok := respDict["chain_stats"].(map[string]interface{})
	if !ok {
		return "0", ErrHttpResponseParse
	}
	funded, ok1 := chainStats["funded_txo_sum"].(float64)
	spend, ok2 := chainStats["spent_txo_sum"].(float64)
	if !ok1 || !ok2 {
		return "0", ErrHttpResponseParse
	}

	balance := int64(math.Max(0, funded-spend))
	return strconv.FormatInt(balance, 10), nil
}

// Deprecated: QueryBalance is deprecated. Please Use Chain.QueryBalanceWithAddress() instead.
func QueryBalance(address, chainnet string) (string, error) {
	return queryBalance(address, chainnet)
}

// Deprecated: QueryBalancePubkey is deprecated. Please Use Chain.QueryBalanceWithPublicKey() instead.
func QueryBalancePubkey(pubkey, chainnet string) (string, error) {
	return queryBalancePubkey(pubkey, chainnet)
}
