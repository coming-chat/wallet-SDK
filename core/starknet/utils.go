package starknet

import (
	"math/rand"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/xiang-xx/starknet.go/rpc"
)

var latestBlockId = rpc.BlockID{Tag: "latest"}

func mustFelt(str string) *felt.Felt {
	f, _ := new(felt.Felt).SetString(str)
	return f
}

func random(max uint64) uint64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Uint64() % max
}
