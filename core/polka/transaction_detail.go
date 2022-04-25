package polka

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

// MARK - Implement the protocol Chain.TransactionDetail

func (c *Chain) FetchTransactionDetail(hashString string) (*base.TransactionDetail, error) {
	if c.ScanUrl == "" {
		return nil, errors.New("Scan url is Empty.")
	}
	url := strings.TrimSuffix(c.ScanUrl, "/") + "/" + hashString

	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}
	respDict := make(map[string]interface{})
	err = json.Unmarshal(response.Body, &respDict)
	if err != nil {
		return nil, err
	}

	// decode informations
	amount, _ := respDict["txAmount"].(string)
	fee, _ := respDict["fee"].(string)
	from, _ := respDict["signer"].(string)
	to, _ := respDict["txTo"].(string)
	timestamp, _ := respDict["blockTime"].(float64)

	status := base.TransactionStatusNone
	finalized, _ := respDict["finalized"].(bool)
	if finalized {
		success, _ := respDict["success"].(bool)
		if success {
			status = base.TransactionStatusSuccess
		} else {
			status = base.TransactionStatusFailure
		}
	} else {
		status = base.TransactionStatusPending
	}

	return &base.TransactionDetail{
		HashString:      hashString,
		Amount:          amount,
		EstimateFees:    fee,
		FromAddress:     from,
		ToAddress:       to,
		Status:          status,
		FinishTimestamp: int64(timestamp),
	}, nil
}

func (c *Chain) FetchTransactionStatus(hashString string) base.TransactionStatus {
	detail, err := c.FetchTransactionDetail(hashString)
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
