package starknet

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/NethermindEth/juno/utils"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
)

const (
	BaseRpcUrlMainnet = gateway.MAINNET_BASE
	BaseRpcUrlGoerli  = gateway.GOERLI_BASE

	NetworkMainnet = int(utils.MAINNET)
	NetworkGoerli  = int(utils.GOERLI)
	NetworkGoerli2 = int(utils.GOERLI2)
)

var (
	MAX_FEE = caigo.MAX_FEE.String()
)

type Chain struct {
	gw      *gateway.Gateway
	network utils.Network
}

func NewChainWithRpc(baseRpc string, network int) *Chain {
	gw := gateway.NewClient(gateway.WithBaseURL(baseRpc))
	return &Chain{
		gw:      gw,
		network: utils.Network(network),
	}
}

// MARK - Implement the protocol Chain

func (c *Chain) NewToken(tokenAddress string) (*Token, error) {
	return NewToken(c, tokenAddress)
}

func (c *Chain) MainToken() base.Token {
	t, _ := NewToken(c, ETHTokenAddress)
	return t
}

func (c *Chain) BalanceOfAddress(address string) (b *base.Balance, err error) {
	return c.BalanceOf(address, ETHTokenAddress)
}
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	address, err := encodePublicKeyToAddressArgentX(publicKey)
	if err != nil {
		return nil, err
	}
	return c.BalanceOf(address, ETHTokenAddress)
}
func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOf(account.Address(), ETHTokenAddress)
}

func (c *Chain) BalanceOf(ownerAddress, erc20Address string) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	res, err := c.gw.Call(context.Background(), types.FunctionCall{
		ContractAddress:    types.HexToHash(erc20Address),
		EntryPointSelector: "balanceOf",
		Calldata: []string{
			types.HexToBN(ownerAddress).String(),
		},
	}, "")
	if err != nil {
		return
	}
	low := types.StrToFelt(res[0])
	hi := types.StrToFelt(res[1])
	if low == nil || hi == nil {
		return nil, errors.New("balance response error")
	}

	balance, err := types.NewUint256(low, hi)
	if err != nil {
		return
	}
	return &base.Balance{
		Total:  balance.String(),
		Usable: balance.String(),
	}, nil
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (string, error) {
	return "", base.ErrUnsupportedFunction
}

func (c *Chain) SendSignedTransaction(signedTxn base.SignedTransaction) (hash *base.OptionalString, err_ error) {
	defer base.CatchPanicAndMapToBasicError(&err_)
	txn := signedTxn.(*SignedTransaction)
	if txn == nil {
		return nil, base.ErrInvalidTransactionType
	}

	if txn.depolyTxn != nil {
		resp, err := c.gw.DeployAccount(context.Background(), *txn.depolyTxn)
		if err != nil {
			return nil, err
		}
		return &base.OptionalString{Value: resp.TransactionHash}, nil
	}
	if txn.invokeTxn != nil && txn.Account != nil {
		caigoAccount, err := caigoAccount(c, txn.Account)
		if err != nil {
			return nil, err
		}
		resp, err := caigoAccount.Execute(context.Background(), txn.invokeTxn.calls, txn.invokeTxn.details)
		if err != nil {
			return nil, err
		}
		return &base.OptionalString{Value: resp.TransactionHash}, nil
	}
	return nil, base.ErrMissingTransaction
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	return nil, nil
}

// Fetch transaction status through transaction hash
func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	return base.TransactionStatusFailure
}

// Batch fetch the transaction status, the hash list and the return value,
// which can only be passed as strings separated by ","
// @param hashListString The hash of the transactions to be queried in batches, a string concatenated with ",": "hash1,hash2,hash3"
// @return Batch transaction status, its order is consistent with hashListString: "status1,status2,status3"
func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	return ""
}

// unsupported
func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

// unsupported
func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction base.Transaction, pubkey string) (fee *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

func (c *Chain) EstimateTransactionFeeUseAccount(transaction base.Transaction, acc *Account) (fee *base.OptionalString, err_ error) {
	defer base.CatchPanicAndMapToBasicError(&err_)

	if txn, ok := transaction.(*DeployAccountTransaction); ok {
		fee := txn.txn.MaxFee.BigInt(big.NewInt(0))
		return &base.OptionalString{Value: fee.String()}, nil
	}
	if txn, ok := transaction.(*Transaction); ok {
		caigoAcc, err := caigoAccount(c, acc)
		if err != nil {
			return nil, err
		}
		res, err := caigoAcc.EstimateFee(context.Background(), txn.calls, txn.details)
		b, _ := big.NewInt(0).SetString(strings.TrimLeft(string(res.OverallFee), "0x"), 16)
		return &base.OptionalString{Value: b.String()}, nil
	}
	return nil, base.ErrInvalidTransactionType
}

func (c *Chain) BuildDeployAccountTransaction(publicKey string) (*DeployAccountTransaction, error) {
	txn, err := deployAccountTxnForArgentX(publicKey)
	if err != nil {
		return nil, err
	}

	return &DeployAccountTransaction{
		txn:     txn,
		network: c.network,
	}, nil
}

func caigoAccount(chain *Chain, acc *Account) (*caigo.Account, error) {
	return caigo.NewGatewayAccount(acc.privateKey.String(), acc.Address(), &gateway.GatewayProvider{Gateway: *chain.gw}, caigo.AccountVersion1)
}
