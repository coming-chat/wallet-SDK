package btc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type Chain struct {
	*Util
}

func NewChainWithChainnet(chainnet string) (*Chain, error) {
	util, err := NewUtilWithChainnet(chainnet)
	if err != nil {
		return nil, err
	}

	return &Chain{Util: util}, nil
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return c
}

func (c *Chain) BalanceOfAddress(address string) (*base.Balance, error) {
	b, err := queryBalance(address, c.Chainnet)
	if err != nil {
		return nil, err
	}
	return &base.Balance{
		Total:  b,
		Usable: b,
	}, nil
}
func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	b, err := queryBalancePubkey(publicKey, c.Chainnet)
	if err != nil {
		return nil, err
	}
	return &base.Balance{
		Total:  b,
		Usable: b,
	}, nil
}
func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfPublicKey(account.PublicKeyHex())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (string, error) {
	tx, err := DecodeTx(signedTx)
	if err != nil {
		return "", err
	}
	hash, err := c.sendWireMsgTx(tx)
	if err != nil {
		return "", err
	}
	return hash.Value, nil
}

func (c *Chain) SendSignedTransaction(signedTxn base.SignedTransaction) (hash *base.OptionalString, err error) {
	if brc20MintTxn, ok := signedTxn.(*Brc20MintTransaction); ok {
		return brc20MintTxn.PublishWithChain(c)
	}
	if psbtTxn, ok := signedTxn.(*SignedPsbtTransaction); ok {
		msgTx, err := PsbtPacketToMsgTx(&psbtTxn.Packet)
		if err != nil {
			return nil, base.MapAnyToBasicError(err)
		}
		return c.sendWireMsgTx(msgTx)
	}
	if txn, ok := signedTxn.(*SignedTransaction); ok {
		return c.sendWireMsgTx(txn.msgTx)
	}
	return nil, base.ErrInvalidTransactionType
}

func (c *Chain) sendWireMsgTx(tx *wire.MsgTx) (hash *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	client, err := rpcClientOf(c.Chainnet)
	if err != nil {
		return
	}
	hashv, err := client.SendRawTransaction(tx, false)
	if err != nil {
		return
	}
	return base.NewOptionalString(hashv.String()), nil
}

// Fetch transaction details through transaction hash
// Note: The input parsing of bitcoin is very complex and the network cost is relatively high,
// So only the status and timestamp can be queried.
func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	return fetchTransactionDetail(hash, c.Chainnet)
}

func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	return fetchTransactionStatus(hash, c.Chainnet)
}

func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	return sdkBatchTransactionStatus(hashListString, c.Chainnet)
}

func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}
func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction base.Transaction, pubkey string) (fee *base.OptionalString, err error) {
	return c.EstimateTransactionFee(transaction)
}

type FeeRate struct {
	Low     int64
	Average int64
	High    int64
}

func (c *Chain) SuggestFeeRate() (rates *FeeRate, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	url := ""
	switch c.Chainnet {
	case ChainBitcoin, ChainMainnet:
		url = "https://mempool.space/api/v1/fees/recommended"
	case ChainTestnet:
		url = "https://mempool.space/testnet/api/v1/fees/recommended"
	case ChainSignet:
		url = "https://mempool.space/signet/api/v1/fees/recommended"
	default:
		return &FeeRate{1, 1, 1}, nil
	}
	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return
	}
	if response.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}

	var feeRates struct {
		MinimumFee  float64 `json:"minimumFee"`
		HalfHourFee float64 `json:"halfHourFee"`
		FastestFee  float64 `json:"fastestFee"`
		// economyFee, hourFee
	}
	err = json.Unmarshal(response.Body, &feeRates)
	if err != nil {
		return &FeeRate{1, 1, 1}, nil
	}
	return &FeeRate{
		Low:     int64(feeRates.MinimumFee),
		Average: int64(feeRates.HalfHourFee),
		High:    int64(feeRates.FastestFee),
	}, nil
}

func (c *Chain) PushPsbt(psbtHex string) (hash *base.OptionalString, err error) {
	packet, err := DecodePsbtTxToPacket(psbtHex)
	if err != nil {
		return nil, err
	}
	err = EnsurePsbtFinalize(packet)
	if err != nil {
		return nil, errors.New("transaction signature error")
	}
	signedTxn := SignedPsbtTransaction{*packet}
	return c.SendSignedTransaction(&signedTxn)
}
