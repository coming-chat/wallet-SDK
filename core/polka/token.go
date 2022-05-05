package polka

import (
	"errors"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type Token struct {
	chain *Chain
}

// Warning: initial unavailable, You must create based on Chain.MainToken()
func NewToken() (*Token, error) {
	return nil, errors.New("Token initial unavailable, You must create based on Chain.MainToken()")
}

// MARK - Implement the protocol Token

func (t *Token) Chain() base.Chain {
	return t.chain
}

// Warning: polka chain is not currently supported
func (t *Token) TokenInfo() (*base.TokenInfo, error) {
	return nil, errors.New("Polka chain is not currently supported")
}

func (t *Token) BalanceOfAddress(address string) (*base.Balance, error) {
	return t.chain.BalanceOfAddress(address)
}
func (t *Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return t.chain.BalanceOfPublicKey(publicKey)
}
func (t *Token) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.chain.BalanceOfAccount(account)
}
