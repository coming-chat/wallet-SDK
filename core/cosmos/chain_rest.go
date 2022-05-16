package cosmos

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type AccountInfo struct {
	AccountNumber string `json:"account_number"`
	Sequence      string `json:"sequence"`
}

// Is aslias of Sequence
func (i *AccountInfo) Nonce() string {
	return i.Sequence
}

type denomBalance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
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

func (c *Chain) BalanceOfAddressAndDenom(address, denom string) (b *base.Balance, err error) {
	b = base.EmptyBalance()

	// url := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s/%s", c.RestUrl, address, denom) // the api is unusable.
	url := fmt.Sprintf("%s/cosmos/bank/v1beta1/balances/%s", c.RestUrl, address)
	body, err := httpUtil.Get(url, nil)
	if err != nil {
		return
	}

	balances := struct {
		Balances []denomBalance `json:"balances"`
	}{}
	err = json.Unmarshal(body, &balances)
	if err != nil {
		return b, errors.New("The balance cannot be found, the account may not exist")
	}

	if len(balances.Balances) <= 0 {
		return
	}

	var balance *denomBalance = nil
	if len(denom) <= 0 {
		// If no denom is specified, get the first balance
		balance = &balances.Balances[0]
	} else {
		for _, bal := range balances.Balances {
			if bal.Denom == denom {
				balance = &bal
				break
			}
		}
	}

	if balance == nil {
		return b, errors.New("Unmatched coin: " + denom)
	}

	b.Total = balance.Amount
	b.Usable = balance.Amount
	return b, nil
}
