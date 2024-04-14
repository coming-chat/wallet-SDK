package runes

type Cenotaph struct {
	Etching *Rune
	Flaws   uint32
	Mint    *RuneId
}

func (c Cenotaph) artifact() {
}
