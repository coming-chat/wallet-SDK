package crypto

import "crypto/sha256"

// marshal and unmarshal to/from hex

// NewHash ...
// TODO: redeclared. which to use...
//func NewHash(input []byte) *Hash {
//h := new(Hash)
//copy(h[:], input)
//return h
//}

// Value ...
func (h Hash) Value() [sha256.Size]uint8 {
	return h
}
