package polka

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/coming-chat/wallet-SDK/core/base"
)

// 功能和 GetSignData 相同，不需要提供 nonce, version 等参数，但需要提供 chain 对象和地址
func (c *Chain) GetSignData(t *Transaction, walletAddress string) (data []byte, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cl, err := getConnectedPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}
	var nonce int64
	err = client.CallWithBlockHash(cl.api.Client, &nonce, "system_accountNextIndex", nil, walletAddress)
	if err != nil {
		return
	}
	genesisHash, err := cl.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return
	}
	runtimeVersion, err := cl.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return
	}

	return t.GetSignData(genesisHash.Hex(), nonce, int32(runtimeVersion.SpecVersion), int32(runtimeVersion.TransactionVersion))
}

type MiniXScriptHash struct {
	ScriptHash  string
	BlockNumber int32
}

// 获取 mini 多签转账时需要的 scriptHash
// @param transferTo 转账目标地址
// @param amount 要转出的金额
func (c *Chain) FetchScriptHashForMiniX(transferTo, amount string) (sh *MiniXScriptHash, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cl, err := getConnectedPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}

	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		return
	}

	signedBlock, err := cl.api.RPC.Chain.GetBlockLatest()
	if err != nil {
		return
	}
	blockNumber := uint64(signedBlock.Block.Header.Number)
	arrNumber := make([]uint64, 0)
	arrNumber = append(arrNumber, blockNumber)
	arrNumber = append(arrNumber, blockNumber+1000)
	arr := make([]interface{}, 0)
	arr = append(arr, transferTo, "Transfer", amountInt, arrNumber)

	param := make(map[string]interface{})
	param["id"] = 1
	param["jsonrpc"] = "2.0"
	param["method"] = "ts_computeScriptHash"
	param["params"] = arr
	body, err := c.post(c.RpcUrl, param)
	if err != nil {
		return
	}
	value := make(map[string]interface{})
	err = json.Unmarshal(body, &value)
	if err != nil {
		return
	}

	scriptHash, ok := value["result"].(string)
	if !ok {
		return nil, errors.New("mini http get script hash result error")
	}
	return &MiniXScriptHash{
		ScriptHash:  scriptHash,
		BlockNumber: int32(blockNumber),
	}, nil
}

func (c *Chain) post(baseUrl string, param map[string]interface{}) (body []byte, err error) {
	client := &http.Client{}
	bytesData, _ := json.Marshal(param)
	req, err := http.NewRequest("POST", baseUrl, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	//http.Header.Set(req.Header, "Content-Type", "application/json")
	resp, err := client.Do(req)
	if resp == nil || err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Post " + baseUrl + " response code = " + resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}
