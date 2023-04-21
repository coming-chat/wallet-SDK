package sui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type dddd struct {
}

func (d *dddd) ReachabilityDidReceiveNode(tester *base.ReachMonitor, latency *base.RpcLatency) {
	fmt.Printf(".... delegate did receive height %v, latency %v, node %v\n", latency.Height, latency.Latency, latency.RpcUrl)
}

func (d *dddd) ReachabilityDidFailNode(tester *base.ReachMonitor, latency *base.RpcLatency) {
	fmt.Printf(".... delegate did fail height %v, latency %v, node %v\n", latency.Height, latency.Latency, latency.RpcUrl)
	// tester.StopConnectivity()
}

func (d *dddd) ReachabilityDidFinish(tester *base.ReachMonitor, overview string) {
	fmt.Printf(".... delegate did finish %v\n", overview)
}

func TestRpcReachability_Test(t *testing.T) {
	reach := NewRestReachability()
	monitor := base.NewReachMonitorWithReachability(reach)
	monitor.ReachCount = 3
	monitor.Delay = 3000
	monitor.Timeout = 2500
	t.Log(reach)

	rpcUrls := []string{TestnetRpcUrl, DevnetRpcUrl}
	rpcListString := strings.Join(rpcUrls, ",")
	// res := reach.StartConnectivitySync(rpcListString)
	// t.Log(res)

	delegate := &dddd{}
	monitor.StartConnectivityDelegate(rpcListString, delegate)
}
