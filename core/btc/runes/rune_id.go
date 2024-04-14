package runes

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

var (
	ErrIllegalRuneIdString = errors.New("illegal rune id string")
)

type RuneId struct {
	Block uint64 `json:"block"`
	Tx    uint32 `json:"tx"`
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

func (r *RuneId) String() string {
	return fmt.Sprintf("%d:%d", r.Block, r.Tx)
}

func NewRuneIdFromStr(runeIdStr string) (*RuneId, error) {
	runeIdData := strings.Split(runeIdStr, ":")
	if len(runeIdData) != 2 {
		return nil, ErrIllegalRuneIdString
	}
	block, err := strconv.ParseUint(runeIdData[0], 10, 64)
	if err != nil {
		return nil, err
	}
	tx, err := strconv.ParseUint(runeIdData[1], 10, 64)
	if err != nil {
		return nil, err
	}
	return &RuneId{block, uint32(tx)}, nil
}
