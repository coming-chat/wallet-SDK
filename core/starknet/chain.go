package starknet

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
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
	// NetworkGoerli2 = int(utils.GOERLI2)
)

var (
	MAX_FEE = caigo.MAX_FEE.String()

	erc20TransferSelectorHash = types.BigToHex(types.GetSelectorFromName("transfer"))
)

type Chain struct {
	gw      *gateway.Gateway
	network utils.Network
}

func NewChainWithRpc(baseRpc string, network int) (*Chain, error) {
	var chainIdOpt gateway.Option
	switch network {
	case NetworkMainnet:
		chainIdOpt = gateway.WithChain(gateway.MAINNET_ID)
	case NetworkGoerli:
		chainIdOpt = gateway.WithChain(gateway.GOERLI_ID)
	default:
		return nil, errors.New("invalid starknet network")
	}
	gw := gateway.NewClient(gateway.WithBaseURL(baseRpc), chainIdOpt)
	return &Chain{
		gw:      gw,
		network: utils.Network(network),
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
func (c *Chain) FetchTransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	ctx := context.Background()
	opt := makeTransactionOpt(hash)
	txn, err := c.gw.Transaction(ctx, opt)
	if err != nil {
		return nil, err
	}
	calldata := txn.Transaction.Calldata
	if len(calldata) < 9 && calldata[2] != erc20TransferSelectorHash {
		return nil, base.ErrNotCoinTransferTxn
	}

	receiver := calldata[6]
	amountHex := calldata[7]
	amountInt := hexToBigInt(amountHex)
	detail = &base.TransactionDetail{
		HashString: txn.Transaction.TransactionHash,

		FromAddress: txn.Transaction.SenderAddress,
		ToAddress:   receiver,

		Amount: amountInt.String(),
		Status: mapTransactionStatus(txn.Status),
		// EstimateFees string
		// FinishTimestamp int64
		// FailureMessage string
	}

	switch detail.Status {
	case base.TransactionStatusFailure:
		detail.FailureMessage, _ = c.fetchTransactionFailureMessage(detail.HashString)
	case base.TransactionStatusSuccess:
		detail.FinishTimestamp, _ = c.fetchBlockTimestamp(txn.BlockHash)
		detail.EstimateFees, _ = c.fetchTransactionFee(detail.HashString)
	}
	return detail, nil
}

func (c *Chain) fetchBlockTimestamp(blockHash string) (time int64, err error) {
	block, err := c.gw.Block(context.Background(), &gateway.BlockOptions{
		BlockHash: blockHash,
	})
	if err != nil {
		return 0, err
	}
	return int64(block.Timestamp), nil
}

func (c *Chain) fetchTransactionFee(hash string) (string, error) {
	receipt, err := c.gw.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return "", err
	}
	for i := len(receipt.Events) - 1; i >= 0; i-- {
		event := receipt.Events[i]
		data, err := json.Marshal(event)
		if err != nil {
			return "", err
		}
		var ee gateway.Event
		err = json.Unmarshal(data, &ee)
		if err != nil {
			return "", err
		}
		transferSelectorHash := types.BigToHex(types.GetSelectorFromName("Transfer"))
		if len(ee.Keys) > 0 && ee.Keys[0].String() == transferSelectorHash && len(ee.Data) > 2 {
			return ee.Data[2].Int.String(), nil
		}
	}
	return "", nil
}

func (c *Chain) fetchTransactionFailureMessage(hash string) (string, error) {
	status, err := c.gw.TransactionStatus(context.Background(), gateway.TransactionStatusOptions{
		TransactionHash: hash,
	})
	if err != nil {
		return "", err
	}
	return status.TxFailureReason.ErrorMessage, nil
}

// Fetch transaction status through transaction hash
func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	opt := makeTransactionOpt(hash)
	status, err := c.gw.TransactionStatus(context.Background(), gateway.TransactionStatusOptions{
		TransactionHash: opt.TransactionHash,
		TransactionId:   opt.TransactionId,
	})
	if err != nil {
		return base.TransactionStatusFailure
	}
	return mapTransactionStatus(status.TxStatus)
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
		if err != nil {
			return nil, err
		}
		b := hexToBigInt(string(res.OverallFee))
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

func makeTransactionOpt(hashOrId string) gateway.TransactionOptions {
	if strings.HasPrefix(hashOrId, "0x") {
		return gateway.TransactionOptions{
			TransactionHash: hashOrId,
		}
	} else {
		id, err := strconv.ParseUint(hashOrId, 10, 64)
		if err != nil {
			return gateway.TransactionOptions{}
		}
		return gateway.TransactionOptions{
			TransactionId: id,
		}
	}
}

func mapTransactionStatus(status string) base.TransactionStatus {
	switch types.TransactionState(status) {
	case types.TransactionNotReceived:
		return base.TransactionStatusNone

	case types.TransactionReceived,
		types.TransactionPending:
		return base.TransactionStatusPending

	case types.TransactionAcceptedOnL1,
		types.TransactionAcceptedOnL2:
		return base.TransactionStatusSuccess

	case types.TransactionRejected:
		return base.TransactionStatusFailure
	}
	return base.TransactionStatusNone
}

// return zero if number is invalid
func hexToBigInt(hexNumber string) big.Int {
	hexNumber = strings.TrimPrefix(hexNumber, "0x")
	if res, ok := big.NewInt(0).SetString(hexNumber, 16); ok {
		return *res
	} else {
		return big.Int{}
	}
}
