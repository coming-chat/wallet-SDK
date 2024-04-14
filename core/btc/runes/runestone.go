package runes

import (
	"bytes"
	"errors"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/wallet-SDK/core/btc/runes/runestone"
	"golang.org/x/exp/maps"
	"math/big"
	"slices"
)

const (
	scriptVersion = 0
	MagicNumber   = txscript.OP_13
)

var (
	ErrNotFoundPayload = errors.New("not found payload")
)

type Artifact interface {
	artifact()
}

type Runestone struct {
	Edicts  []Edict
	Etching *Etching
	Mint    *RuneId
	Pointer *uint32
}

func (r Runestone) artifact() {
}

func Decipher(transaction *wire.MsgTx) (Artifact, error) {
	pl, err := payload(transaction)
	if err != nil {
		return nil, err
	}
	var payloadScript []byte
	switch pl.(type) {
	case InvalidPayload:
		return Cenotaph{
			Flaws: pl.(InvalidPayload).Flag(),
		}, nil
	case ValidPayload:
		payloadScript = pl.(ValidPayload)
	}
	integers := integers(payloadScript)
	if integers == nil {
		return Cenotaph{
			Flaws: Varint.Flag(),
		}, nil
	}

	message := NewMessageFromIntegers(transaction, *integers)

	flags := runestone.Take[big.Int](runestone.Flags, message.Fields, 1, func(b ...big.Int) *big.Int {
		if len(b) != 1 {
			return nil
		}
		return &b[0]
	})
	if flags == nil {
		flags = big.NewInt(0)
	}

	var etching *Etching
	if runestone.Etching.Take(flags) {
		etching = &Etching{
			Divisibility: runestone.Take[byte](runestone.Divisibility, message.Fields, 1, func(b ...big.Int) *byte {
				if len(b) != 1 || b[0].BitLen() > 8 {
					return nil
				}

				divisibility := byte(b[0].Uint64())
				if divisibility <= MaxDivisibility {
					return &divisibility
				}
				return nil
			}),
			Premine: runestone.Take[big.Int](runestone.Premine, message.Fields, 1, func(b ...big.Int) *big.Int {
				if len(b) != 1 {
					return nil
				}
				return &b[0]
			}),
			Rune: runestone.Take[Rune](runestone.Rune, message.Fields, 1, func(b ...big.Int) *Rune {
				if len(b) != 1 {
					return nil
				}
				return &Rune{b[0]}
			}),
			Spacers: runestone.Take[uint32](runestone.Spacers, message.Fields, 1, func(b ...big.Int) *uint32 {
				if len(b) != 1 || b[0].BitLen() > 32 {
					return nil
				}
				spacers := uint32(b[0].Uint64())
				if spacers <= MaxSpacers {
					return &spacers
				}
				return nil
			}),
			Symbol: runestone.Take[rune](runestone.Symbol, message.Fields, 1, func(b ...big.Int) *rune {
				if len(b) != 1 || b[0].BitLen() > 32 {
					return nil
				}
				symbol := rune(b[0].Uint64())
				return &symbol
			}),
		}
		if runestone.Terms.Take(flags) {
			etching.Terms = &Terms{
				Cap: runestone.Take[big.Int](runestone.Cap, message.Fields, 1, func(b ...big.Int) *big.Int {
					if len(b) != 1 {
						return nil
					}
					return &b[0]
				}),
				Height: [2]*uint64{
					runestone.Take[uint64](runestone.HeightStart, message.Fields, 1, func(b ...big.Int) *uint64 {
						if len(b) != 1 || b[0].BitLen() > 64 {
							return nil
						}
						startHeight := b[0].Uint64()
						return &startHeight
					}),
					runestone.Take[uint64](runestone.HeightEnd, message.Fields, 1, func(b ...big.Int) *uint64 {
						if len(b) != 1 || b[0].BitLen() > 64 {
							return nil
						}
						endHeight := b[0].Uint64()
						return &endHeight
					}),
				},
				Amount: runestone.Take[big.Int](runestone.Amount, message.Fields, 1, func(b ...big.Int) *big.Int {
					if len(b) != 1 {
						return nil
					}
					return &b[0]
				}),
				Offset: [2]*uint64{
					runestone.Take[uint64](runestone.OffsetStart, message.Fields, 1, func(b ...big.Int) *uint64 {
						if len(b) != 1 || b[0].BitLen() > 64 {
							return nil
						}
						startOffset := b[0].Uint64()
						return &startOffset
					}),
					runestone.Take[uint64](runestone.OffsetEnd, message.Fields, 1, func(b ...big.Int) *uint64 {
						if len(b) != 1 || b[0].BitLen() > 64 {
							return nil
						}
						endOffset := b[0].Uint64()
						return &endOffset
					}),
				},
			}
		}

		etching.Turbo = runestone.Turbo.Take(flags)
	}

	mint := runestone.Take[RuneId](runestone.Mint, message.Fields, 2, func(b ...big.Int) *RuneId {
		if len(b) != 2 || b[0].BitLen() > 64 || b[1].BitLen() > 32 {
			return nil
		}
		return &RuneId{Block: b[0].Uint64(), Tx: uint32(b[0].Uint64())}
	})

	pointer := runestone.Take[uint32](runestone.Pointer, message.Fields, 1, func(b ...big.Int) *uint32 {
		if len(b) != 1 || b[0].BitLen() > 32 {
			return nil
		}
		pointer := uint32(b[0].Uint64())
		if uint64(pointer) < uint64(len(transaction.TxOut)) {
			return &pointer
		}
		return nil
	})

	if etching != nil && etching.Supply() == nil {
		message.Flaws |= SupplyOverflow.Flag()
	}

	if flags.Cmp(big.NewInt(0)) != 0 {
		message.Flaws |= UnrecognizedFlag.Flag()
	}

	if slices.IndexFunc(maps.Keys(message.Fields), func(s string) bool {
		tag, _ := new(big.Int).SetString(s, 10)
		return new(big.Int).Mod(tag, big.NewInt(2)).Cmp(big.NewInt(0)) == 0
	}) != -1 {
		message.Flaws |= UnrecognizedEvenTag.Flag()
	}

	if message.Flaws != 0 {
		cenotaph := Cenotaph{
			Flaws: message.Flaws,
			Mint:  mint,
		}
		if etching != nil {
			cenotaph.Etching = etching.Rune
		}
		return cenotaph, nil
	}

	return Runestone{
		Edicts:  message.Edicts,
		Etching: etching,
		Mint:    mint,
		Pointer: pointer,
	}, nil
}

