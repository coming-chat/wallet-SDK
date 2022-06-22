package eth

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
)

const reachFailedTime int64 = -1

type RpcReachability struct {
	// The number of network connectivity tests to be performed per rpc. 0 means infinite, default is 1
	ReachCount int
	// Timeout for each connectivity test (ms). default 20000ms
	Timeout int64
	// Time interval between two network connectivity tests (ms). default 1500ms
	Delay int64
}

func NewRpcReachability() *RpcReachability {
	return &RpcReachability{
		ReachCount: 1,
		Timeout:    20000,
		Delay:      1500,
	}
}

type RpcLatency struct {
	RpcUrl  string `json:"rpcUrl"`
	Latency int64  `json:"latency"`
	Height  int64  `json:"height"`
}

// @param rpcList string of rpcs like "rpc1,rpc2,rpc3,..."
// @return jsonString sorted array base of tatency like "[{rpcUrl:rpc1,latency:100}, {rpcUrl:rpc2, latency:111}, ...]" latency unit is ms. -1 means the connection failed
func (r *RpcReachability) StartConnectivityTest(rpcList string) string {
	rpcUrlList := strings.Split(rpcList, ",")
	list := make([]interface{}, len(rpcUrlList))
	for i, s := range rpcUrlList {
		list[i] = s
	}
	temp, _ := base.MapListConcurrent(list, func(i interface{}) (interface{}, error) {
		var totalCost int64 = 0
		var latestLatency *RpcLatency
		successTimes := r.ReachCount
		url := i.(string)
		for c := 0; c < r.ReachCount; c++ {
			latestLatency = r.latencyOf(url)
			// fmt.Printf("... connect %v %v, cost: %v \n", c, url, cost)
			if latestLatency.Latency == reachFailedTime {
				successTimes--
			} else {
				totalCost += latestLatency.Latency
			}
			if c < r.ReachCount-1 {
				time.Sleep(time.Duration(r.Delay * int64(time.Millisecond)))
			}
		}
		if successTimes == 0 {
			latestLatency.Latency = reachFailedTime
		} else {
			latestLatency.Latency = totalCost / int64(successTimes)
			// fmt.Printf("... connect finish %v, total: %v, avg: %v \n", url, totalCost, latency.Latency)
		}
		return latestLatency, nil
	})

	sort.Slice(temp, func(i, j int) bool {
		ii := temp[i].(*RpcLatency)
		jj := temp[j].(*RpcLatency)
		if ii.Latency == reachFailedTime {
			return false
		}
		if jj.Latency == reachFailedTime {
			return true
		}
		return ii.Latency < jj.Latency
	})
	data, err := json.Marshal(temp)
	if err != nil {
		return ""
	}

	return string(data)
}

// @return latency (ms) of rpc query blockNumber. -1 means the connection failed.
func (r *RpcReachability) latencyOf(rpc string) (l *RpcLatency) {
	l = &RpcLatency{
		RpcUrl:  rpc,
		Latency: reachFailedTime,
		Height:  -1,
	}
	chain, err := GetConnection(rpc)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.Timeout*int64(time.Millisecond)))
	defer cancel()
	timeStart := time.Now()
	height, err := chain.RemoteRpcClient.BlockNumber(ctx)
	timeCost := time.Since(timeStart)
	if err != nil {
		return
	}

	l.Height = int64(height)
	l.Latency = timeCost.Milliseconds()
	return l
}
