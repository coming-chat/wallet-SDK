package multi_signature_check

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"strings"
)

type JsonRpcReq struct {
	Id      int64         `json:"id"`
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func RequestJsonRpc(url, method string, params ...interface{}) (interface{}, error) {
	data := &JsonRpcReq{
		Id:      1,
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
	}
	body, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	resp, err := httpUtil.Request(http.MethodPost, url, map[string]string{"Content-Type": "application/json; charset=utf-8"}, body)
	if err != nil {
		return nil, err
	}
	if resp.Code != http.StatusOK {
		return nil, fmt.Errorf("code: %d, body: %s", resp.Code, string(resp.Body))
	}
	respData := make(map[string]interface{})
	err = json.Unmarshal(resp.Body, &respData)
	if err != nil {
		return nil, err
	}
	if errData, ok := respData["error"]; ok {
		errData := errData.(map[string]interface{})
		if data, ok := errData["data"]; ok {
			return nil, fmt.Errorf("rpc method or params err{ code: %d, message: %s, data: %s }", int64(errData["code"].(float64)), errData["message"].(string), data.(string))
		}
		return nil, fmt.Errorf("rpc method or params err{ code: %d, message: %s }", int64(errData["code"].(float64)), errData["message"].(string))
	}
	if resultData, ok := respData["result"]; ok {
		return resultData, nil
	}
	return nil, fmt.Errorf("code: %d, body: %s", resp.Code, string(resp.Body))
}

func XGatewayCommonWithdrawalListWithFeeInfo(url string, assertId int) (string, error) {
	result, err := RequestJsonRpc(url, "xgatewaycommon_withdrawalListWithFeeInfo", assertId)
	if err != nil {
		return "", err
	}
	data, ok := result.(map[string]interface{})
	if !ok {
		return "", errors.New("resolve result.map data failed")
	}
	idMap := make(map[string]uint32)
	for k, v := range data {
		v, ok := v.([]interface{})
		if !ok {
			return "", errors.New("resolve result.map.list data failed")
		}
		if len(v) < 2 {
			return "", fmt.Errorf("result data value len is %d, not equal 2", len(v))
		}
		index0, ok := v[0].(map[string]interface{})
		if !ok {
			return "", errors.New("resolve.map.list.map result data failed")
		}
		index1, ok := v[1].(map[string]interface{})
		if !ok {
			return "", errors.New("resolve.map.list.map result data failed")
		}
		if "Applying" == index0["state"].(string) {
			id, err := strconv.ParseUint(k, 10, 32)
			if err != nil {
				return "", err
			}
			balance, err := decimal.NewFromString(index0["balance"].(string))
			if err != nil {
				return "", err
			}
			fee, err := decimal.NewFromString(index1["fee"].(string))
			if err != nil {
				return "", err
			}

			idMap[AddressAmountKey(index0["addr"].(string), balance.Sub(fee).String())] = uint32(id)
		}
	}
	jsonb, err := json.Marshal(idMap)
	if err != nil {
		return "", err
	}
	return string(jsonb), nil
}

func XGatewayBitcoinVerifyTxValid(url, rawTx, withdrawalIds string, isFullAmount bool) (bool, error) {
	withdrawalIdList := strings.Split(withdrawalIds, ",")
	withdrawalIdsU32 := make([]uint32, 0)
	for _, v := range withdrawalIdList {
		if len(v) == 0 {
			continue
		}
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return false, err
		}
		withdrawalIdsU32 = append(withdrawalIdsU32, uint32(id))
	}
	result, err := RequestJsonRpc(url, "xgatewaybitcoin_verifyTxValid", rawTx, withdrawalIdsU32, isFullAmount)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

func AddressAmountKey(address, amount string) string {
	return fmt.Sprintf("%s:%s", address, amount)
}
