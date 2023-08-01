package btc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

// the v2 resp.Body should like `{code: *, msg: *, data: *}`
func decodeUnisatResponseV2(resp httpUtil.Res, out interface{}) error {
	err := responseJsonCheck(resp)
	if err != nil {
		return err
	}
	var data struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}
	err = json.Unmarshal(resp.Body, &data)
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
	err := responseJsonCheck(resp)
	if err != nil {
		return err
	}
	var data struct {
		Code string          `json:"status"`
		Msg  string          `json:"message"`
		Data json.RawMessage `json:"result"`
	}
	err = json.Unmarshal(resp.Body, &data)
	if err != nil {
		return err
	}
	if data.Code != "1" || data.Msg != "OK" {
		return fmt.Errorf("code: %v, message: %v", data.Code, data.Msg)
	}
	return json.Unmarshal(data.Data, out)
}

func responseJsonCheck(resp httpUtil.Res) error {
	if resp.Code != http.StatusOK {
		return fmt.Errorf("code: %v, body: %v", resp.Code, string(resp.Body))
	}
	contentType := resp.Header["Content-Type"]
	if len(contentType) <= 0 {
		return fmt.Errorf("response content error")
	}
	if contentType[0] == "text/html" {
		titleRegexp := regexp.MustCompile(`\<title\>(.*)\</title\>`)
		matches := titleRegexp.FindStringSubmatch(string(resp.Body))
		if len(matches) >= 2 {
			return fmt.Errorf("error: %v", matches[1])
		}
		return fmt.Errorf("response content error")
	}
	if !strings.Contains(contentType[0], "json") {
		return fmt.Errorf("response content error")
	}
	return nil
}

func unisatRequestHeader() map[string]string {
	return map[string]string{
		"X-Client":   "UniSat Wallet",
		"X-Version":  "1.1.25",
		"x-address":  "bc1pdq423fm5dv00sl2uckmcve8y3w7guev8ka6qfweljlu23mmsw63qpjc9k7",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
	}
}
