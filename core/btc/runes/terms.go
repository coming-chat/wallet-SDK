package runes

import "math/big"

type Terms struct {
	Amount *big.Int
	Cap    *big.Int
	Height [2]*uint64
	Offset [2]*uint64
}
