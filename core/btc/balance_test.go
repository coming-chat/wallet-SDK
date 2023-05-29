package btc

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
)

func TestBatchQueryBalance(t *testing.T) {
	addresses := base.StringArray{
		Values: []string{
			"123", "456", // invalid address
			"15MdAHnkxt9TMC2Rj595hsg8Hnv693pPBB", "bc1qa5wkzvf775vxddzaaru2hacd3mj0ehsh3g4anx",
		},
	}

	balances, err := BatchQueryBalance(&addresses, ChainMainnet)
	t.Log(balances, err)
}
