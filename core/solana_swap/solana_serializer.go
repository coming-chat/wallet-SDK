package solanaswap

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"math/big"

	"github.com/blocto/solana-go-sdk/common"
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

type SolanaDeserializer struct {
	data []byte
	idx  int
}

func NewSolanaDeserializer(b []byte) *SolanaDeserializer {
	return &SolanaDeserializer{
		data: b,
	}
}

func (s *SolanaDeserializer) IsWhirlpoolDataType(typeName string) bool {
	if len(s.data) < ACCOUNT_DISCRIMINATOR_SIZE {
		return false
	}
	return bytes.Compare(s.data[:ACCOUNT_DISCRIMINATOR_SIZE], nameDiscriminator(typeName)) == 0
}

func (s *SolanaDeserializer) StartWhirlpoolDataParse() {
	s.idx = 8
}

const ACCOUNT_DISCRIMINATOR_SIZE = 8

func nameDiscriminator(name string) []byte {
	key := `account:` + name // camelcase
	temp := (sha256.Sum256([]byte(key)))
	return temp[:ACCOUNT_DISCRIMINATOR_SIZE]
}

func (s *SolanaDeserializer) TakeBytes(len int) []byte {
	res := s.data[s.idx : s.idx+len]
	s.idx = s.idx + len
	return res
}

func (s *SolanaDeserializer) TakeBool() bool {
	return s.TakeBytes(1)[0] > 0
}

func (s *SolanaDeserializer) TakeU8() uint8 {
	return s.TakeBytes(1)[0]
}

func (s *SolanaDeserializer) TakeU16() uint16 {
	return binary.LittleEndian.Uint16(s.TakeBytes(2))
}

func (s *SolanaDeserializer) TakeU32() uint32 {
	return binary.LittleEndian.Uint32(s.TakeBytes(4))
}

func (s *SolanaDeserializer) TakeI32() int32 {
	return int32(binary.LittleEndian.Uint32(s.TakeBytes(4)))
}

func (s *SolanaDeserializer) TakeU64() uint64 {
	return binary.LittleEndian.Uint64(s.TakeBytes(8))
}

func (s *SolanaDeserializer) TakeU128() *big.Int {
	return newIntLittleEndianBytes(s.TakeBytes(16))
}

func (s *SolanaDeserializer) TakeI128() *big.Int {
	bts := s.TakeBytes(16)
	u128 := newIntLittleEndianBytes(bts)
	if bts[0]&0x80 == 0 { // positive
		return u128
	} else { // negtive
		pow128 := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil)
		return u128.Sub(u128, pow128)
	}
}

func (s *SolanaDeserializer) TakePublicKey() common.PublicKey {
	return common.PublicKeyFromBytes(s.TakeBytes(common.PublicKeyLength))
}

func (s *SolanaDeserializer) EndLength() int {
	return len(s.data) - s.idx
}

func newIntLittleEndianBytes(b []byte) *big.Int {
	l := len(b)
	for i := 0; i < l/2; i++ {
		b[i], b[l-i-1] = b[l-i-1], b[i]
	}
	return big.NewInt(0).SetBytes(b)
}
