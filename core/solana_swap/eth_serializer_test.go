package solanaswap

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_eth_serialize_append(t *testing.T) {
	sl := NewEthSerializer(10)

	sl.AppendU8(14)
	require.Equal(t, sl.Bytes(), []byte{14})

	sl.AppendU16(257)
	require.Equal(t, sl.Bytes(), []byte{14, 1, 1})

	sl.AppendU64(0x123bcd)
	require.Equal(t, sl.Bytes(), []byte{
		14, 1, 1,
		0, 0, 0, 0, 0, 0x12, 0x3b, 0xcd,
	})
}

func Test_eth_serialize(t *testing.T) {
	tests := []struct {
		name      string
		serialize func(s *EthSerializer)
		want      []byte
	}{
		{
			serialize: func(s *EthSerializer) { s.AppendU8(10) },
			want:      []byte{10},
		},
		{
			serialize: func(s *EthSerializer) { s.AppendU16(0x12df) },
			want:      []byte{0x12, 0xdf},
		},
		{
			serialize: func(s *EthSerializer) { s.AppendU256(*big.NewInt(0xabcdef)) },
			want: []byte{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0xab, 0xcd, 0xef},
		},
		{
			serialize: func(s *EthSerializer) { s.AppendBytesWithLenth(nil) },
			want:      []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			serialize: func(s *EthSerializer) { s.AppendBytesWithLenth([]byte("hello")) },
			want: []byte{0, 0, 0, 0, 0, 0, 0, 5,
				'h', 'e', 'l', 'l', 'o'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewEthSerializer(10)
			tt.serialize(s)
			require.Equal(t, tt.want, s.Bytes())
		})
	}
}
