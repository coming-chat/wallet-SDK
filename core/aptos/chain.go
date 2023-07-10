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
	MaxGasAmount = 20000
	GasPrice     = 100
	TxExpireSec  = 600

	FaucetUrlDevnet  = "https://faucet.devnet.aptoslabs.com"
	FaucetUrlTestnet = "https://faucet.testnet.aptoslabs.com"
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
	return NewMainToken(c)
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

func (c *Chain) SendSignedTransaction(signedTxn base.SignedTransaction) (hash *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	txn := AsSignedTransaction(signedTxn)
	if txn == nil {
		return nil, base.ErrInvalidTransactionType
	}
	client, err := c.client()
	if err != nil {
		return
	}
	resultTx, err := client.SubmitSignedBCSTransaction(txn.SignedBytes)
	if err != nil {
		return
	}
	return &base.OptionalString{Value: resultTx.Hash}, nil
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	txn, err := c.fetchDetail(hash)
	if err != nil {
		return nil, err
	}
	return toBaseTransaction(txn)
}

func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	transaction, err := c.fetchDetail(hash)
	if err != nil {
		return base.TransactionStatusNone
	}
	if transaction.Success {
		return base.TransactionStatusSuccess
	} else if transaction.VmStatus == "" {
		return base.TransactionStatusPending
	} else {
		return base.TransactionStatusFailure
	}
}

func (c *Chain) fetchDetail(hashOrVersion string) (txn *aptostypes.Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	client, err := c.client()
	if err != nil {
		return
	}
	if strings.HasPrefix(hashOrVersion, "0x") {
		return client.GetTransactionByHash(hashOrVersion)
	}
	_, err = strconv.ParseUint(hashOrVersion, 10, 64)
	if err == nil {
		return client.GetTransactionByVersion(hashOrVersion)
	} else {
		return client.GetTransactionByHash(hashOrVersion)
	}
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := base.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}

func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	return nil, base.ErrEstimateGasNeedPublicKey
}

func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction base.Transaction, pubkey string) (fee *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	txn, ok := transaction.(*Transaction)
	if !ok {
		return nil, base.ErrInvalidTransactionType
	}
	pubData, err := types.HexDecodeString(pubkey)
	if err != nil {
		return nil, base.ErrInvalidPublicKey
	}
	gas, err := c.EstimateMaxGasAmountBCS(pubData, &txn.RawTxn)
	if err != nil {
		return
	}
	txn.RawTxn.MaxGasAmount = gas // reset MaxGasAmount
	gasFee := gas * txn.RawTxn.GasUnitPrice
	return &base.OptionalString{Value: strconv.FormatUint(gasFee, 10)}, nil
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

func (c *Chain) EstimateMaxGasAmountBCS(publicKey []byte, rawTxn *txbuilder.RawTransaction) (uint64, error) {
	signedTx, err := txbuilder.GenerateBCSSimulation(publicKey, rawTxn)
	if err != nil {
		return 0, err
	}
	client, err := c.client()
	if err != nil {
		return 0, err
	}
	txns, err := client.SimulateSignedBCSTransaction(signedTx)
	if err != nil {
		return 0, err
	}
	return handleGasAmount(txns)
}

func handleGasAmount(txns []*aptostypes.Transaction) (uint64, error) {
	if len(txns) <= 0 {
		return 0, errors.New("Query gas fee failed.")
	}
	txn := txns[0]
	if !txn.Success {
		return 0, errors.New(txn.VmStatus)
	}
	gasFloat := float64(txn.GasUsed)*1.2 + 1
	return uint64(gasFloat), nil
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

func (c *Chain) SignAndSendTransaction(account base.Account, hexData string) (*base.OptionalString, error) {
	var (
		err         error
		client      *aptosclient.RestClient
		txAbi       *txbuilder.RawTransaction
		submittedTx *aptostypes.Transaction
		signedTxn   []byte
	)
	txData, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, err
	}
	txAbi = &txbuilder.RawTransaction{}
	err = lcs.Unmarshal(txData, txAbi)
	if err != nil {
		return nil, err
	}
	if client, err = c.client(); err != nil {
		return nil, err
	}
	aptAccount, ok := account.(*Account)
	if !ok {
		return nil, errors.New("invalid account type")
	}
	if signedTxn, err = txbuilder.GenerateBCSTransaction(aptAccount.account, txAbi); err != nil {
		return nil, err
	}
	if submittedTx, err = client.SubmitSignedBCSTransaction(signedTxn); err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: submittedTx.Hash}, err
}

