package aptos

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClaimTokenFromHash(t *testing.T) {
	privateKey := os.Getenv("PriPetra1")
	receiver, err := AccountWithPrivateKey(privateKey)
	t.Log("receiver address:", receiver.Address())
	require.Nil(t, err)

	offerHash := "0x8a9673937e7d4f01a7f305dfd8ad18d29998540eb216e4aa304f4c68a3717f46"

	chain := NewChainWithRestUrl(testnetRestUrl)
	signedTxn, err := chain.ClaimTokenFromHash(receiver, offerHash)
	require.Nil(t, err)
	t.Log(signedTxn)
}
