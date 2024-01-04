package starknet

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/coming-chat/wallet-SDK/core/base"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/xiang-xx/starknet.go/account"
	"github.com/xiang-xx/starknet.go/rpc"
	"github.com/xiang-xx/starknet.go/utils"
)

const (
	BaseRpcUrlMainnet = "https://starknet-mainnet.public.blastapi.io"

	InvokeMaxFee = 2e14 // 0.0002 ETH
)

var (
	erc20TransferSelectorFelt = utils.GetSelectorFromNameFelt("transfer")
	erc20TransferSelectorHash = erc20TransferSelectorFelt.String()
)

type Chain struct {
	rpc *account.Account
}

func NewChainWithRpc(baseRpc string) (*Chain, error) {
	cli, err := ethrpc.DialContext(context.Background(), baseRpc)
	if err != nil {
		return nil, err
	}
	provider := rpc.NewProvider(cli)
	rpc, err := account.NewAccount(provider, nil, "", nil, 0)
	if err != nil {
		return nil, err
	}
	return &Chain{
		rpc: rpc,
	}, nil
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
	address, err := EncodePublicKeyToAddress(publicKey)
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

	ownerFelt, err := utils.HexToFelt(ownerAddress)
	if err != nil {
		return
	}
	erc20Felt, err := utils.HexToFelt(erc20Address)
	if err != nil {
		return
	}
	tx := rpc.FunctionCall{
		ContractAddress:    erc20Felt,
		EntryPointSelector: utils.GetSelectorFromNameFelt("balanceOf"),
		Calldata:           []*felt.Felt{ownerFelt},
	}

	resp, err := c.rpc.Call(context.Background(), tx, latestBlockId)
	if err != nil {
		return
	}
	if len(resp) != 2 {
		return nil, errors.New("balance response error")
	}
	low := utils.FeltToBigInt(resp[0])
	high := utils.FeltToBigInt(resp[1])
	balance := new(big.Int).Add(new(big.Int).Lsh(high, 128), low)
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

func (c *Chain) SendSignedTransaction(signedTxn base.SignedTransaction) (hash *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	txn := AsSignedTransaction(signedTxn)
	if txn == nil {
		return nil, base.ErrInvalidTransactionType
	}

	if txn.depolyTxn != nil {
		resp, err := c.rpc.AddDeployAccountTransaction(context.Background(), rpc.BroadcastDeployAccountTxn{
			DeployAccountTxn: *txn.depolyTxn,
		})
		if err != nil {
			return nil, err
		}
		return &base.OptionalString{Value: resp.TransactionHash.String()}, nil
	}
	if txn.invokeTxn != nil {
		resp, err := c.rpc.AddInvokeTransaction(context.Background(), txn.invokeTxn)
		if err != nil {
			if !txn.NeedAutoDeploy {
				return nil, err
			}
			deployed, err_ := c.IsContractAddressDeployed(txn.invokeTxn.SenderAddress.String())
			if err_ != nil || deployed.Value == true {
				return nil, err // if query failed, return the previous error.
			}

			// we need deploy the account firstly, and resend the original txn with fixed Nonce 1
			pubHex := txn.Account.PublicKeyHex()
			pubFelt, err := utils.HexToFelt(pubHex)
			if err != nil {
				return nil, err
			}
			version, err := CheckCairoVersionFelt(txn.invokeTxn.SenderAddress, pubFelt)
			if err != nil {
				return nil, err
			}
			deployTxn, err := c.deployAccountTxnWithVersion(pubHex, "", version == 0)
			if err != nil {
				return nil, err
			}
			signedDeployTxn, err := deployTxn.SignedTransactionWithAccount(txn.Account)
			if err != nil {
				return nil, err
			}
			_, err = c.SendSignedTransaction(signedDeployTxn)
			if err != nil {
				return nil, err
			}
			// now resend the original txn, the nonce is 1 now
			txn.invokeTxn.Nonce = new(felt.Felt).SetUint64(1)
			resp, err = c.rpc.AddInvokeTransaction(context.Background(), txn.invokeTxn)
			if err != nil {
				return nil, err
			}
		}
		return &base.OptionalString{Value: resp.TransactionHash.String()}, nil
	}
	return nil, base.ErrMissingTransaction
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	hashFelt, err := utils.HexToFelt(hash)
	if err != nil {
		return nil, base.ErrInvalidTransactionHash
	}

	ctx := context.Background()
	rawTxn, err := c.rpc.TransactionByHash(ctx, hashFelt)
	if err != nil {
		return nil, err
	}
	txn, ok := rawTxn.(rpc.InvokeTxnV1)
	if !ok {
		return nil, base.ErrNotCoinTransferTxn
	}
	calldata := txn.Calldata
	dataLen := len(calldata)
	if dataLen < 7 || calldata[2].Cmp(erc20TransferSelectorFelt) != 0 {
		return nil, base.ErrNotCoinTransferTxn
	}
	detail = &base.TransactionDetail{
		HashString:   hash,
		FromAddress:  txn.SenderAddress.String(),
		ToAddress:    calldata[dataLen-3].String(),
		Amount:       calldata[dataLen-2].Text(10),
		EstimateFees: txn.MaxFee.Text(10),
		Status:       base.TransactionStatusPending,
	}

	receipt, err := c.rpc.TransactionReceipt(ctx, hashFelt)
	if err != nil {
		return detail, nil
	}
	invokeReceipt, ok := receipt.(rpc.InvokeTransactionReceipt)
	if !ok {
		return detail, nil
	}
	if invokeReceipt.ActualFee.Amount != nil {
		detail.EstimateFees = invokeReceipt.ActualFee.Amount.Text(10)
	}
	if invokeReceipt.ExecutionStatus == rpc.TxnExecutionStatusREVERTED {
		detail.Status = base.TransactionStatusFailure
		detail.FailureMessage = invokeReceipt.RevertReason
		return detail, nil
	}
	blockInfo, err := c.rpc.BlockWithTxHashes(ctx, rpc.WithBlockHash(invokeReceipt.BlockHash))
	if err != nil {
		detail.Status = base.TransactionStatusPending
		return detail, nil
	}
	detail.Status = base.TransactionStatusSuccess
	detail.FinishTimestamp = int64(blockInfo.(*rpc.BlockTxHashes).BlockHeader.Timestamp)
	return detail, nil
}

// Fetch transaction status through transaction hash
func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	hashFelt, err := utils.HexToFelt(hash)
	if err != nil {
		return base.TransactionStatusNone
	}
	receipt, err := c.rpc.TransactionReceipt(context.Background(), hashFelt)
	if err != nil {
		return base.TransactionStatusNone
	}
	if receipt.GetExecutionStatus() == rpc.TxnExecutionStatusREVERTED {
		return base.TransactionStatusFailure
	} else {
		return base.TransactionStatusSuccess
	}
}

