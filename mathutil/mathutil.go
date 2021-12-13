package mathutil

import (
	"math"
	"math/big"
)

// Pow ...
func Pow(i, e *big.Int) *big.Int {
	return new(big.Int).Exp(i, e, nil)
}

// FromTwos ...
func FromTwos(value *big.Int, width int) *big.Int {
	a := inotn(value, width)
	b := new(big.Int).Add(a, big.NewInt(1))
	return new(big.Int).Neg(b)
}

// ToTwos ...
func ToTwos(value *big.Int, width int) *big.Int {
	if value.Cmp(big.NewInt(0)) == -1 {
		a := Abs(value)
		b := inotn(a, width)
		c := iaddn(b, 1)
		return c
	}

	return Clone(value)
}

// Abs ...
func Abs(value *big.Int) *big.Int {
	ret := Clone(value)
	return ret.Neg(value)
}

// Clone ...
func Clone(value *big.Int) *big.Int {
	ret := new(big.Int)
	ret.SetBits(value.Bits())
	return value
}

// Iaddn ...
func Iaddn(value *big.Int, num int) *big.Int {
	return iaddn(value, num)
}

func iaddn(value *big.Int, num int) *big.Int {
	return new(big.Int).Add(value, big.NewInt(int64(num)))
}

// Inotn ...
func Inotn(n *big.Int, width int) *big.Int {
	return inotn(n, width)
}

// inotn ...
func inotn(n *big.Int, width int) *big.Int {
	isNegative := n.Cmp(big.NewInt(0)) == -1
	var bytesNeeded = int(math.Ceil(float64(width) / float64(26)))
	var bitsLeft = width % 26
	var length = BitLen(n)

	words := make([]big.Word, bytesNeeded)
	copy(words[:], n.Bits())

	// extend the buffer with leading zeroes
	for length < bytesNeeded {
		words[length] = 0
		length = length + 1
	}

	if bitsLeft > 0 {
		bytesNeeded--
	}

	var i int
	// handle complete words
	for i = 0; i < bytesNeeded; i++ {
		words[i] = big.Word(^int(words[i]) & 0x3ffffff)
	}

	// handle the residue
	if bitsLeft > 0 {
		words[i] = big.Word(^int(words[i]) & (0x3ffffff >> (26 - uint(bitsLeft))))
	}

	words, _ = stripZerosFromWords(words, isNegative)

	ret := new(big.Int)
	ret.SetBits(words)
	return ret
}

func stripZerosFromWords(words []big.Word, isNegative bool) ([]big.Word, bool) {
	var i int
	for len(words) > 1 && words[len(words)-1] == 0 {
		i++
	}

	words = words[:len(words)-i]

	return words, normSign(words, isNegative)
}

func strip(value *big.Int) *big.Int {
	words := value.Bits()
	var i int
	for len(words) > 1 && words[len(words)-1] == 0 {
		i++
	}

	words = words[:len(words)-i]

	ret := new(big.Int)
	ret.SetBits(words)
	return ret
}

func normSign(words []big.Word, isNegative bool) bool {
	// -0 = 0
	if len(words) == 1 && words[0] == 0 {
		return false
	}

	return isNegative
}

// ToUint8Slice ...
func ToUint8Slice(value *big.Int, isLittleEndian bool, length int) []uint8 {
	byteLength := int(math.Ceil(float64(BitLen(value)) / float64(8)))
	var reqLength int
	if length > 0 {
		reqLength = length
	} else {
		reqLength = int(math.Max(float64(1), float64(byteLength)))
	}

	if byteLength > reqLength {
		panic("byte array longer than desired length")
	}
	if reqLength <= 0 {
		panic("requested array length <= 0")
	}

	value = strip(value)
	res := make([]uint8, reqLength)

	q := new(big.Int)
	q.SetBits(value.Bits())
	if isLittleEndian {
		var i int
		for ; !isZero(q); i++ {
			b := andln(q, 0xff)
			q = iushrn(q, 8, -1, false)

			res[i] = uint8(b)
		}

		for ; i < reqLength; i++ {
			res[i] = 0
		}
	} else {
		for i := 0; i < reqLength-byteLength; i++ {
			res[i] = 0
		}

		for i := 0; !isZero(q); i++ {
			b := andln(q, 0xff)
			q = iushrn(q, 8, -1, false)

			res[reqLength-i-1] = uint8(b)
		}
	}

	return res
}

