package eth

import (
	"encoding/json"
	"errors"
	"math/big"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

type TokenProtocol interface {
	base.Token

	// need `fromAddress`, `receiverAddress`, `gasPrice`, `gasLimit`, `amount`
	EstimateGasFeeLayer2(msg *CallMsg) (*OptimismLayer2Gas, error)
	EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string) (string, error)
	BuildTransferTx(privateKey string, transaction *Transaction) (*base.OptionalString, error)
	BuildTransferTxWithAccount(account *Account, transaction *Transaction) (*base.OptionalString, error)
}

type Token struct {
	chain *Chain
}

func NewToken(chain *Chain) *Token {
	return &Token{chain}
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
	return t.chain.BalanceOfAddress(address)
}

func (t *Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return t.BalanceOfAddress(publicKey)
}

func (t *Token) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.BalanceOfAddress(account.Address())
}

// MARK - Eth TokenProtocol

func (t *Token) EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string) (string, error) {
	msg := NewCallMsg()
	msg.SetFrom(fromAddress)
	msg.SetTo(receiverAddress)
	msg.SetGasPrice(gasPrice)
	msg.SetValue(amount)

	res, err := t.chain.EstimateGasLimit(msg)
	if err != nil {
		return "", err
	}
	return res.Value, nil
}

func (t *Token) BuildTransferTx(privateKey string, transaction *Transaction) (*base.OptionalString, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	return t.chain.buildTransfer(privateKeyECDSA, transaction)
}

func (t *Token) BuildTransferTxWithAccount(account *Account, transaction *Transaction) (*base.OptionalString, error) {
	return t.chain.buildTransfer(account.privateKeyECDSA, transaction)
}

func (t *Token) BuildTransfer(sender, receiver, amount string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}
func (t *Token) CanTransferAll() bool {
	return false
}
func (t *Token) BuildTransferAll(sender, receiver string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}

type OptimismLayer2Gas struct {
	L1GasLimit string
	L1GasPrice string
	L2GasLimit string
	L2GasPrice string
}

// l1GasLimit * l1GasPrice + l2Gaslimit * l2GasPrice
func (g *OptimismLayer2Gas) GasFee() string {
	l1Limit, ok := big.NewInt(0).SetString(g.L1GasLimit, 10)
	if !ok {
		l1Limit = big.NewInt(0)
	}
	l1Price, ok := big.NewInt(0).SetString(g.L1GasPrice, 10)
	if !ok {
		l1Price = big.NewInt(0)
	}
	l2Limit, ok := big.NewInt(0).SetString(g.L2GasLimit, 10)
	if !ok {
		l2Limit = big.NewInt(0)
	}
	l2Price, ok := big.NewInt(0).SetString(g.L2GasPrice, 10)
	if !ok {
		l2Price = big.NewInt(0)
	}
	l1Fee := big.NewInt(0).Mul(l1Limit, l1Price)
	l2Fee := big.NewInt(0).Mul(l2Limit, l2Price)
	return big.NewInt(0).Add(l1Fee, l2Fee).String()
}

func (t *Token) EstimateGasFeeLayer2(msg *CallMsg) (*OptimismLayer2Gas, error) {
	// We need fetch the ethereum mainnet Gas Price
	ethMainRpc := "https://geth-mainnet.coming.chat"
	l1GasPriceString, err := NewChainWithRpc(ethMainRpc).SuggestGasPrice()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(msg.msg)
	if err != nil {
		return nil, err
	}
	l1GasLimit := calculateL1GasLimit(data, overhead)

	return &OptimismLayer2Gas{
		L1GasPrice: l1GasPriceString.Value,
		L1GasLimit: l1GasLimit.String(),
		L2GasPrice: msg.GetGasPrice(),
		L2GasLimit: msg.GetGasLimit(),
	}, nil
}

const overhead uint64 = 200 * params.TxDataNonZeroGasEIP2028

func calculateL1GasLimit(data []byte, overhead uint64) *big.Int {
	zeroes, ones := zeroesAndOnes(data)
	zeroesCost := zeroes * params.TxDataZeroGas
	onesCost := ones * params.TxDataNonZeroGasEIP2028
	gasLimit := zeroesCost + onesCost + overhead
	return new(big.Int).SetUint64(gasLimit)
}

func zeroesAndOnes(data []byte) (uint64, uint64) {
	var zeroes uint64
	var ones uint64
	for _, byt := range data {
		if byt == 0 {
			zeroes++
		} else {
			ones++
		}
	}
	return zeroes, ones
}
