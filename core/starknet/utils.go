package starknet

import (
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
)

var latestBlockId = rpc.BlockID{Tag: "latest"}

func mustFelt(str string) *felt.Felt {
	f, _ := new(felt.Felt).SetString(str)
	return f
}
