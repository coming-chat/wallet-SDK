package solanaswap

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_serialize(t *testing.T) {
	tests := []struct {
		name      string
		serialize func(buf *[]byte)
		want      []byte
	}{
		{
			serialize: func(buf *[]byte) { serialize_u8(buf, 10) },
			want:      []byte{10},
		},
		{
			serialize: func(buf *[]byte) { serialize_u16(buf, 0x12df) },
			want:      []byte{0x12, 0xdf},
		},
		{
			serialize: func(buf *[]byte) { serialize_u256(buf, *big.NewInt(0xabcdef)) },
			want: []byte{
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0xab, 0xcd, 0xef},
		},
		{
			serialize: func(buf *[]byte) { serialize_vector_with_length(buf, nil) },
			want:      []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			serialize: func(buf *[]byte) { serialize_vector_with_length(buf, []byte("hello")) },
			want: []byte{0, 0, 0, 0, 0, 0, 0, 5,
				'h', 'e', 'l', 'l', 'o'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := []byte{}
			tt.serialize(&buf)
			require.Equal(t, buf, tt.want)
		})
	}
}
