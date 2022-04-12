package wallet

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/eth"
	CustomType "github.com/coming-chat/wallet-SDK/core/substrate/types"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
	"github.com/itering/subscan/util/ss58"
)

type PolkaBalance struct {
	Total  string
	Usable string
}

func emptyBalance() *PolkaBalance {
	return &PolkaBalance{
		Total:  "0",
		Usable: "0",
	}
}

type PolkaChain struct {
	RpcUrl  string
	ScanUrl string
}

// 通过 url 创建对象
// @param rpcUrl 链端 rpc 地址
// @param scanUrl 浏览器地址(查询交易详情需要的)
// 		chainx 线上: https://multiscan-api.coming.chat/chainx
// 		minix  测试: https://multiscan-api-pre.coming.chat/minix
func NewPolkaChain(rpcUrl, scanUrl string) *PolkaChain {
	return &PolkaChain{
		RpcUrl:  rpcUrl,
		ScanUrl: scanUrl,
	}
}

// 通过 url 和 metadata string 创建对象
// @param rpcUrl 链端 rpc 地址
// @param scanUrl 浏览器地址(查询交易详情需要的)
func NewPolkaChainWithRpc(rpcUrl, scanUrl string, metadataString string) (*PolkaChain, error) {
	_, err := getPolkaClientWithMetadata(rpcUrl, metadataString)
	if err != nil {
		return nil, err
	}
	return &PolkaChain{
		RpcUrl:  rpcUrl,
		ScanUrl: scanUrl,
	}, nil
}

// 获取该链的 metadata string (如果没有会自动下载)
func (c *PolkaChain) GetMetadataString() (s string, err error) {
	client, err := getPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}
	return client.MetadataString()
}

// 刷新最新的 metadata (可以从返回值读取到最新的 metadata)
func (c *PolkaChain) ReloadMetadata() (s string, err error) {
	client, err := getPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}

	err = client.ReloadMetadata()
	if err != nil {
		return
	}

	return client.MetadataString()
}

// 通过 address 查询余额
func (c *PolkaChain) QueryBalance(address string) (*PolkaBalance, error) {
	ss58Format := base58.Decode(address)
	pubkey, err := hex.DecodeString(ss58.Decode(address, int(ss58Format[0])))
	if err != nil {
		return emptyBalance(), err
	}
	return c.queryBalance(pubkey)
}

// 通过 public key 查询余额
func (c *PolkaChain) QueryBalancePubkey(pubkey string) (*PolkaBalance, error) {
	pubkey = strings.TrimPrefix(pubkey, "0x")
	data, err := hex.DecodeString(pubkey)
	if err != nil {
		return emptyBalance(), ErrPublicKey
	}
	return c.queryBalance(data)
}

func (c *PolkaChain) queryBalance(pubkey []byte) (b *PolkaBalance, err error) {
	b = emptyBalance()

	client, err := getPolkaClient(c.RpcUrl)
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

	freeInt := data.Data.Free.Int
	total := freeInt.Add(freeInt, data.Data.Reserved.Int)

	locked := data.Data.MiscFrozen.Int
	if data.Data.MiscFrozen.Cmp(data.Data.FeeFrozen.Int) <= 0 {
		locked = data.Data.FeeFrozen.Int
	}
	usable := freeInt.Sub(freeInt, locked)

	return &PolkaBalance{
		Total:  total.String(),
		Usable: usable.String(),
	}, nil
}

