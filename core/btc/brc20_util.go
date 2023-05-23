package btc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

// the v2 resp.Body should like `{code: *, msg: *, data: *}`
func decodeUnisatResponseV2(resp httpUtil.Res, out interface{}) error {
	if resp.Code != http.StatusOK {
		return fmt.Errorf("code: %v, body: %v", resp.Code, string(resp.Body))
	}
	var data struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}
	err := json.Unmarshal(resp.Body, &data)
	if err != nil {
		return err
	}
	if data.Code != 0 || data.Msg != "ok" {
		return fmt.Errorf("code: %v, message: %v", data.Code, data.Msg)
	}
	return json.Unmarshal(data.Data, out)
}

// the v4 resp.Body should like `{status: *, message: *, result: *}`
func decodeUnisatResponseV4(resp httpUtil.Res, out interface{}) error {
	if resp.Code != http.StatusOK {
		return fmt.Errorf("code: %v, body: %v", resp.Code, string(resp.Body))
	}
	var data struct {
		Code string          `json:"status"`
		Msg  string          `json:"message"`
		Data json.RawMessage `json:"result"`
	}
	err := json.Unmarshal(resp.Body, &data)
	if err != nil {
		return err
	}
	if data.Code != "1" || data.Msg != "OK" {
		return fmt.Errorf("code: %v, message: %v", data.Code, data.Msg)
	}
	return json.Unmarshal(data.Data, out)
}
