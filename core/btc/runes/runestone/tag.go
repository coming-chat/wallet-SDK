package runestone

import (
	"math/big"
)

type Tag uint64

const (
	Body Tag = iota
	Divisibility
	Flags
	Spacers
	Rune
	Symbol
	Premine
	_
	Cap
	_
	Amount
	_
	HeightStart
	_
	HeightEnd
	_
	OffsetStart
	_
	OffsetEnd
	_
	Mint
	_
	Pointer
	// unused
	Cenotaph = iota + 103
	// unused
	Nop
)

func (t Tag) ToBigInt() *big.Int {
	return new(big.Int).SetUint64(uint64(t))
}

func Take[T any](tag Tag, fields map[string][]big.Int, n int, with func(...big.Int) *T) *T {
	field, ok := fields[tag.ToBigInt().String()]
	if !ok {
		return nil
	}
	values := make([]big.Int, n)

	for i := range values {
		if i >= len(field) {
			return nil
		}
		values[i] = field[i]
	}

	value := with(values...)
	if value == nil {
		return nil
	}

	field = field[n:]
	if len(field) == 0 {
		delete(fields, tag.ToBigInt().String())
	}
	return value
}
