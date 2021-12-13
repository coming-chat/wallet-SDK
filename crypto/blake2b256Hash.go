package crypto

import "golang.org/x/crypto/blake2b"

// NewBlake2b256Hash ...
// TODO: redeclared, which to use...
//func NewBlake2b256Hash(input []byte) *Blake2b256Hash {
//b := New(Blake2b256Hash)
//copy(b[:], input)
//return b
//}

// Value ...
func (b Blake2b256Hash) Value() [blake2b.Size256]uint8 {
	return b
}
