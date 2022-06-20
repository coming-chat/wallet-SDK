package doge

import "github.com/coming-chat/wallet-SDK/core/base"

// MARK - Implement the protocol Token

func (c *Chain) Chain() base.Chain {
	return c
}

func (c *Chain) TokenInfo() (*base.TokenInfo, error) {
	name, err := nameOf(c.Chainnet)
	if err != nil {
		return nil, err
	}
	return &base.TokenInfo{
		Name:    name,
		Symbol:  name,
		Decimal: 8,
	}, nil
}
