package solana

import (
	"context"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/rpc"

	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/portto/solana-go-sdk/types"
)

const (
	DevnetRPCEndpoint  = rpc.DevnetRPCEndpoint
	TestnetRPCEndpoint = rpc.TestnetRPCEndpoint
	MainnetRPCEndpoint = rpc.MainnetRPCEndpoint
)

type Chain struct {
	*Util
	RpcUrl string
}

func NewChainWithRpc(rpcUrl string) *Chain {
	util := NewUtil()
	return &Chain{
		Util:   util,
		RpcUrl: rpcUrl,
	}
}

func (c *Chain) client() *client.Client {
	return client.NewClient(c.RpcUrl)
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return &Token{chain: c}
}

func (c *Chain) BalanceOfAddress(address string) (*base.Balance, error) {
	client := c.client()
	balance, err := client.GetBalance(context.Background(), address)
	if err != nil {
		return nil, err
	}
	b := strconv.FormatUint(balance, 10)
	return &base.Balance{Total: b, Usable: b}, nil
}
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	address, err := EncodePublicKeyToAddress(publicKey)
	if err != nil {
		return nil, err
	}
	return c.BalanceOfAddress(address)
}
func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (string, error) {
	client := c.client()
	bytes, err := hexTypes.HexDecodeString(signedTx)
	if err != nil {
		return "", err
	}
	transaction, err := types.TransactionDeserialize(bytes)
	if err != nil {
		return "", err
	}
	res, err := client.SendTransaction(context.Background(), transaction)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (c *Chain) SendSignedTransaction(signedTxn base.SignedTransaction) (hash *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	response, err := c.client().GetTransaction(context.Background(), hash)
	if err != nil {
		return nil, err
	}
	detail := &base.TransactionDetail{HashString: hash}
	err = decodeTransaction(response, detail)
	if err != nil {
		return nil, err
	}
	return detail, nil
}

func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	response, err := c.client().GetTransaction(context.Background(), hash)
	if err != nil {
		return base.TransactionStatusNone
	}
	if response == nil || response.Meta == nil {
		return base.TransactionStatusPending
	}
	if response.Meta.Err == nil {
		return base.TransactionStatusSuccess
	}
	return base.TransactionStatusFailure
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}
func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction base.Transaction, pubkey string) (fee *base.OptionalString, err error) {
	return c.EstimateTransactionFee(transaction)
}

func decodeTransaction(tx *client.GetTransactionResponse, to *base.TransactionDetail) error {
	base.CatchPanicAndMapToBasicError(nil)

	if tx == nil || tx.BlockTime == nil {
		to.Status = base.TransactionStatusPending
		return nil
	}

	message := tx.Transaction.Message
	for _, instruction := range message.Instructions {
		program := message.Accounts[instruction.ProgramIDIndex]
		if program != common.SystemProgramID {
			continue
		}

		// We only support decode amount transfer currently.
		data := instruction.Data
		instruct := binary.LittleEndian.Uint32(data[:4])
		toidx := -1
		switch sysprog.Instruction(instruct) {
		case sysprog.InstructionTransfer:
			toidx = instruction.Accounts[1]
		case sysprog.InstructionTransferWithSeed:
			toidx = instruction.Accounts[2]
		}
		if toidx == -1 {
			continue
		}

		fromidx := instruction.Accounts[0]
		to.FromAddress = message.Accounts[fromidx].ToBase58()
		to.ToAddress = message.Accounts[toidx].ToBase58()
		amount := binary.LittleEndian.Uint64(data[4:12])
		to.Amount = strconv.FormatUint(amount, 10)
		to.EstimateFees = strconv.FormatUint(tx.Meta.Fee, 10)
		to.FinishTimestamp = *tx.BlockTime

		if tx.Meta.Err == nil {
			to.Status = base.TransactionStatusSuccess
		} else {
			to.Status = base.TransactionStatusFailure
			to.FailureMessage = tx.Meta.LogMessages[len(tx.Meta.LogMessages)-1]
		}

		return nil
	}
	return errors.New("The transaction does not contain an amount transfer")
}
