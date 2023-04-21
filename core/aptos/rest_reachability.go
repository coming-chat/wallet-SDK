package aptos

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
	l = &base.RpcLatency{
		RpcUrl:  rpc,
		Latency: -1,
		Height:  -1,
	}

	timeStart := time.Now() // Time Start
	body, err := httpUtil.Get(rpc, nil)
	if err != nil {
		return l, err
	}

	model := struct {
		BlockHeight string `json:"block_height"`
	}{}
	err = json.Unmarshal(body, &model)
	if err != nil {
		return l, err
	}
	heightInt, err := strconv.ParseInt(model.BlockHeight, 10, 64)
	if err != nil {
		heightInt = 0
		err = nil
	}
	timeCost := time.Since(timeStart) // Time End

	l.Height = heightInt
	l.Latency = timeCost.Milliseconds()
	return l, nil
}
