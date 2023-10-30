package solanaswap

import (
	"math/big"
	"testing"

	"github.com/blocto/solana-go-sdk/pkg/bincode"
	"github.com/stretchr/testify/require"
)

func TestSerializeU256(t *testing.T) {
	u64, err := bincode.SerializeData(uint64(1))
	require.Nil(t, err)
	require.Equal(t, u64, []byte{1, 0, 0, 0, 0, 0, 0, 0})

	bb, _ := big.NewInt(0).SetString("1", 16)
	u256, err := LittleEndianSerializeU256(bb)
	require.Nil(t, err)
	require.Equal(t, u256, []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func TestSerializeBytes(t *testing.T) {
	bytes := []byte("abcdef")
	data := make([]byte, 0)
	buf := &data
	LittleEndianSerializeBytesWithLength(buf, bytes)
	require.Equal(t, data, []byte{6, 0, 0, 0, 0, 0, 0, 0, 'a', 'b', 'c', 'd', 'e', 'f'})
}
