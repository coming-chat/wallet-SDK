package runes

import "math/big"

type RuneId struct {
	Block uint64
	Tx    uint32
}

func (r *RuneId) Next(block big.Int, tx big.Int) *RuneId {
	if block.BitLen() > 64 || tx.BitLen() > 32 {
		return nil
	}
	newBlock := r.Block + block.Uint64()
	var newTx uint32
	if newBlock == 0 {
		newTx = r.Tx + uint32(tx.Uint64())
	} else {
		newTx = uint32(tx.Uint64())
	}
	return &RuneId{newBlock, newTx}
}
