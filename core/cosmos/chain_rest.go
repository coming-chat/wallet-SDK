package cosmos

import (
	"encoding/json"

	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type AccountInfo struct {
	AccountNumber string `json:"account_number"`
	Sequence      string `json:"sequence"`
}

func (c *Chain) AccountOf(address string) (*AccountInfo, error) {
	url := c.RestUrl + "/cosmos/auth/v1beta1/accounts/" + address
	body, err := httpUtil.Get(url, nil)
	if err != nil {
		return nil, err
	}

	account := struct {
		Account AccountInfo `json:"account"`
	}{}
	err = json.Unmarshal(body, &account)
	if err != nil {
		return nil, err
	}

	return &account.Account, nil
}
