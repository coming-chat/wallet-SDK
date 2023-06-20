package starknet

import (
	"errors"

	"github.com/coming-chat/wallet-SDK/core/base"
)

const (
	ETHTokenContractAddressMainnet = ""
	ETHTokenContractAddressGoerli  = "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"
)

type Token struct {
	chain           *Chain
	ContractAddress string
}

func NewToken(chain *Chain, contractAddress string) *Token {
	return &Token{
		chain:           chain,
		ContractAddress: contractAddress,
	}
}

// MARK - Implement the protocol Token

func (t *Token) Chain() base.Chain {
	return t.chain
}

// Warning: Main token does not support
func (t *Token) TokenInfo() (*base.TokenInfo, error) {
	return nil, errors.New("Main token does not support")
}

func (t *Token) BalanceOfAddress(address string) (*base.Balance, error) {
	return t.chain.BalanceOf(address, t.ContractAddress)
}

func (t *Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	address, err := encodePublicKeyToAddressArgentX(publicKey)
	if err != nil {
		return nil, err
	}
	return t.chain.BalanceOf(address, t.ContractAddress)
}

func (t *Token) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.chain.BalanceOf(account.Address(), t.ContractAddress)
}

func (t *Token) BuildTransfer(sender, receiver, amount string) (txn base.Transaction, err error) {
	return nil, nil
}

// Before invoking this method, it is best to check `CanTransferAll()`
func (t *Token) CanTransferAll() bool {
	return false
}
func (t *Token) BuildTransferAll(sender, receiver string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}
