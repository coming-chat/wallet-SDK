package runes

import (
	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/wallet-SDK/core/btc/runes/runestone"
	"math/big"
)

type Message struct {
	Flaws  uint32
	Edicts []Edict
	Fields map[string][]big.Int
}

func NewMessageFromIntegers(tx *wire.MsgTx, payload []big.Int) *Message {
	var edicts []Edict
	fields := make(map[string][]big.Int)
	flaws := uint32(0)

	for i := 0; i < len(payload); i += 2 {
		tag := payload[i]
		if tag.Cmp(runestone.Body.ToBigInt()) == 0 {
			id := RuneId{}
			for j := i + 1; j < len(payload); j += 4 {
				if len(payload[j:]) < 4 {
					flaws |= TrailingIntegers.Flag()
					break
				}

				next := id.Next(payload[j], payload[j+1])
				if next == nil {
					flaws |= EdictRuneId.Flag()
					break
				}

				edict := NewEdictFromIntegers(tx, *next, payload[j+2], payload[j+3])
				if edict == nil {
					flaws |= EdictOutput.Flag()
					break
				}

				id = *next
				edicts = append(edicts, *edict)
			}
			break
		}

		if i+1 >= len(payload) {
			flaws |= TruncatedField.Flag()
			break
		}
		value := payload[i+1]

		fields[tag.String()] = append(fields[tag.String()], value)
	}
	return &Message{
		flaws,
		edicts,
		fields,
	}
}
