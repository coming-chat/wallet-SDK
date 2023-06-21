package starknet

import (
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/utils"
	"github.com/stretchr/testify/require"
)

func TestDeployAccountTransactionHash(t *testing.T) {
	pub := "0xc24ee1f993471a03a72d4d5e9f21e91296db50811e2298cc36224552afd5c1"

	txn, err := deployAccountTxnForArgentX(pub)
	require.Nil(t, err)

	txn.MaxFee = new(felt.Felt).SetUint64(10000000000)
	txhash, err := deployAccountTransactionHash(txn, utils.Network(utils.GOERLI))
	require.Nil(t, err)
	require.Equal(t, txhash.String(), "0x794ddc51a8298b57064667cd8fb9ef79d7410c71d8f8ad8098b4462520f582e")
}