// Batch fetch the transaction status, the hash list and the return value,
// which can only be passed as strings separated by ","
// @param hashListString The hash of the transactions to be queried in batches, a string concatenated with ",": "hash1,hash2,hash3"
// @return Batch transaction status, its order is consistent with hashListString: "status1,status2,status3"
func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

// unsupported
func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

// unsupported
func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction base.Transaction, pubkey string) (fee *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

func (c *Chain) EstimateTransactionFeeUseAccount(transaction base.Transaction, acc *Account) (fee *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if txn, ok := transaction.(*DeployAccountTransaction); ok {
		return base.NewOptionalString(txn.MaxFee.Text(10)), nil
	}
	if txn, ok := transaction.(*Transaction); ok {
		txn.txnV1.Signature, err = acc.SignHash(txn.txnHash)
		if err != nil {
			return nil, err
		}
		defer func() {
			txn.txnV1.Signature = nil // clean when finished
		}()

		resp, err_ := c.rpc.SimulateTransactions(context.Background(), latestBlockId,
			[]rpc.Transaction{txn.txnV1}, []rpc.SimulationFlag{rpc.SKIP_FEE_CHARGE})
		if err_ != nil {
			return nil, err_
		}
		if len(resp) <= 0 {
			return nil, errors.New("estimate fee failed")
		}
		fee := utils.FeltToBigInt(resp[0].OverallFee)
		result := base.BigIntMultiply(fee, 1.2)
		txn.txnV1.MaxFee = utils.BigIntToFelt(result)             // update maxfee
		txn.txnHash, err = c.rpc.TransactionHashInvoke(txn.txnV1) // update txn hash
		if err != nil {
			return nil, err
		}
		return base.NewOptionalString(result.String()), nil
	}
	return nil, base.ErrInvalidTransactionType
}

// BuildDeployAccountTransaction
// @param maxFee default is 0.0002
func (c *Chain) BuildDeployAccountTransaction(publicKey string, maxFee string) (*DeployAccountTransaction, error) {
	return c.deployAccountTxnWithVersion(publicKey, maxFee, false)
}

// BuildDeployAccountTransaction
// @param maxFee default is 0.0002
func (c *Chain) BuildDeployAccountTransactionCairo0(publicKey string, maxFee string) (*DeployAccountTransaction, error) {
	return c.deployAccountTxnWithVersion(publicKey, maxFee, true)
}

func (c *Chain) deployAccountTxnWithVersion(publicKey string, maxFee string, isCairo0 bool) (*DeployAccountTransaction, error) {
	var feeInt *big.Int
	var ok bool
	if maxFee == "" {
		feeInt = big.NewInt(0).SetUint64(2e14 + random(1e12)) // 0.000200xx
	} else {
		if feeInt, ok = big.NewInt(0).SetString(maxFee, 10); !ok {
			return nil, base.ErrInvalidAmount
		}
	}
	return NewDeployAccountTransaction(publicKey, feeInt, c.rpc, isCairo0)
}

func (c *Chain) IsContractAddressDeployed(contractAddress string) (b *base.OptionalBool, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	addressFelt, err := utils.HexToFelt(contractAddress)
	if err != nil {
		return nil, base.ErrInvalidAccountAddress
	}
	nonce, err := c.rpc.Nonce(context.Background(), latestBlockId, addressFelt)
	if err != nil {
		if err.Error() == rpc.ErrContractNotFound.Error() {
			return base.NewOptionalBool(false), nil
		}
		return nil, err
	}
	deployed := nonce != nil
	return base.NewOptionalBool(deployed), nil
}
