package solanaswap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerializeBytes(t *testing.T) {
	bytes := []byte("abcdef")
	sl := NewSolanaSerializer(4)
	sl.AppendBytesWithLenth(bytes)
	require.Equal(t, sl.Bytes(), []byte{
		6, 0, 0, 0,
		'a', 'b', 'c', 'd', 'e', 'f',
	})
}
