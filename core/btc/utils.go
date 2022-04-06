package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/wallet-SDK/core/eth"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

// 检查地址是否有效
// @param address 比特币地址
// @param chainnet 链名称
func IsValidAddress(address, chainnet string) bool {
	var netParams *chaincfg.Params
	switch chainnet {
	case chainSignet:
		netParams = &chaincfg.SigNetParams
	case chainMainnet, chainBitcoin:
		netParams = &chaincfg.MainNetParams
	default:
		return false
	}

	_, err := btcutil.DecodeAddress(address, netParams)
	return err == nil
}

// 根据地址查余额
func QueryBalance(address, chainnet string) (string, error) {
	host, err := hostOf(chainnet)
	if err != nil {
		return "0", err
	}
	url := host + "/address/" + address

	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return "0", err
	}

	return parseBalanceResponse(response)
}

// 根据公钥查余额
func QueryBalancePubkey(pubkey, chainnet string) (string, error) {
	host, err := hostOf(chainnet)
	if err != nil {
		return "0", err
	}
	url := host + "/pubkey/" + pubkey

	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return "0", err
	}

	return parseBalanceResponse(response)
}

func parseBalanceResponse(response *httpUtil.Res) (string, error) {
	if response.Code != http.StatusOK {
		return "0", fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}
	respDict := make(map[string]interface{})
	err := json.Unmarshal(response.Body, &respDict)
	if err != nil {
		return "0", err
	}

	chainStats, ok := respDict["chain_stats"].(map[string]interface{})
	if !ok {
		return "0", ErrHttpResponseParse
	}
	funded, ok1 := chainStats["funded_txo_sum"].(float64)
	spend, ok2 := chainStats["spent_txo_sum"].(float64)
	if !ok1 || !ok2 {
		return "0", ErrHttpResponseParse
	}

	balance := int64(math.Max(0, funded-spend))
	return strconv.FormatInt(balance, 10), nil
}

// 对交易进行广播
// @param txHex 签名的tx
// @return 交易 hash
func SendRawTransaction(txHex string, chainnet string) (string, error) {
	client, err := getClientFor(chainnet)
	if err != nil {
		return "", err
	}

	tx, err := decodeTx(txHex)
	if err != nil {
		return "", err
	}

	hash, err := client.SendRawTransaction(tx, false)
	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

func decodeTx(txHex string) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	raw, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	if err = tx.Deserialize(bytes.NewReader(raw)); err != nil {
		return nil, err
	}
	return tx, nil
}

func FetchTransactionDetail(hashString, chainnet string) (*eth.TransactionDetail, error) {
	client, err := getClientFor(chainnet)
	if err != nil {
		return nil, err
	}

	hash, err := chainhash.NewHashFromStr(hashString)
	if err != nil {
		return nil, err
	}

	rawResult, err := client.GetRawTransactionVerbose(hash)
	if err != nil {
		return nil, err
	}

	// map TxRawResult to TransactionDetail
	status := eth.TransactionStatusPending
	if rawResult.Confirmations > 0 {
		status = eth.TransactionStatusSuccess
	}

	// TODO: need decode
	fromAddress := ""
	toAddress := ""
	amount := "0"
	fees := "0"
	for _, txin := range rawResult.Vin {
		println(txin)
	}

	return &eth.TransactionDetail{
		HashString: hashString,

		Amount:       amount,
		EstimateFees: fees,

		FromAddress: fromAddress,
		ToAddress:   toAddress,

		Status:          status,
		FinishTimestamp: rawResult.Time,
		// FailureMessage: "", // if rawResult not nil, then there is no failure
	}, nil
}

func hostOf(chainnet string) (string, error) {
	switch chainnet {
	case chainSignet:
		return "https://electrs-pre.coming.chat", nil
	case chainMainnet, chainBitcoin:
		return "https://electrs-mainnet.coming.chat", nil
	default:
		return "", ErrUnsupportedChain
	}
}
