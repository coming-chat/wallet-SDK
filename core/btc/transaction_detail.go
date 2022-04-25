package btc

import (
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/coming-chat/wallet-SDK/core/base"
)

// Note: The input parsing of bitcoin is very complex and the network cost is relatively high,
// So only the status and timestamp can be queried.
func fetchTransactionDetail(hashString, chainnet string) (*base.TransactionDetail, error) {
	client, err := rpcClientOf(chainnet)
	if err != nil {
		return nil, err
	}

	hash, err := chainhash.NewHashFromStr(hashString)
	if err != nil {
		return nil, err
	}

	rawResult, err := client.GetRawTransactionVerbose(hash)
	if err != nil {
		return nil, base.MapToBasicError(err)
	}

	status := base.TransactionStatusPending
	if rawResult.Confirmations > 0 {
		status = base.TransactionStatusSuccess
	}
	return &base.TransactionDetail{
		HashString:      hashString,
		Status:          status,
		FinishTimestamp: rawResult.Time,
	}, nil
}

func fetchTransactionStatus(hashString string, chainnet string) base.TransactionStatus {
	detail, err := FetchTransactionDetail(hashString, chainnet)
	if err != nil {
		return base.TransactionStatusNone
	}
	return detail.Status
}

func batchTransactionStatus(hashList []string, chainnet string) ([]string, error) {
	return base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(FetchTransactionStatus(s, chainnet)), nil
	})
}

// Batch function for mobile sdk for iOS & Android.
// Because the api exported by gomobile does not support arrays.
func sdkBatchTransactionStatus(hashListString string, chainnet string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(FetchTransactionStatus(s, chainnet)), nil
	})
	return strings.Join(statuses, ",")
}

// Deprecated: FetchTransactionDetail is deprecated. Please Use Chain.FetchTransactionDetail() instead.
func FetchTransactionDetail(hashString, chainnet string) (*base.TransactionDetail, error) {
	return fetchTransactionDetail(hashString, chainnet)
}

// Deprecated: FetchTransactionStatus is deprecated. Please Use Chain.FetchTransactionStatus() instead.
func FetchTransactionStatus(hashString string, chainnet string) base.TransactionStatus {
	return fetchTransactionStatus(hashString, chainnet)
}

// Deprecated: SdkBatchTransactionStatus is deprecated. Please Use Chain.BatchFetchTransactionStatus() instead.
func SdkBatchTransactionStatus(hashListString string, chainnet string) string {
	return sdkBatchTransactionStatus(hashListString, chainnet)
}
