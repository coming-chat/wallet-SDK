package aptos

import (
	"context"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
	txbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/lcs"
	"github.com/coming-chat/wallet-SDK/core/base"
)

const (
	MaxGasAmount = 2000
	GasPrice     = 1
	TxExpireSec  = 600
)

type IChain interface {
	base.Chain
	SubmitTransactionPayloadBCS(account base.Account, data []byte) (string, error)
	EstimatePayloadGasFeeBCS(account base.Account, data []byte) (*base.OptionalString, error)
	GetClient() (*aptosclient.RestClient, error)
}

type Chain struct {
	restClient *aptosclient.RestClient
	RestUrl    string
}

func NewChainWithRestUrl(restUrl string) *Chain {
	return &Chain{RestUrl: restUrl}
}

func (c *Chain) client() (*aptosclient.RestClient, error) {
	if c.restClient != nil {
		return c.restClient, nil
	}
	var err error
	c.restClient, err = aptosclient.Dial(context.Background(), c.RestUrl)
	return c.restClient, err
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return &Token{chain: c}
}

func (c *Chain) BalanceOfAddress(address string) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	client, err := c.client()
	if err != nil {
		return
	}

	balance, err := client.AptosBalanceOf(address)
	if err != nil {
		return
	}

	return &base.Balance{
		Total:  balance.String(),
		Usable: balance.String(),
	}, nil
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
func (c *Chain) SendRawTransaction(signedTx string) (hash string, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	bytes, err := types.HexDecodeString(signedTx)
	if err != nil {
		return
	}
	client, err := c.client()
	if err != nil {
		return
	}
	resultTx, err := client.SubmitSignedBCSTransaction(bytes)
	if err != nil {
		return
	}

	return resultTx.Hash, nil
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (detail *base.TransactionDetail, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	client, err := c.client()
	if err != nil {
		return
	}

	transaction, err := client.GetTransactionByHash(hash)
	if err != nil {
		return
	}
	return toBaseTransaction(transaction)
}

func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	client, err := c.client()
	if err != nil {
		return base.TransactionStatusNone
	}
	transaction, err := client.GetTransactionByHash(hash)
	if err != nil {
		return base.TransactionStatusNone
	}
	if transaction.Success {
		return base.TransactionStatusSuccess
	} else {
		return base.TransactionStatusFailure
	}
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

// MARK - Implement the protocol IChain

func (c *Chain) GetClient() (*aptosclient.RestClient, error) {
	return c.client()
}

func (c *Chain) EstimatePayloadGasFeeBCS(account base.Account, data []byte) (*base.OptionalString, error) {
	var (
		err   error
		txAbi *txbuilder.RawTransaction
	)
	payload := txbuilder.TransactionPayloadEntryFunction{}
	if err := lcs.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	if txAbi, err = c.createTransactionFromPayloadBCS(account, payload); err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: strconv.FormatUint(txAbi.MaxGasAmount, 10)}, nil
}

func (c *Chain) EstimateMaxGasAmount(account base.Account, transaction *aptostypes.Transaction) (uint64, error) {
	client, err := c.client()
	if err != nil {
		return 0, err
	}
	transaction.MaxGasAmount = base.Max(transaction.MaxGasAmount, 5000) // as big as possible

	commitedTxs, err := client.SimulateTransaction(transaction, account.PublicKeyHex())
	if err != nil {
		return 0, err
	}
	if len(commitedTxs) <= 0 {
		return 0, errors.New("Query gas fee failed.")
	}

	tx := commitedTxs[0]
	maxGas := (tx.GasUsed*15 + 9) / 10 // ceil(fee * 1.5)
	return maxGas, nil
}

func (c *Chain) EstimateMaxGasAmountBCS(publicKey []byte, rawTxn *txbuilder.RawTransaction) (uint64, error) {
	signedTx, err := txbuilder.GenerateBCSSimulation(publicKey, rawTxn)
	if err != nil {
		return 0, err
	}
	client, err := c.client()
	if err != nil {
		return 0, err
	}
	commitedTxs, err := client.SimulateSignedBCSTransaction(signedTx)
	if err != nil {
		return 0, err
	}
	if len(commitedTxs) <= 0 {
		return 0, errors.New("Query gas fee failed.")
	}

	tx := commitedTxs[0]
	maxGas := (tx.GasUsed*15 + 9) / 10 // ceil(fee * 1.5)
	return maxGas, nil
}

