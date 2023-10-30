package solanaswap

import (
	"bytes"
	"encoding/binary"
)

type SolanaSerializer struct {
	buf bytes.Buffer
}

func NewSolanaSerializer(capacity int) *SolanaSerializer {
	return &SolanaSerializer{buf: *bytes.NewBuffer(make([]byte, 0, capacity))}
}

func (s *SolanaSerializer) Bytes() []byte {
	return s.buf.Bytes()
}

func (s *SolanaSerializer) AppendU16(u uint16) {
	binary.Write(&s.buf, binary.LittleEndian, u)
}

func (s *SolanaSerializer) AppendU32(u uint32) {
	binary.Write(&s.buf, binary.LittleEndian, u)
}

func (s *SolanaSerializer) AppendBytesWithLenth(bs []byte) {
	s.AppendU32(uint32(len(bs)))
	s.buf.Write(bs)
}
