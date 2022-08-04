package polka

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type RpcReachability struct {
}

func NewRpcReachability() *RpcReachability {
	return &RpcReachability{}
}

// @return latency (ms) of rpc query blockNumber. -1 means the connection failed.
func (r *RpcReachability) LatencyOf(rpc string, timeout int64) (l *base.RpcLatency, err error) {
	l = &base.RpcLatency{
		RpcUrl:  rpc,
		Latency: -1,
		Height:  -1,
	}

	timeStart := time.Now() // Time Start
	params := httpUtil.RequestParams{
		Header:  map[string]string{"Content-Type": "application/json"},
		Body:    []byte(`{"jsonrpc":"2.0","method":"chain_getHeader","id":13}`),
		Timeout: time.Duration(timeout * int64(time.Millisecond)),
	}
	response, err := httpUtil.Post(rpc, params)
	if err != nil {
		return l, err
	}

	model := struct {
		Result struct {
			Number string `json:"number"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(response, &model)
	if err != nil {
		return l, err
	}
	heightInt, ok := big.NewInt(0).SetString(model.Result.Number, 0)
	if !ok {
		heightInt = big.NewInt(0)
	}
	timeCost := time.Since(timeStart) // Time End

	l.Height = heightInt.Int64()
	l.Latency = timeCost.Milliseconds()
	return l, nil
}
