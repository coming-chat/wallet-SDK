package solanaswap

import (
	"bytes"
	"encoding/binary"
	"math/big"
)

type EthSerializer struct {
	buf bytes.Buffer
}

func NewEthSerializer(capacity int) *EthSerializer {
	return &EthSerializer{buf: *bytes.NewBuffer(make([]byte, 0, capacity))}
}

func (s *EthSerializer) Bytes() []byte {
	return s.buf.Bytes()
}

func (s *EthSerializer) AppendU8(u uint8) {
	s.buf.Write([]byte{u})
}

func (s *EthSerializer) AppendU16(u uint16) {
	binary.Write(&s.buf, binary.BigEndian, u)
}

func (s *EthSerializer) AppendU32(u uint32) {
	binary.Write(&s.buf, binary.BigEndian, u)
}

func (s *EthSerializer) AppendU64(u uint64) {
	binary.Write(&s.buf, binary.BigEndian, u)
}

func (s *EthSerializer) AppendU128(u big.Int) {
	cop := big.NewInt(0).Set(&u)

	data := make([]byte, 0, 128/8)
	for t := 0; t < 128/64; t++ {
		head := binary.BigEndian.AppendUint64(make([]byte, 0), cop.Uint64())
		data = append(head, data...)
		cop = cop.Rsh(cop, 64)
	}
	s.buf.Write(data)
}

func (s *EthSerializer) AppendU256(u big.Int) {
	cop := big.NewInt(0).Set(&u)

	data := make([]byte, 0, 256/8)
	for t := 0; t < 256/64; t++ {
		head := binary.BigEndian.AppendUint64(make([]byte, 0), cop.Uint64())
		data = append(head, data...)
		cop = cop.Rsh(cop, 64)
	}
	s.buf.Write(data)
}

func (s *EthSerializer) AppendBytesWithLenth(bs []byte) {
	s.AppendU64(uint64(len(bs)))
	s.buf.Write(bs)
}
