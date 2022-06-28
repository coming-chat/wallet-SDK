package doge

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
	"github.com/coming-chat/wallet-SDK/util/hexutil"
)

func queryBalance(address, chainnet string) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	b = &base.Balance{
		Total:  "0",
		Usable: "0",
	}

	restUrl, err := restUrlOf(chainnet)
	if err != nil {
		return
	}
	if IsValidAddress(address, chainnet) == false {
		return b, errors.New("Invalid address")
	}

	// https://api.blockcypher.com/v1/doge/main/addrs/DBx1XSBxpSUnEK79nA8VtrKh2qr2LupZ6G/balance
	url := fmt.Sprintf("%v/addrs/%v/balance", restUrl, address)
	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return
	}
	if response.Code != http.StatusOK {
		return b, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}

	var balance = struct {
		Balance uint64 `json:"balance"`
	}{}
	err = json.Unmarshal(response.Body, &balance)
	if err != nil {
		return
	}

	b.Total = strconv.FormatUint(balance.Balance, 10)
	b.Usable = b.Total
	return b, nil
}

func fetchTransactionDetail(hash, chainnet string) (d *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	restUrl, err := restUrlOf(chainnet)
	if err != nil {
		return
	}

	// https://api.blockcypher.com/v1/doge/main/txs/7bc313903372776e1eb81d321e3fe27c9721ce8e71a9bcfee1bde6baea31b5c2
	hash = strings.TrimPrefix(hash, "0x")
	url := fmt.Sprintf("%v/txs/%v", restUrl, hash)
	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return
	}
	if response.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}

	var detail = Transaction{}
	err = json.Unmarshal(response.Body, &detail)
	return &detail, err
}

// @param limit Specify how many the latest utxos to fetch
func fetchUtxos(address, chainnet string, limit int) (l *UTXOList, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	restUrl, err := restUrlOf(chainnet)
	if err != nil {
		return
	}
	if IsValidAddress(address, chainnet) == false {
		return nil, errors.New("Invalid address")
	}

	// https://api.blockcypher.com/v1/doge/main/addrs/D8aDCsK4TA9NYhmwiqw1BjZ4CP8LQ814Ea?limit=5&unspentOnly=true
	url := fmt.Sprintf("%v/addrs/%v?limit=%v&unspentOnly=true", restUrl, address, limit)
	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return
	}
	if response.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}

	var list = UTXOList{}
	err = json.Unmarshal(response.Body, &list)
	return &list, err
}

func suggestFeeRate(chainnet string) (f *FeeRate, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	restUrl, err := restUrlOf(chainnet)
	if err != nil {
		return
	}

	// https://api.blockcypher.com/v1/doge/main
	response, err := httpUtil.Request(http.MethodGet, restUrl, nil, nil)
	if err != nil {
		return
	}
	if response.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}

	f = &FeeRate{}
	err = json.Unmarshal(response.Body, f)
	if err != nil {
		return nil, err
	}

	f.High = f.High / 1024
	f.Average = f.Average / 1024
	f.Low = f.Low / 1024
	return f, nil
}

func sendRawTransaction(txHex, chainnet string) (t *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	restUrl, err := restUrlOf(chainnet)
	if err != nil {
		return
	}
	txHex = strings.TrimPrefix(txHex, "0x")
	if !hexutil.ValidHex(txHex) {
		return nil, errors.New("Invalid raw transaction hex string")
	}

	url := restUrl + "/txs/push"
	params := map[string]string{
		"tx": txHex,
	}
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	response, err := httpUtil.Request(http.MethodPost, url, nil, data)
	if response.Code != http.StatusOK && response.Code != 201 {
		return nil, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}

	var res = struct {
		Tx Transaction `json:"tx"`
	}{}

	err = json.Unmarshal(response.Body, &res)
	return &res.Tx, err
}
