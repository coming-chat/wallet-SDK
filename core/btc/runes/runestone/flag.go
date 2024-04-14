package runestone

import "math/big"

type Flag uint

const (
	Etching Flag = iota
	Terms
	Turbo
	// unused
	FlagCenotaph = iota + 124
)

func (f Flag) mask() *big.Int {
	return new(big.Int).Lsh(big.NewInt(1), uint(f))
}

func (f Flag) Take(flags *big.Int) bool {
	mask := f.mask()
	set := new(big.Int).And(flags, mask).Cmp(big.NewInt(0)) != 0
	flags.Set(new(big.Int).AndNot(flags, mask))
	return set
}
