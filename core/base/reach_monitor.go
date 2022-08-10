package base

import (
	"encoding/json"
	"sort"
	"strings"
	"time"
)

const reachFailedTime int64 = -1

type RpcLatency struct {
	RpcUrl  string `json:"rpcUrl"`
	Latency int64  `json:"latency"`
	Height  int64  `json:"height"`
}

type RpcReachability interface {
	LatencyOf(rpc string, timeout int64) (l *RpcLatency, err error)
}

type ReachMonitorDelegate interface {
	// A node has received a response
	ReachabilityDidReceiveNode(monitor *ReachMonitor, latency *RpcLatency)
	// A node request failed
	ReachabilityDidFailNode(monitor *ReachMonitor, latency *RpcLatency)
	// The entire network connection test task is over
	// @param overview Overview of the results of all connection tests
	ReachabilityDidFinish(monitor *ReachMonitor, overview string)
}

type ReachMonitor struct {
	// The number of network connectivity tests to be performed per rpc. 0 means infinite, default is 1
	ReachCount int
	// Timeout for each connectivity test (ms). default 20000ms
	Timeout int64
	// Time interval between two network connectivity tests (ms). default 1500ms
	Delay int64

	reachability RpcReachability

	stoped bool
}

func NewReachMonitorWithReachability(reachability RpcReachability) *ReachMonitor {
	return &ReachMonitor{
		reachability: reachability,

		ReachCount: 1,
		Timeout:    20000,
		Delay:      1500,
	}
}

func (r *ReachMonitor) StopConnectivity() {
	r.stoped = true
}

// @param rpcList string of rpcs like "rpc1,rpc2,rpc3,..."
func (r *ReachMonitor) StartConnectivityDelegate(rpcList string, delegate ReachMonitorDelegate) {
	if delegate == nil {
		println("You execute the method StartConnectivityDelegate() without listening for any callbacks!!! You should better set a delegate.")
	} else {
		r.startConnectivity(rpcList, delegate)
	}
}

// @param rpcList string of rpcs like "rpc1,rpc2,rpc3,..."
// @return jsonString sorted array base of tatency like "[{rpcUrl:rpc1,latency:100}, {rpcUrl:rpc2, latency:111}, ...]" latency unit is ms. -1 means the connection failed
func (r *ReachMonitor) StartConnectivitySync(rpcList string) string {
	return r.startConnectivity(rpcList, nil)
}

func (r *ReachMonitor) startConnectivity(rpcList string, delegate ReachMonitorDelegate) string {
	r.stoped = false
	successCall, successOk := delegate.(interface {
		ReachabilityDidReceiveNode(tester *ReachMonitor, latency *RpcLatency)
	})
	failCall, failOk := delegate.(interface {
		ReachabilityDidFailNode(tester *ReachMonitor, latency *RpcLatency)
	})

	rpcUrlList := strings.Split(rpcList, ",")
	list := make([]interface{}, len(rpcUrlList))
	for i, s := range rpcUrlList {
		list[i] = s
	}
	temp, _ := MapListConcurrent(list, 0, func(i interface{}) (interface{}, error) {
		var totalCost int64 = 0
		var latestLatency *RpcLatency
		successTimes := 0
		url := i.(string)
		// The result of the first request may be highly skewed, So we need to request one more time
		for c := 0; c < r.ReachCount+1; c++ {
			if r.stoped {
				break
			}
			latency, err := r.reachability.LatencyOf(url, r.Timeout)
			if c == 0 {
				// Ignore the result of the first request
				continue
			}
			// fmt.Printf("... connect %v %v, cost: %v \n", c, url, latency.Latency)
			if err != nil {
				latency.Latency = reachFailedTime
				if failOk {
					failCall.ReachabilityDidFailNode(r, latency)
				}
			} else {
				if successOk {
					successCall.ReachabilityDidReceiveNode(r, latency)
				}
				successTimes++
				totalCost += latency.Latency
				latestLatency = latency
			}
			if c <= r.ReachCount-1 {
				time.Sleep(time.Duration(r.Delay * int64(time.Millisecond)))
			}
		}
		if successTimes == 0 {
			latestLatency = &RpcLatency{
				RpcUrl:  url,
				Latency: reachFailedTime,
				Height:  -1,
			}
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
		ReachabilityDidFinish(tester *ReachMonitor, overview string)
	})
	if finishOk {
		finishCall.ReachabilityDidFinish(r, res)
	}
	return res
}
