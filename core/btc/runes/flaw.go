package runes

type Flaw byte

const (
	EdictOutput Flaw = iota
	EdictRuneId
	InvalidScript
	Opcode
	SupplyOverflow
	TrailingIntegers
	TruncatedField
	UnrecognizedEvenTag
	UnrecognizedFlag
	Varint
)

func (f Flaw) Flag() uint32 {
	return 1 << uint32(f)
}
