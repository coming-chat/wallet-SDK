package sui

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/coming-chat/go-sui/client"
	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

const (
	MaxGasBudget   = 10000
	MaxGasForMerge = 10000

	MaxGasForPay      = 150
	MaxGasForTransfer = 200

	FaucetUrlTestnet = "https://faucet.testnet.sui.io/gas"
)

type Chain struct {
	rpcClient *client.Client
	RpcUrl    string
}

func NewChainWithRpcUrl(rpcUrl string) *Chain {
	return &Chain{RpcUrl: rpcUrl}
}

func (c *Chain) client() (*client.Client, error) {
	if c.rpcClient != nil {
		return c.rpcClient, nil
	}
	var err error
	c.rpcClient, err = client.Dial(c.RpcUrl)
	return c.rpcClient, err
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return NewTokenMain(c)
}

func (c *Chain) BalanceOfAddress(address string) (b *base.Balance, err error) {
	return c.MainToken().BalanceOfAddress(address)
}
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return c.MainToken().BalanceOfPublicKey(publicKey)
}
func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (hash string, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	bytes, err := types.NewBase64Data(signedTx)
	if err != nil {
		return
	}
	signedTxn := types.SignedTransaction{}
	err = json.Unmarshal(bytes.Data(), &signedTxn)
	if err != nil {
		return
	}
	cli, err := c.client()
	if err != nil {
		return
	}
	response, err := cli.ExecuteTransaction(context.Background(), signedTxn, types.TxnRequestTypeWaitForTxCert)
	if err != nil {
		return
	}

	hash = response.TransactionDigest().String()
	return
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	digest, err := types.NewBase64Data(hash)
	if err != nil {
		return
	}
	cli, err := c.client()
	if err != nil {
		return
	}
	resp, err := cli.GetTransaction(context.Background(), *digest)
	if err != nil {
		return nil, err
	}

	var firstRecipient *types.HexData
	var total uint64
	for _, txn := range resp.Certificate.Data.Transactions {
		if tsui := txn.TransferSui; tsui != nil {
			if firstRecipient == nil {
				firstRecipient = &tsui.Recipient
				total = tsui.Amount
			} else if bytes.Compare(firstRecipient.Data(), tsui.Recipient.Data()) == 0 {
				total = total + tsui.Amount
			}
		} else if tobject := txn.TransferObject; tobject != nil {
			if firstRecipient == nil {
				firstRecipient = &tobject.Recipient
			}
		}
	}
	if firstRecipient == nil {
		return nil, errors.New("Invalid coin transfer transaction.")
	}

	detail = &base.TransactionDetail{
		HashString:      hash,
		FromAddress:     resp.Certificate.Data.Sender.String(),
		ToAddress:       firstRecipient.String(),
		Amount:          strconv.FormatUint(total, 10),
		EstimateFees:    strconv.FormatUint(resp.Effects.GasFee(), 10),
		FinishTimestamp: int64(resp.TimestampMs / 1000),
	}
	status := resp.Effects.Status
	if status.Status == types.TransactionStatusSuccess {
		detail.Status = base.TransactionStatusSuccess
	} else {
		detail.Status = base.TransactionStatusFailure
		detail.FailureMessage = status.Error
	}
	return
}

func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	detail, err := c.FetchTransactionDetail(hash)
	if err != nil {
		return base.TransactionStatusNone
	}
	return detail.Status
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

// @param gasId gas object to be used in this transaction, the gateway will pick one from the signer's possession if not provided
func (c *Chain) TransferObject(sender, receiver, objectId, gasId string, gasBudget int64) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	senderAddress, err := types.NewAddressFromHex(sender)
	if err != nil {
		return
	}
	receiverAddress, err := types.NewAddressFromHex(receiver)
	if err != nil {
		return
	}
	nftObject, err := types.NewHexData(objectId)
	if err != nil {
		return nil, err
	}
	var gas *types.ObjectId = nil
	if gasId != "" {
		gas, err = types.NewHexData(gasId)
		if err != nil {
			return nil, errors.New("Invalid gas object id")
		}
	}
	client, err := c.client()
	if err != nil {
		return
	}
	tx, err := client.TransferObject(context.Background(), *senderAddress, *receiverAddress, *nftObject, gas, uint64(gasBudget))
	if err != nil {
		return
	}
	return &Transaction{Txn: *tx}, nil
}

/**
 * @param address Hex-encoded 16 bytes Sui account address wich mints tokens
 * @param faucetUrl default https://faucet.testnet.sui.io/gas
 * @return digest of transfer transaction.
 */
func FaucetFundAccount(address string, faucetUrl string) (h *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	_, err = types.NewAddressFromHex(address)
	if err != nil {
		return
	}
	if faucetUrl == "" {
		faucetUrl = FaucetUrlTestnet
	}

	var authority = ""
	if strings.Contains(faucetUrl, "devnet") {
		authority = "faucet.devnet.sui.io"
	} else {
		authority = "faucet.testnet.sui.io"
	}
	paramJson := fmt.Sprintf(`{"FixedAmountRequest":{"recipient":"%v"}}`, address)
	params := httpUtil.RequestParams{
		Header: map[string]string{
			"Content-Type": "application/json",
			"Authority":    authority,
			"authority":    authority,
		},
		Body: []byte(paramJson),
	}
	res, err := httpUtil.Post(faucetUrl, params)
	if err != nil {
		return
	}
	response := struct {
		TransferredObjects []struct {
			Amount uint64         `json:"amount"`
			Id     types.ObjectId `json:"id"`
			Digest types.Digest   `json:"transfer_tx_digest"`
		} `json:"transferred_gas_objects,omitempty"`
		Error string `json:"error,omitempty"`
	}{}
	err = json.Unmarshal(res, &response)
	if err != nil {
		return
	}
	if strings.TrimSpace(response.Error) != "" {
		return nil, errors.New(response.Error)
	}
	if len(response.TransferredObjects) <= 0 {
		return nil, errors.New("Transaction not found.")
	}

	digest := response.TransferredObjects[0].Digest.String()
	return &base.OptionalString{Value: digest}, nil
}
