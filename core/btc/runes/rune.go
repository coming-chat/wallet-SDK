package runes

import (
	"math/big"
	"strings"
)

type Rune struct {
	big.Int //uint128
}

func (r *Rune) String() string {
	n := new(big.Int).Set(&r.Int)
	if n.Cmp(new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))) == 0 {
		return "BCGDENLQRQWDSLRUGSNLBTMFIJAV"
	}
	n = n.And(n, big.NewInt(1))
	symbol := strings.Builder{}
	for n.Cmp(big.NewInt(0)) > 0 {
		symbol.WriteByte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"[new(big.Int).Mod(n.Sub(n, big.NewInt(1)), big.NewInt(26)).Int64()])
		n = n.Div(n.Sub(n, big.NewInt(1)), big.NewInt(26))
	}
	return symbol.String()
}