func integers(payload []byte) *[]big.Int {
	var integers []big.Int
	i := 0
	for i < len(payload) {
		integer, length := Decode(payload[i:])
		if integer == nil {
			return nil
		}
		integers = append(integers, *integer)
		i += *length
	}
	return &integers
}

type Payload interface {
	payload()
}

type ValidPayload []byte

func (v ValidPayload) payload() {
}

type InvalidPayload = Flaw

func (i InvalidPayload) payload() {}

func payload(transaction *wire.MsgTx) (Payload, error) {
	for _, out := range transaction.TxOut {
		tokenizer := txscript.MakeScriptTokenizer(scriptVersion, out.PkScript)
		if !tokenizer.Next() || tokenizer.Err() != nil {
			continue
		}
		if tokenizer.Opcode() != txscript.OP_RETURN || !tokenizer.Next() || tokenizer.Err() != nil {
			continue
		}
		if tokenizer.Opcode() != MagicNumber {
			continue
		}
		payload := bytes.NewBuffer([]byte{})
		for tokenizer.Next() {
			switch {
			case tokenizer.Data() != nil:
				payload.Write(tokenizer.Data())
			default:
				return InvalidPayload(Opcode), nil
			}
		}
		if tokenizer.Err() != nil {
			return InvalidPayload(InvalidScript), nil
		}
		return ValidPayload(payload.Bytes()), nil
	}
	return nil, ErrNotFoundPayload
}
