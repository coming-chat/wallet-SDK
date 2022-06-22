package eth

import (
	"crypto/ecdsa"
	"errors"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

type Erc20TokenInfo struct {
	*base.TokenInfo

	ContractAddress string
	ChainId         string
	TokenIcon       string

	// Deprecated: Balance is not a token's info.
	Balance string
}

type Erc20Token struct {
	*Token

	ContractAddress string
}

// Warning: initial unavailable, You must create based on Chain.Erc20Token()
func NewErc20Token() (*Erc20Token, error) {
	return nil, errors.New("Token initial unavailable, You must create based on Chain.MainToken()")
}

// MARK - Implement the protocol Token, Override

// cannot get balance
func (t *Erc20Token) Erc20TokenInfo() (*Erc20TokenInfo, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return nil, err
	}
	baseInfo, err := t.TokenInfo()
	if err != nil {
		return nil, err
	}

	return &Erc20TokenInfo{
		TokenInfo:       baseInfo,
		ContractAddress: t.ContractAddress,
		ChainId:         chain.chainId.String(),
	}, nil
}

func (t *Erc20Token) TokenInfo() (*base.TokenInfo, error) {
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return nil, err
	}

	info := &base.TokenInfo{}
	info.Name, err = chain.TokenName(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	info.Symbol, err = chain.TokenSymbol(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	info.Decimal, err = chain.TokenDecimal(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (t *Erc20Token) BalanceOfAddress(address string) (*base.Balance, error) {
	b := base.EmptyBalance()
	chain, err := GetConnection(t.chain.RpcUrl)
	if err != nil {
		return b, err
	}

	balance, err := chain.TokenBalance(t.ContractAddress, address)
	if err != nil {
		return b, err
	}
	return &base.Balance{
		Total:  balance,
		Usable: balance,
	}, nil
}

func (t *Erc20Token) EstimateGasFeeLayer2(msg *CallMsg) (*OptimismLayer2Gas, error) {
	data, err := EncodeErc20Transfer(msg.GetTo(), msg.GetValue())
	if err != nil {
		return nil, err
	}
	msg.SetData(data)
	return t.Token.EstimateGasFeeLayer2(msg)
}

func (t *Erc20Token) EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string) (string, error) {
	msg := NewCallMsg()
	msg.SetFrom(fromAddress)
	msg.SetTo(t.ContractAddress)
	msg.SetGasPrice(gasPrice)

	data, err := EncodeErc20Transfer(receiverAddress, amount)
	if err != nil {
		return "", err
	}
	msg.SetData(data)

	gasLimit, err := t.chain.EstimateGasLimit(msg)
	if err != nil {
		return "", err
	}
	gasLimitDecimal, err := decimal.NewFromString(gasLimit.Value)
	if err != nil {
		return "", err
	}
	return gasLimitDecimal.Mul(decimal.NewFromFloat32(1.3)).Round(0).String(), nil
}

func (t *Erc20Token) BuildTransferTx(privateKey string, transaction *Transaction) (*base.OptionalString, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	return t.buildTransfer(privateKeyECDSA, transaction)
}

func (t *Erc20Token) BuildTransferTxWithAccount(account *Account, transaction *Transaction) (*base.OptionalString, error) {
	return t.buildTransfer(account.privateKeyECDSA, transaction)
}

func (t *Erc20Token) buildTransfer(privateKey *ecdsa.PrivateKey, transaction *Transaction) (*base.OptionalString, error) {
	err := transaction.TransformToErc20Transaction(t.ContractAddress)
	if err != nil {
		return nil, err
	}
	return t.Token.buildTransfer(privateKey, transaction)
}
