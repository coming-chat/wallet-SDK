package ordinal

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	BRC20_OP_DEPLOY   = "deploy"
	BRC20_OP_MINT     = "mint"
	BRC20_OP_TRANSFER = "transfer"

	BRC20_HISTORY_TYPE_INSCRIBE_DEPLOY   = "inscribe-deploy"
	BRC20_HISTORY_TYPE_INSCRIBE_MINT     = "inscribe-mint"
	BRC20_HISTORY_TYPE_INSCRIBE_TRANSFER = "inscribe-transfer"
	BRC20_HISTORY_TYPE_TRANSFER          = "transfer"
	BRC20_HISTORY_TYPE_SEND              = "send"
	BRC20_HISTORY_TYPE_RECEIVE           = "receive"
)

type InscriptionBRC20Content struct {
	Proto        string `json:"p"`
	Operation    string `json:"op"`
	BRC20Tick    string `json:"tick"`
	BRC20Max     string `json:"max,omitempty"`
	BRC20Limit   string `json:"lim,omitempty"`
	BRC20Amount  string `json:"amt,omitempty"`
	BRC20To      string `json:"to,omitempty"`
	BRC20Decimal string `json:"dec,omitempty"`
}

func CheckBrc20FromContent(contect []byte) (*InscriptionBRC20Content, error) {
	var bodyMap = make(map[string]interface{}, 8)
	if err := json.Unmarshal(contect, &bodyMap); err != nil {
		return nil, err
	}
	var body InscriptionBRC20Content
	if v, ok := bodyMap["p"].(string); ok {
		body.Proto = v
	}
	if v, ok := bodyMap["op"].(string); ok {
		body.Operation = v
	}
	if v, ok := bodyMap["tick"].(string); ok {
		body.BRC20Tick = v
	}
	if v, ok := bodyMap["max"].(string); ok {
		body.BRC20Max = v
	}
	if _, ok := bodyMap["lim"]; !ok {
		body.BRC20Limit = body.BRC20Max
	} else {
		if v, ok := bodyMap["lim"].(string); ok {
			body.BRC20Limit = v
		}
	}
	if v, ok := bodyMap["amt"].(string); ok {
		body.BRC20Amount = v
	}
	if v, ok := bodyMap["to"].(string); ok {
		body.BRC20To = v
	}

	if _, ok := bodyMap["dec"]; !ok {
		body.BRC20Decimal = "18"
	} else {
		if v, ok := bodyMap["dec"].(string); ok {
			body.BRC20Decimal = v
		}
	}

	if body.Proto != "brc-20" || len(body.BRC20Tick) != 4 {
		return nil, fmt.Errorf("the proto or tick not valid")
	}

	uniqueLowerTicker := strings.ToLower(body.BRC20Tick)
	switch body.Operation {
	case BRC20_OP_DEPLOY:
		if body.BRC20Max == "" { // without max
			return nil, fmt.Errorf("ProcessUpdateLatestBRC20 deploy, but max missing. ticker: %s", uniqueLowerTicker)
		}
		// dec
		dec, err := strconv.ParseUint(body.BRC20Decimal, 10, 64)
		if err != nil || dec > 18 {
			// dec invalid
			return nil, fmt.Errorf("ProcessUpdateLatestBRC20 deploy, but dec invalid. ticker: %s, dec: %s",
				uniqueLowerTicker,
				body.BRC20Decimal,
			)
		}

		// max
		if max, precision, err := NewDecimalFromString(body.BRC20Max); err != nil {
			// max invalid
			return nil, fmt.Errorf("ProcessUpdateLatestBRC20 deploy, but max invalid. ticker: %s, max: '%s'",
				uniqueLowerTicker,
				body.BRC20Max,
			)
		} else {
			if max.Sign() <= 0 || max.IsOverflowUint64() || precision > int(dec) {

			}
		}

		// lim
		if lim, precision, err := NewDecimalFromString(body.BRC20Limit); err != nil {
			// limit invalid
			return nil, fmt.Errorf("ProcessUpdateLatestBRC20 deploy, but limit invalid. ticker: %s, limit: '%s'",
				uniqueLowerTicker,
				body.BRC20Limit,
			)
		} else {
			if lim.Sign() <= 0 || lim.IsOverflowUint64() || precision > int(dec) {
				return nil, fmt.Errorf("ticker: %s, not valid lim", uniqueLowerTicker)
			}
		}
		return &body, nil
	case BRC20_OP_MINT:
		// check mint amount
		_, _, err := NewDecimalFromString(body.BRC20Amount)
		if err != nil {
			return nil, fmt.Errorf("ProcessUpdateLatestBRC20 mint, but amount invalid. ticker: %s, amount: '%s'",
				uniqueLowerTicker,
				body.BRC20Amount,
			)
		}
		//if precision > int(body.Decimal) {
		//	return nil, fmt.Errorf()
		//}
		//if amt.Sign() <= 0 || amt.Cmp(body.Limit) > 0 {
		//	return nil, fmt.Errorf()
		//}
		return &body, nil
	case BRC20_OP_TRANSFER:
		// check amount
		_, _, err := NewDecimalFromString(body.BRC20Amount)
		if err != nil {
			return nil, fmt.Errorf("ProcessUpdateLatestBRC20 inscribe transfer, but amount invalid. ticker: %s, amount: '%s'",
				body.BRC20Tick,
				body.BRC20Amount,
			)
		}

		//if precision > int(body.BRC20Decimal) {
		//	return nil, fmt.Errorf()
		//}
		//if amt.Sign() <= 0 || amt.Cmp(tinfo.Max) > 0 {
		//	return nil, fmt.Errorf()
		//}
		return &body, nil
	default:
		return nil, fmt.Errorf("ticket: %s, has unknown op: %s", uniqueLowerTicker, body.Operation)
	}
}
