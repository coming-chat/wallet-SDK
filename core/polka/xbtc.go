package polka

import (
	"errors"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type XBTCToken struct {
	chain *Chain
}

// Warning: initial unavailable, You must create based on Chain.XBTCToken()
func NewXBTCToken() (*Token, error) {
	return nil, errors.New("Token initial unavailable, You must create based on Chain.XBTCToken()")
}

func (c *Chain) XBTCToken() *XBTCToken {
	return &XBTCToken{chain: c}
}

// MARK - Implement the protocol Token, Override

func (t *XBTCToken) TokenInfo() (*base.TokenInfo, error) {
	return &base.TokenInfo{
		Name:    "XBTC",
		Symbol:  "XBTC",
		Decimal: 8,
	}, nil
}

func (t *XBTCToken) BalanceOfAddress(address string) (*base.Balance, error) {
	return t.chain.QueryBalanceXBTC(address)
}
func (t *XBTCToken) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	address, err := t.chain.EncodePublicKeyToAddress(publicKey)
	if err != nil {
		return nil, err
	}
	return t.chain.QueryBalanceXBTC(address)
}
func (t *XBTCToken) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.chain.QueryBalanceXBTC(account.Address())
}