// @return The raw transaction that `MaxGasAmount` has obtained from the chain in real time.
func (c *Chain) createTransactionFromPayloadBCS(account base.Account, payload txbuilder.TransactionPayload) (*txbuilder.RawTransaction, error) {
	txn, err := c.buildTransactionFromPayloadBCS(account.Address(), payload)
	if err != nil {
		return nil, err
	}
	// estimate gas and resign
	maxGas, err := c.EstimateMaxGasAmountBCS(account.PublicKey(), txn)
	if err != nil {
		return nil, err
	}
	txn.MaxGasAmount = maxGas
	return txn, nil
}

func (c *Chain) buildTransactionFromPayloadBCS(sender string, payload txbuilder.TransactionPayload) (txn *txbuilder.RawTransaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	senderKey, err := txbuilder.NewAccountAddressFromHex(sender)
	if err != nil {
		return nil, err
	}
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	accountData, err := client.GetAccount(sender)
	if err != nil {
		return nil, err
	}
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.EstimateGasPrice()
	if err != nil {
		return nil, err
	}
	txAbi := &txbuilder.RawTransaction{
		Sender:                  *senderKey,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            MaxGasAmount,
		GasUnitPrice:            gasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + TxExpireSec,
		ChainId:                 uint8(ledgerInfo.ChainId),
	}
	return txAbi, nil
}

func (c *Chain) EstimateGasPrice() (*base.OptionalString, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	price, err := client.EstimateGasPrice()
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: strconv.FormatUint(price, 10)}, nil
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
	hashs, _ = base.MapListConcurrentStringToString(hashs, func(s string) (string, error) {
		if !strings.HasPrefix(s, "0x") {
			s = "0x" + s
		}
		return s, nil
	})
	return &base.OptionalString{Value: strings.Join(hashs[:], ",")}, nil
}

func toBaseTransaction(transaction *aptostypes.Transaction) (*base.TransactionDetail, error) {
	if transaction.Type != aptostypes.TypeUserTransaction ||
		transaction.Payload.Type != aptostypes.EntryFunctionPayload {
		return nil, errors.New("Invalid transfer transaction.")
	}

	gasFee := transaction.GasUnitPrice * transaction.GasUsed
	timestamp := transaction.Timestamp / 1e6
	detail := &base.TransactionDetail{
		HashString:      transaction.Hash,
		FromAddress:     transaction.Sender,
		EstimateFees:    strconv.FormatUint(gasFee, 10),
		FinishTimestamp: int64(timestamp),
	}
	if transaction.Success {
		detail.Status = base.TransactionStatusSuccess
	} else if transaction.VmStatus == "" {
		detail.Status = base.TransactionStatusPending
	} else {
		detail.Status = base.TransactionStatusFailure
		detail.FailureMessage = transaction.VmStatus
	}

	function := transaction.Payload.Function
	args := transaction.Payload.Arguments
	switch {
	case function == "0x3::token_transfers::offer_script":
		if len(args) >= 4 {
			detail.ToAddress = args[0].(string)
			detail.TokenName = args[3].(string)
		}
	case function == "0x3::token_transfers::claim_script":
		if len(args) >= 4 {
			detail.FromAddress = args[0].(string)
			detail.ToAddress = transaction.Sender
			detail.TokenName = args[3].(string)
		}
	case strings.HasSuffix(function, "::cid::cid_token_transfer"):
		if len(args) >= 2 {
			detail.ToAddress = args[1].(string)
			detail.CIDNumber = strings.TrimSuffix(args[0].(string), ".aptos")
		}
	case strings.HasSuffix(function, "::cid::token_trasfer"):
		if len(args) >= 5 {
			detail.ToAddress = args[4].(string)
			detail.TokenName = args[2].(string)
		}
	case function == "0x1::coin::transfer" || function == "0x1::aptos_account::transfer":
		if len(args) >= 2 {
			detail.ToAddress = args[0].(string)
			detail.Amount = args[1].(string)
		}
	}

	return detail, nil
}

func getAuthKey(account base.Account) txbuilder.AccountAddress {
	key, _ := hex.DecodeString(account.Address()[2:])
	var a txbuilder.AccountAddress
	copy(a[:], key)
	return a
}
