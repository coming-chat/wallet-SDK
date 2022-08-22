package eth

import (
	"crypto/ecdsa"
	"math/big"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type IChain interface {
	base.Chain
	SubmitTransactionData(account base.Account, to string, data []byte, value string) (string, error)
	GetEthChain() (*EthChain, error)
}

type Chain struct {
	RpcUrl string
}

func NewChainWithRpc(rpcUrl string) *Chain {
	return &Chain{
		RpcUrl: rpcUrl,
	}
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return &Token{chain: c}
}

func (c *Chain) MainEthToken() TokenProtocol {
	return &Token{chain: c}
}

func (c *Chain) Erc20Token(contractAddress string) TokenProtocol {
	return &Erc20Token{
		Token:           &Token{chain: c},
		ContractAddress: contractAddress,
	}
}

func (c *Chain) BalanceOfAddress(address string) (*base.Balance, error) {
	b := base.EmptyBalance()

	eip55Address, err := TransformEIP55Address(address)
	if err != nil {
		return b, err
	}

	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return b, err
	}

	balance, err := chain.Balance(eip55Address)
	if err != nil {
		return b, err
	}
	return &base.Balance{
		Total:  balance,
		Usable: balance,
	}, nil
}

func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return c.BalanceOfAddress(publicKey)
}

func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (string, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return "", err
	}
	return chain.SendRawTransaction(signedTx)
}

// Fetch transaction details through transaction hash
// Support normal or erc20 transfer
func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return nil, err
	}
	detail, msg, err := chain.FetchTransactionDetail(hash)
	if err != nil {
		return nil, err
	}
	if data := msg.Data(); len(data) > 0 {
		method, params, err := DecodeContractParams(Erc20AbiStr, data)
		if err == nil && method == ERC20_METHOD_TRANSFER {
			detail.ToAddress = params[0].(common.Address).String()
			detail.Amount = params[1].(*big.Int).String()
		}
	}
	return detail, nil
}

// Fetch transaction status through transaction hash
func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return base.TransactionStatusNone
	}
	return chain.FetchTransactionStatus(hash)
}

// Batch fetch the transaction status, the hash list and the return value,
// which can only be passed as strings separated by ","
// @param hashListString The hash of the transactions to be queried in batches, a string concatenated with ",": "hash1,hash2,hash3"
// @return Batch transaction status, its order is consistent with hashListString: "status1,status2,status3"
func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return ""
	}
	return chain.SdkBatchTransactionStatus(hashListString)
}

func (c *Chain) BuildTransferTx(privateKey string, transaction *Transaction) (*base.OptionalString, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	return c.buildTransfer(privateKeyECDSA, transaction)
}

func (c *Chain) BuildTransferTxWithAccount(account *Account, transaction *Transaction) (*base.OptionalString, error) {
	return c.buildTransfer(account.privateKeyECDSA, transaction)
}

func (c *Chain) buildTransfer(privateKey *ecdsa.PrivateKey, transaction *Transaction) (*base.OptionalString, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return nil, err
	}

	if transaction.Nonce == "" || transaction.Nonce == "0" {
		address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
		nonce, err := chain.Nonce(address)
		if err != nil {
			nonce = "0"
			err = nil
		}
		transaction.Nonce = nonce
	}

	rawTx, err := transaction.GetRawTx()
	if err != nil {
		return nil, err
	}

	txResult, err := chain.buildTxWithTransaction(rawTx, privateKey)
	if err != nil {
		return nil, err
	}

	return &base.OptionalString{Value: txResult.TxHex}, nil
}

// MARK - Implement the protocol IChain

func (c *Chain) SubmitTransactionData(account base.Account, to string, data []byte, value string) (string, error) {
	gasPrice, err := c.SuggestGasPrice()
	if err != nil {
		return "", err
	}
	msg := NewCallMsg()
	msg.SetFrom(account.Address())
	msg.SetTo(to)
	msg.SetGasPrice(gasPrice.Value)
	msg.SetData(data)
	msg.SetValue(value)

	gasLimit, err := c.EstimateGasLimit(msg)
	if err != nil {
		gasLimit = &base.OptionalString{Value: "200000"}
		err = nil
	}
	msg.SetGasLimit(gasLimit.Value)
	tx := msg.TransferToTransaction()
	privateKeyHex, err := account.PrivateKeyHex()
	if err != nil {
		return "", err
	}
	signedTx, err := c.SignTransaction(privateKeyHex, tx)
	if err != nil {
		return "", err
	}
	return c.SendRawTransaction(signedTx.Value)
}

func (c *Chain) GetEthChain() (*EthChain, error) {
	return GetConnection(c.RpcUrl)
}
