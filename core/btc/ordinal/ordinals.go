package ordinal

import (
	"bytes"
	"errors"
	"github.com/btcsuite/btcd/txscript"
)

type Ord struct {
	Type        string
	ContentType string
	Content     []byte
}

func DecodeOrdFromWitness(witness []byte) (*Ord, error) {
	const scriptVersion = 0
	tokenizer := txscript.MakeScriptTokenizer(scriptVersion, witness)
	for tokenizer.Next() {
		if tokenizer.Opcode() != txscript.OP_FALSE || tokenizer.Done() {
			continue
		}
		internalTokenizer := tokenizer
		if !internalTokenizer.Next() {
			continue
		}
		if internalTokenizer.Opcode() != txscript.OP_IF {
			continue
		}
		if !internalTokenizer.Next() {
			continue
		}
		if !bytes.Equal([]byte("ord"), internalTokenizer.Data()) {
			continue
		}
		if !internalTokenizer.Next() {
			continue
		}
		if !bytes.Equal([]byte{1}, internalTokenizer.Data()) {
			continue
		}
		if !internalTokenizer.Next() {
			continue
		}
		contentType := string(internalTokenizer.Data())
		if !internalTokenizer.Next() {
			continue
		}
		if internalTokenizer.Opcode() != txscript.OP_0 {
			continue
		}
		if !internalTokenizer.Next() {
			continue
		}
		content := internalTokenizer.Data()
		return &Ord{
			Type:        "ord",
			ContentType: contentType,
			Content:     content,
		}, nil
	}
	if tokenizer.Err() != nil {
		return nil, tokenizer.Err()
	}
	return nil, errors.New("not found ord")
}