// 特殊查询 XBTC 的余额
// 只能通过 chainx 链对象来查询，其他链会抛出 error
func (c *PolkaChain) QueryBalanceXBTC(address string) (b *PolkaBalance, err error) {
	b = emptyBalance()

	client, err := getPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}

	err = client.LoadMetadataIfNotExists()
	if err != nil {
		return
	}

	ss58Format := base58.Decode("5QUEnWNMDFqsbUGpvvtgWGUgiiojnEpLf7581ELLAQyQ1xnT")
	publicKey, err := hex.DecodeString(ss58.Decode("5QUEnWNMDFqsbUGpvvtgWGUgiiojnEpLf7581ELLAQyQ1xnT", int(ss58Format[0])))
	if err != nil {
		return
	}

	assetId, err := types.EncodeToBytes(uint32(1))
	if err != nil {
		return
	}

	metadata := client.metadata
	call, err := types.CreateStorageKey(metadata, "XAssets", "AssetBalance", publicKey, assetId)
	if err != nil {
		return
	}
	entryMetadata, err := metadata.FindStorageEntryMetadata("XAssets", "AssetBalance")
	if err != nil {
		return
	}
	i := entryMetadata.(types.StorageEntryMetadataV14).Type.AsMap.Value
	kIndex := metadata.AsMetadataV14.EfficientLookup[i.Int64()].Params[0].Type.Int64()
	vValue := metadata.AsMetadataV14.EfficientLookup[i.Int64()].Params[1].Type.Int64()
	data := CustomType.NewMap(metadata.AsMetadataV14.EfficientLookup[kIndex], metadata.AsMetadataV14.EfficientLookup[vValue])
	_, err = client.api.RPC.State.GetStorageLatest(call, &data)
	if err != nil {
		return
	}

	usable, ok := data.Data["Usable"]
	if !ok {
		return b, errors.New("No usable balance")
	}
	usableInt := usable.(types.U128).Int

	return &PolkaBalance{
		Usable: usableInt.String(),
		Total:  usableInt.String(),
	}, nil
}

// 发起交易
func (c *PolkaChain) SendRawTransaction(txHex string) (s string, err error) {
	client, err := getPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}

	var hashString string
	err = client.api.Client.Call(&hashString, "author_submitExtrinsic", txHex)
	if err != nil {
		return
	}

	return hashString, nil
}

// 查询交易详情
func (c *PolkaChain) FetchTransactionDetail(hashString string) (t *eth.TransactionDetail, err error) {
	if c.ScanUrl == "" {
		return nil, errors.New("Scan url is Empty.")
	}
	url := strings.TrimSuffix(c.ScanUrl, "/") + "/extrinsics/" + hashString

	response, err := httpUtil.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return
	}

	if response.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", response.Code, string(response.Body))
	}
	respDict := make(map[string]interface{})
	err = json.Unmarshal(response.Body, &respDict)
	if err != nil {
		return
	}

	// decode informations
	amount, _ := respDict["txAmount"].(string)
	fee, _ := respDict["fee"].(string)
	from, _ := respDict["signer"].(string)
	to, _ := respDict["txTo"].(string)
	timestamp, _ := respDict["blockTime"].(float64)

	status := eth.TransactionStatusNone
	finalized, _ := respDict["finalized"].(bool)
	if finalized {
		success, _ := respDict["success"].(bool)
		if success {
			status = eth.TransactionStatusSuccess
		} else {
			status = eth.TransactionStatusFailure
		}
	} else {
		status = eth.TransactionStatusPending
	}

	return &eth.TransactionDetail{
		HashString:      hashString,
		Amount:          amount,
		EstimateFees:    fee,
		FromAddress:     from,
		ToAddress:       to,
		Status:          status,
		FinishTimestamp: int64(timestamp),
	}, nil
}

// 获取交易的状态
// @param hashString 交易的 hash
func (c *PolkaChain) FetchTransactionStatus(hashString string) eth.TransactionStatus {
	detail, err := c.FetchTransactionDetail(hashString)
	if err != nil {
		return eth.TransactionStatusFailure
	}
	return detail.Status
}

// SDK 批量获取交易的转账状态，hash 列表和返回值，都只能用字符串，逗号隔开传递
// @param hashListString 要批量查询的交易的 hash，用逗号拼接的字符串："hash1,hash2,hash3"
// @return 批量的交易状态，它的顺序和 hashListString 是保持一致的: "status1,status2,status3"
func (c *PolkaChain) SdkBatchTransactionStatus(hashListString string) string {
	hashList := strings.Split(hashListString, ",")
	statuses, _ := eth.MapListConcurrentStringToString(hashList, func(s string) (string, error) {
		return strconv.Itoa(c.FetchTransactionStatus(s)), nil
	})
	return strings.Join(statuses, ",")
}
