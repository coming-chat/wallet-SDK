package doge

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
