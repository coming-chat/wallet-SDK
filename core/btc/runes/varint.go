package runes

import (
	"errors"
	"math/big"
)

var (
	ErrOverlong     = errors.New("too long")
	ErrOverflow     = errors.New("overflow")
	ErrUnterminated = errors.New("unterminated")
)

func Decode(buffer []byte) (*big.Int, *int) {
	u128, i, err := tryDecode(buffer)
	if err != nil {
		return nil, nil
	}
	return u128, &i
}

func tryDecode(buffer []byte) (*big.Int, int, error) {
	n := big.NewInt(0)

	for i, b := range buffer {
		if i > 18 {
			return nil, 0, ErrOverlong
		}

		value := new(big.Int).And(big.NewInt(int64(b)), big.NewInt(0b0111_1111))

		if i == 18 && new(big.Int).And(value, big.NewInt(0b0111_1100)).Cmp(big.NewInt(0)) != 0 {
			return nil, 0, ErrOverflow
		}

		n = new(big.Int).Or(n, new(big.Int).Lsh(value, uint(7*i)))

		if b&0b1000_0000 == 0 {
			return n, i + 1, nil
		}
	}
	return nil, 0, ErrUnterminated
}