func (c *Chain) SubmitTransactionPayloadBCS(account base.Account, data []byte) (string, error) {
	var (
		err         error
		client      *aptosclient.RestClient
		txAbi       *txbuilder.RawTransaction
		submittedTx *aptostypes.Transaction
		signedTxn   []byte
	)
	if client, err = c.client(); err != nil {
		return "", err
	}
	payload := txbuilder.TransactionPayloadEntryFunction{}
	if err := lcs.Unmarshal(data, &payload); err != nil {
		return "", err
	}
	if txAbi, err = c.createTransactionFromPayloadBCS(account, payload); err != nil {
		return "", err
	}
	aptAccount, ok := account.(*Account)
	if !ok {
		return "", errors.New("invalid account type")
	}
	if signedTxn, err = txbuilder.GenerateBCSTransaction(aptAccount.account, txAbi); err != nil {
		return "", err
	}
	if submittedTx, err = client.SubmitSignedBCSTransaction(signedTxn); err != nil {
		return "", err
	}
	return submittedTx.Hash, err
}

func (c *Chain) signTransaction(account *Account, transaction *txbuilder.RawTransaction) ([]byte, error) {
	return txbuilder.GenerateBCSTransaction(account.account, transaction)
}

// @return The raw transaction that `MaxGasAmount` has obtained from the chain in real time.
func (c *Chain) createTransactionFromPayloadBCS(account base.Account, payload txbuilder.TransactionPayload) (*txbuilder.RawTransaction, error) {
	var (
		err         error
		client      *aptosclient.RestClient
		accountData *aptostypes.AccountCoreData
		ledgerInfo  *aptostypes.LedgerInfo
	)

	if client, err = c.client(); err != nil {
		return nil, err
	}
	if accountData, err = client.GetAccount(account.Address()); err != nil {
		return nil, err
	}
	if ledgerInfo, err = client.LedgerInfo(); err != nil {
		return nil, err
	}
	txAbi := &txbuilder.RawTransaction{
		Sender:                  getAuthKey(account),
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            MaxGasAmount,
		GasUnitPrice:            GasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + TxExpireSec,
		ChainId:                 uint8(ledgerInfo.ChainId),
	}
	// estimate gas and resign
	maxGas, err := c.EstimateMaxGasAmountBCS(account.PublicKey(), txAbi)
	if err != nil {
		return nil, err
	}
	txAbi.MaxGasAmount = maxGas
	return txAbi, nil
}

/**
 * This creates an account if it does not exist and mints the specified amount of
 * coins into that account
 * @param address Hex-encoded 16 bytes Aptos account address wich mints tokens
 * @param amount Amount of tokens to mint
 * @param faucetUrl default https://faucet.devnet.aptoslabs.com
 * @returns Hashes of submitted transactions, e.g. "hash1,has2,hash3,..."
 */
func FaucetFundAccount(address string, amount int64, faucetUrl string) (h *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	hashs, err := aptosclient.FaucetFundAccount(address, uint64(amount), faucetUrl)
	if err != nil {
		return
	}
	return &base.OptionalString{Value: strings.Join(hashs[:], ",")}, nil
}

func toBaseTransaction(transaction *aptostypes.Transaction) (*base.TransactionDetail, error) {
	if transaction.Type != aptostypes.TypeUserTransaction ||
		transaction.Payload.Type != aptostypes.EntryFunctionPayload {
		return nil, errors.New("Invalid transfer transaction.")
	}

	detail := &base.TransactionDetail{
		HashString:  transaction.Hash,
		FromAddress: transaction.Sender,
	}

	gasFee := transaction.GasUnitPrice * transaction.GasUsed
	detail.EstimateFees = strconv.FormatUint(gasFee, 10)

	args := transaction.Payload.Arguments
	if len(args) >= 2 {
		detail.ToAddress = args[0].(string)
		detail.Amount = args[1].(string)
	}

	if transaction.Success {
		detail.Status = base.TransactionStatusSuccess
	} else {
		detail.Status = base.TransactionStatusFailure
		detail.FailureMessage = transaction.VmStatus
	}

	timestamp := transaction.Timestamp / 1e6
	detail.FinishTimestamp = int64(timestamp)

	return detail, nil
}

func getAuthKey(account base.Account) txbuilder.AccountAddress {
	key, _ := hex.DecodeString(account.Address()[2:])
	var a txbuilder.AccountAddress
	copy(a[:], key)
	return a
}
