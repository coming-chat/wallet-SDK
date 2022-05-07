package polka

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
	"github.com/decred/base58"
	"github.com/itering/subscan/util/ss58"
)

type Chain struct {
	*Util
	RpcUrl  string
	ScanUrl string
}

// @param rpcUrl will be used to get metadata, query balance, estimate fee, send signed tx.
// @param scanUrl will be used to query transaction details
func NewChainWithRpc(rpcUrl, scanUrl string, network int) (*Chain, error) {
	util := NewUtilWithNetwork(network)
	return &Chain{
		Util:    util,
		RpcUrl:  rpcUrl,
		ScanUrl: scanUrl,
	}, nil
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return &Token{chain: c}
}

// Note: Only chainx have XBTC token.
func (c *Chain) XBTCToken() *XBTCToken {
	return &XBTCToken{chain: c}
}

func (c *Chain) BalanceOfAddress(address string) (*base.Balance, error) {
	ss58Format := base58.Decode(address)
	pubkey, err := hex.DecodeString(ss58.Decode(address, int(ss58Format[0])))
	if err != nil {
		return base.EmptyBalance(), err
	}
	return c.queryBalance(pubkey)
}

func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	publicKey = strings.TrimPrefix(publicKey, "0x")
	data, err := hex.DecodeString(publicKey)
	if err != nil {
		return base.EmptyBalance(), ErrPublicKey
	}
	return c.queryBalance(data)
}

func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfPublicKey(account.PublicKeyHex())
}

func (c *Chain) SendRawTransaction(signedTx string) (s string, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	client, err := getConnectedPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}

	var hashString string
	err = client.api.Client.Call(&hashString, "author_submitExtrinsic", signedTx)
	if err != nil {
		return
	}

	return hashString, nil
}

func (c *Chain) FetchTransactionDetail(hashString string) (*base.TransactionDetail, error) {
	if c.ScanUrl == "" {
		return nil, errors.New("Scan url is Empty.")
	}
	url := strings.TrimSuffix(c.ScanUrl, "/") + "/" + hashString

	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	if response.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}
	respDict := make(map[string]interface{})
	err = json.Unmarshal(response.Body, &respDict)
	if err != nil {
		return nil, err
	}

	// decode informations
	amount, _ := respDict["txAmount"].(string)
	fee, _ := respDict["fee"].(string)
	from, _ := respDict["signer"].(string)
	to, _ := respDict["txTo"].(string)
	timestamp, _ := respDict["blockTime"].(float64)

	status := base.TransactionStatusNone
	finalized, _ := respDict["finalized"].(bool)
	if finalized {
		success, _ := respDict["success"].(bool)
		if success {
			status = base.TransactionStatusSuccess
		} else {
			status = base.TransactionStatusFailure
		}
	} else {
		status = base.TransactionStatusPending
	}

	return &base.TransactionDetail{
		HashString:      hashString,
		Amount:          amount,
		EstimateFees:    fee,
		FromAddress:     from,
		ToAddress:       to,
		Status:          status,
		FinishTimestamp: int64(timestamp),
	}, nil
}

func (c *Chain) FetchTransactionStatus(hashString string) base.TransactionStatus {
	detail, err := c.FetchTransactionDetail(hashString)
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

// query balance with pubkey data.
func (c *Chain) queryBalance(pubkey []byte) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	b = base.EmptyBalance()

	client, err := getConnectedPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}

	err = client.LoadMetadataIfNotExists()
	if err != nil {
		return
	}

	call, err := types.CreateStorageKey(client.metadata, "System", "Account", pubkey)
	if err != nil {
		return
	}

	data := struct {
		Nonce       uint32
		Consumers   uint32
		Providers   uint32
		Sufficients uint32
		Data        struct {
			Free       types.U128
			Reserved   types.U128
			MiscFrozen types.U128
			FeeFrozen  types.U128
		}
	}{}

	// Ok is true if the value is not empty.
	ok, err := client.api.RPC.State.GetStorageLatest(call, &data)
	if err != nil {
		return
	}
	if !ok {
		return
	}

	totalInt := big.NewInt(0).Add(data.Data.Free.Int, data.Data.Reserved.Int)
	locked := base.MaxBigInt(data.Data.MiscFrozen.Int, data.Data.FeeFrozen.Int)
	usableInt := big.NewInt(0).Sub(data.Data.Free.Int, locked)

	return &base.Balance{
		Total:  totalInt.String(),
		Usable: usableInt.String(),
	}, nil
}
