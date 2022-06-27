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