func isZero(value *big.Int) bool {
	if len(value.Bits()) == 1 && value.Bits()[0] == 0 {
		return true
	} else if value.Cmp(big.NewInt(0)) == 0 {
		return true
	}

	return false
}

// Andln ANDs first word and num
func Andln(value *big.Int, num int) int {
	return andln(value, num)
}

func andln(value *big.Int, num int) int {
	words := value.Bits()
	return int(words[0]) & num
}

// Iushrn ...
func Iushrn(value *big.Int, bits uint, hint int, extended bool) *big.Int {
	return iushrn(value, bits, hint, extended)
}

// Shift-right in-place
// NOTE: `hint` is a lowest bit before trailing zeroes
// NOTE: if `extended` is present - it will be filled with destroyed bits
func iushrn(value *big.Int, bits uint, hint int, extended bool) *big.Int {
	var h int
	if hint != -1 {
		h = (hint - (hint % 26)) / 26
	} else {
		h = 0
	}

	r := bits % 26
	s := int(math.Min(float64((bits-r)/26), float64(len(value.Bits()))))
	mask := 0x3ffffff ^ ((0x3ffffff >> r) << r)
	maskedWords := make([]big.Word, len(value.Bits()))

	h = h - s
	h = int(math.Max(float64(0), float64(h)))

	words := value.Bits()
	// Extended mode, copy masked part
	if extended {
		for i := 0; i < s; i++ {
			maskedWords[i] = words[i]
		}
		maskedWords = maskedWords[:s]
	}

	var l = len(value.Bits())
	if s == 0 {
		// No-op, we should not move anything at all
	} else if len(value.Bits()) > s {
		l = l - s
		for i := 0; i < l; i++ {
			words[i] = words[i+s]
		}
	} else {
		words[0] = 0
		l = 1
	}

	var carry uint
	for i := l - 1; i >= 0 && (carry != 0 || i >= h); i-- {
		word := words[i]
		words[i] = big.Word((carry << (26 - r)) | (uint(word) >> r))
		carry = uint(word) & uint(mask)
	}

	// Push carried bits as a mask
	if len(maskedWords) > 0 && carry != 0 {
		j := len(maskedWords)
		// expand slice
		if len(maskedWords) <= j {
			tmp := make([]big.Word, j+1)
			copy(tmp[:], maskedWords)
			maskedWords = tmp
		}
		maskedWords[j] = big.Word(carry)
		maskedWords = maskedWords[:j+1]
	}

	if l == 0 {
		words[0] = 0
		l = 1
	}

	words = words[:l]
	ret := new(big.Int)
	ret.SetBits(words)
	ret = strip(ret)
	return ret
}

// BitLen ...
func BitLen(value *big.Int) int {
	bits := value.Bits()
	var w big.Word
	if len(bits) > 0 {
		w = bits[len(bits)-1]
	}
	hi := CountBits(int(w))
	return (len(bits)-1)*26 + hi
}

// CountBits ...
func CountBits(w int) int {
	t := w
	r := 0
	if t >= 0x1000 {
		r += 13
		t = t >> 13
	}
	if t >= 0x40 {
		r += 7
		t = t >> 7
	}
	if t >= 0x8 {
		r += 4
		t = t >> 4
	}
	if t >= 0x02 {
		r += 2
		t = t >> 2
	}

	return r + t
}
