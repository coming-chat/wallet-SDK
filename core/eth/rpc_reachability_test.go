package eth

import (
	"fmt"
	"strings"
	"testing"
)

type dddd struct {
}

func (d *dddd) ReachabilityDidReceiveNode(tester *RpcReachability, latency *RpcLatency) {
	fmt.Printf(".... delegate did receive node %v, latency %v\n", latency.RpcUrl, latency.Latency)
}

func (d *dddd) ReachabilityDidFailNode(tester *RpcReachability, latency *RpcLatency) {
	fmt.Printf(".... delegate did fail node %v, latency %v\n", latency.RpcUrl, latency.Latency)
	tester.StopConnectivity()
}

func (d *dddd) ReachabilityDidFinish(tester *RpcReachability, overview string) {
	fmt.Printf(".... delegate did finish %v\n", overview)
}

func TestRpcReachability_Test(t *testing.T) {
	reach := NewRpcReachability()
	reach.ReachCount = 3
	reach.Delay = 3000
	reach.Timeout = 150
	t.Log(reach)

	rpcUrls := []string{rpcs.ethereumProd.url, rpcs.binanceTest.url}
	rpcListString := strings.Join(rpcUrls, ",")
	// res := reach.StartConnectivitySync(rpcListString)
	// t.Log(res)

	delegate := &dddd{}
	reach.StartConnectivityDelegate(rpcListString, delegate)
}
