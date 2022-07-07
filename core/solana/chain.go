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

	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/portto/solana-go-sdk/types"
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
		return base.TransactionStatusFailure
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

func decodeTransaction(tx *client.GetTransactionResponse, to *base.TransactionDetail) error {
	base.CatchPanicAndMapToBasicError(nil)

	message := tx.Transaction.Message
	for _, instruction := range message.Instructions {
		program := message.Accounts[instruction.ProgramIDIndex]
		if program != common.SystemProgramID {
			continue
		}

		// We only support decode amount transfer currently.
		data := instruction.Data
		instruct := binary.LittleEndian.Uint32(data[:4])
		isTransfer := false
		switch sysprog.Instruction(instruct) {
		case sysprog.InstructionTransfer:
			isTransfer = true
			to.ToAddress = message.Accounts[1].ToBase58()
		case sysprog.InstructionTransferWithSeed:
			isTransfer = true
			to.ToAddress = message.Accounts[2].ToBase58()
		}
		if !isTransfer {
			continue
		}

		to.FromAddress = message.Accounts[0].ToBase58()
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
