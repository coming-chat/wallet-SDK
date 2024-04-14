package runes

import (
	"github.com/btcsuite/btcd/wire"
	"math/big"
)

type Edict struct {
	Id     RuneId
	Amount big.Int
	Output uint32
}

func NewEdictFromIntegers(tx *wire.MsgTx, id RuneId, amount big.Int, output big.Int) *Edict {
	if output.BitLen() > 32 {
		return nil
	}
	if int(output.Uint64()) > len(tx.TxOut) {
		return nil
	}
	return &Edict{
		Id:     id,
		Amount: amount,
		Output: uint32(output.Uint64()),
	}
}
