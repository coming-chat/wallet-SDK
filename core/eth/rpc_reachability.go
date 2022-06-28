package eth

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type RpcReachabilityDelegate interface {
	// A node has received a response
	ReachabilityDidReceiveNode(tester *RpcReachability, latency *RpcLatency)
	// A node request failed
	ReachabilityDidFailNode(tester *RpcReachability, latency *RpcLatency)
	// The entire network connection test task is over
	// @param overview Overview of the results of all connection tests
	ReachabilityDidFinish(tester *RpcReachability, overview string)
}

const reachFailedTime int64 = -1

type RpcReachability struct {
	// The number of network connectivity tests to be performed per rpc. 0 means infinite, default is 1
	ReachCount int
	// Timeout for each connectivity test (ms). default 20000ms
	Timeout int64
	// Time interval between two network connectivity tests (ms). default 1500ms
	Delay int64

	stoped bool
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

func (r *RpcReachability) StopConnectivity() {
	r.stoped = true
}

// @param rpcList string of rpcs like "rpc1,rpc2,rpc3,..."
func (r *RpcReachability) StartConnectivityDelegate(rpcList string, delegate RpcReachabilityDelegate) {
	if delegate == nil {
		println("You execute the method StartConnectivityDelegate() without listening for any callbacks!!! You should better set a delegate.")
	} else {
		r.startConnectivity(rpcList, delegate)
	}
}

// @param rpcList string of rpcs like "rpc1,rpc2,rpc3,..."
// @return jsonString sorted array base of tatency like "[{rpcUrl:rpc1,latency:100}, {rpcUrl:rpc2, latency:111}, ...]" latency unit is ms. -1 means the connection failed
func (r *RpcReachability) StartConnectivitySync(rpcList string) string {
	return r.startConnectivity(rpcList, nil)
}

func (r *RpcReachability) startConnectivity(rpcList string, delegate RpcReachabilityDelegate) string {
	r.stoped = false
	successCall, successOk := delegate.(interface {
		ReachabilityDidReceiveNode(tester *RpcReachability, latency *RpcLatency)
	})
	failCall, failOk := delegate.(interface {
		ReachabilityDidFailNode(tester *RpcReachability, latency *RpcLatency)
	})

	rpcUrlList := strings.Split(rpcList, ",")
	list := make([]interface{}, len(rpcUrlList))
	for i, s := range rpcUrlList {
		list[i] = s
	}
	temp, _ := base.MapListConcurrent(list, 0, func(i interface{}) (interface{}, error) {
		var totalCost int64 = 0
		var latestLatency *RpcLatency
		successTimes := 0
		url := i.(string)
		for c := 0; c < r.ReachCount; c++ {
			if r.stoped {
				break
			}
			latency, err := r.latencyOf(url)
			// fmt.Printf("... connect %v %v, cost: %v \n", c, url, latency.Latency)
			if err != nil {
				if failOk {
					failCall.ReachabilityDidFailNode(r, latency)
				}
			} else {
				if successOk {
					successCall.ReachabilityDidReceiveNode(r, latency)
				}
				successTimes++
				totalCost += latency.Latency
			}
			if c < r.ReachCount-1 {
				time.Sleep(time.Duration(r.Delay * int64(time.Millisecond)))
			}
			latestLatency = latency
		}
		if successTimes == 0 {
			latestLatency.Latency = reachFailedTime
		} else {
			latestLatency.Latency = totalCost / int64(successTimes)
			// fmt.Printf("... connect finish %v, total: %v, avg: %v \n", url, totalCost, latestLatency.Latency)
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
	res := ""
	if err == nil {
		res = string(data)
	}

	finishCall, finishOk := delegate.(interface {
		ReachabilityDidFinish(tester *RpcReachability, overview string)
	})
	if finishOk {
		finishCall.ReachabilityDidFinish(r, res)
	}
	return res
}

// @return latency (ms) of rpc query blockNumber. -1 means the connection failed.
func (r *RpcReachability) latencyOf(rpc string) (l *RpcLatency, err error) {
	l = &RpcLatency{
		RpcUrl:  rpc,
		Latency: reachFailedTime,
		Height:  -1,
	}
	var connectTimeout int64 = 3000
	if r.Timeout > connectTimeout {
		connectTimeout = r.Timeout
	}
	chain, err := getConnectionWithTimeout(rpc, connectTimeout)
	if err != nil {
		return l, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.Timeout*int64(time.Millisecond)))
	defer cancel()
	timeStart := time.Now()
	height, err := chain.RemoteRpcClient.BlockNumber(ctx)
	timeCost := time.Since(timeStart)
	if err != nil {
		return l, err
	}

	l.Height = int64(height)
	l.Latency = timeCost.Milliseconds()
	return l, nil
}
