package starknet

import (
	"context"
	"errors"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
)

const (
	BaseRpcUrlMainnet = gateway.MAINNET_BASE
	BaseRpcUrlGoerli  = gateway.GOERLI_BASE
)

var (
	MAX_FEE = caigo.MAX_FEE.String()
)

type Chain struct {
	gw *gateway.Gateway
}

func NewChainWithRpc(baseRpc string) *Chain {
	gw := gateway.NewClient(gateway.WithBaseURL(baseRpc))
	return &Chain{gw: gw}
}

// MARK - Implement the protocol Chain

func (c *Chain) NewToken(tokenContractAddress string) *Token {
	return NewToken(c, tokenContractAddress)
}

// warning: please use `chain.NewToken(contractAddress)` instead.
func (c *Chain) MainToken() base.Token {
	return nil
}

func (c *Chain) BalanceOfAddress(address string) (b *base.Balance, err error) {
	return nil, errors.New("please call with `chain.NewToken(contractAddress).BalanceOf...`")
}
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return c.BalanceOfAddress("")
}
func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfAddress("")
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
	return "", nil
}

func (c *Chain) SendSignedTransaction(signedTxn base.SignedTransaction) (hash *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
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

// Most chains can estimate the fee directly to the transaction object
// **But two chains don't work: `aptos`, `starcoin`**
func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	return nil, nil
}

// All chains can call this method to estimate the gas fee.
// **Chain  `aptos`, `starcoin` must pass in publickey**
func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction base.Transaction, pubkey string) (fee *base.OptionalString, err error) {
	return nil, nil
}

func (c *Chain) DeployAccount(contractAddress string) {

}
