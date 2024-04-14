package runes

import "math/big"

type Etching struct {
	Divisibility *uint8   `json:"divisibility"`
	Premine      *big.Int `json:"premine"`
	Rune         *Rune    `json:"rune"`
	Spacers      *uint32  `json:"spacers"`
	Symbol       *rune    `json:"symbol"`
	Terms        *Terms   `json:"terms"`
	Turbo        bool     `json:"turbo"`
}

const (
	MaxDivisibility uint8  = 38
	MaxSpacers      uint32 = 0b00000111_11111111_11111111_11111111
)

func (e *Etching) Supply() *big.Int {
	premine := big.NewInt(0)
	if e.Premine != nil {
		premine = e.Premine
	}
	cap := big.NewInt(0)
	if e.Terms != nil && e.Terms.Cap != nil {
		cap = e.Terms.Cap
	}
	amount := big.NewInt(0)
	if e.Terms != nil && e.Terms.Amount != nil {
		amount = e.Terms.Amount
	}
	mul := new(big.Int).Mul(cap, amount)
	if mul.BitLen() > 128 {
		return nil
	}
	supply := new(big.Int).Add(premine, mul)
	if supply.BitLen() > 128 {
		return nil
	}
	return supply
}
