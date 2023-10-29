package solanaswap

import (
	"encoding/binary"
	"math/big"
)

func serialize_u8(buf *[]byte, n uint8) {
	*buf = append(*buf, n)
}

func serialize_u16(buf *[]byte, n uint16) {
	*buf = binary.BigEndian.AppendUint16(*buf, n)
}

func serialize_u32(buf *[]byte, n uint32) {
	*buf = binary.BigEndian.AppendUint32(*buf, n)
}

func serialize_u64(buf *[]byte, n uint64) {
	*buf = binary.BigEndian.AppendUint64(*buf, n)
}

func serialize_u128(buf *[]byte, n big.Int) {
	cop := big.NewInt(0).Set(&n)

	data := make([]byte, 0, 128/8)
	for t := 0; t < 128/64; t++ {
		head := binary.BigEndian.AppendUint64(make([]byte, 0), cop.Uint64())
		data = append(head, data...)
		cop = cop.Rsh(cop, 64)
	}
	*buf = append(*buf, data...)
}

func serialize_u256(buf *[]byte, n big.Int) {
	cop := big.NewInt(0).Set(&n)

	data := make([]byte, 0, 256/8)
	for t := 0; t < 256/64; t++ {
		head := binary.BigEndian.AppendUint64(make([]byte, 0), cop.Uint64())
		data = append(head, data...)
		cop = cop.Rsh(cop, 64)
	}
	*buf = append(*buf, data...)
}

func serialize_vector_with_length(buf *[]byte, bytes []byte) {
	if len(bytes) == 0 {
		serialize_u64(buf, 0)
		return
	}
	serialize_u64(buf, uint64(len(bytes)))
	*buf = append(*buf, bytes...)
}
