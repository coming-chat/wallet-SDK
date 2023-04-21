package sui

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type RestReachability struct {
}

func NewRestReachability() *RestReachability {
	return &RestReachability{}
}

// @return latency (ms) of rpc query blockNumber. -1 means the connection failed.
func (r *RestReachability) LatencyOf(rpc string, timeout int64) (l *base.RpcLatency, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	l = &base.RpcLatency{
		RpcUrl:  rpc,
		Latency: -1,
		Height:  -1,
	}

	timeStart := time.Now() // Time Start
	params := httpUtil.RequestParams{
		Header:  map[string]string{"Content-Type": "application/json"},
		Body:    []byte(`{"jsonrpc":"2.0","method":"sui_getTotalTransactionBlocks","params":{},"id":33}`),
		Timeout: time.Duration(timeout * int64(time.Millisecond)),
	}
	response, err := httpUtil.Post(rpc, params)
	if err != nil {
		return l, err
	}

	model := struct {
		Result string `json:"result"`
	}{}
	err = json.Unmarshal(response, &model)
	if err != nil {
		return l, err
	}
	heightInt, err := strconv.ParseInt(model.Result, 10, 64)
	if err != nil {
		heightInt = 0
		err = nil
	}
	timeCost := time.Since(timeStart) // Time End

	l.Height = heightInt
	l.Latency = timeCost.Milliseconds()
	return l, nil
}
