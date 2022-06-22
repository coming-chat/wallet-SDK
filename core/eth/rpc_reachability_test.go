package eth

import (
	"strings"
	"testing"
)

func TestRpcReachability_Test(t *testing.T) {
	reach := NewRpcReachability()
	reach.ReachCount = 1
	reach.Delay = 3000
	reach.Timeout = 1000
	t.Log(reach)

	rpcUrls := []string{rpcs.ethereumProd.url, rpcs.binanceTest.url}
	rpcListString := strings.Join(rpcUrls, ",")
	res := reach.StartConnectivityTest(rpcListString)
	t.Log(res)
}
