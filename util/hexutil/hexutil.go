package hexutil

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/coming-chat/wallet-SDK/util/mathutil"
)

var (
	// Prefix hex prefix
	Prefix = "0x"

	// ErrInvalidHex is error for invalid hex
	ErrInvalidHex = errors.New("invalid hex")
)

// HasPrefix tests for the existence of a `0x` prefix. Checks for a valid hex input value and if the start matched `0x`.
func HasPrefix(hexStr string) bool {
	if len(hexStr) < 2 {
		return false
	}

	return len(hexStr) >= 2 && hexStr[0:2] == "0x"
}

// ValidHex tests that the value is a valid string. The `0x` prefix is optional in the test. Empty string returns false.
func ValidHex(hexStr string) bool {
	if len(hexStr) == 0 {
		return false
	}

	regex := regexp.MustCompile(`^(0x)?[\da-fA-F]*$`)
	return regex.Match([]byte(hexStr))
}

// AddPrefix adds the `0x` prefix to string values.
// Returns a `0x` prefixed string from the input value. If the input is already prefixed, it is returned unchanged. Adds extra 0 when `length % 2 == 1`.
func AddPrefix(hexStr string) string {
	if HasPrefix(hexStr) {
		return hexStr
	}

	if len(hexStr)%2 == 1 {
		hexStr = fmt.Sprintf("0%s", hexStr)
	}

	return fmt.Sprintf("%s%s", Prefix, hexStr)
}

// StripPrefix strips any leading `0x` prefix. Tests for the existence of a `0x` prefix, and returns the value without the prefix. Un-prefixed values are returned as-is.
func StripPrefix(hexStr string) string {
	return strings.TrimPrefix(hexStr, Prefix)
}

// HexFixLength shifts a hex string to a specific bitLength.
// Returns a `0x` prefixed string with the specified number of bits contained in the return value. (If bitLength is -1, length checking is not done). Values with more bits are trimmed to the specified length. Input values with less bits are returned as-is by default. When `withPadding` is set, shorter values are padded with `0`.
func HexFixLength(hexStr string, bitLength int, withPadding bool) string {
	strLen := int(math.Ceil(float64(bitLength) / float64(4)))
	hexLen := strLen + 2

	if bitLength == -1 || len(hexStr) == hexLen ||
		(!withPadding && len(hexStr) < hexLen) {
		return AddPrefix(hexStr)
	}

	strippedHexStr := StripPrefix(hexStr)
	strippedHexLen := len(strippedHexStr)

	if len(hexStr) > hexLen {
		return AddPrefix(
			strippedHexStr[strippedHexLen-strLen : strippedHexLen],
		)
	}

	paddedHexStr := fmt.Sprintf("%s%s", strings.Repeat("0", strLen), strippedHexStr)
	return AddPrefix(
		paddedHexStr[len(paddedHexStr)-strLen:],
	)
}

// ToBN creates a math/big big number from a hex string.
func ToBN(hexStr string, isLittleEndian bool, isNegative bool) (*big.Int, error) {
	i := new(big.Int)
	hx := StripPrefix(hexStr)

	if hx == "" {
		return big.NewInt(0), nil
	}

	if isLittleEndian {
		hx = Reverse(hx)
	}

	if _, ok := i.SetString(hx, 16); !ok {
		return nil, errors.New("could not decode to big.Int")
	}

	// NOTE: fromTwos takes as parameter the number of bits,
	// which is the hex length multiplied by 4.
	if isNegative {
		return mathutil.FromTwos(i, i.BitLen()), nil
	}

	return i, nil
}

// ToUint8Slice creates a uint8 array from a hex string. empty inputs returns an empty array result. Hex input values return the actual bytes value converted to a uint8. Anything that is not a hex string (including the `0x` prefix) returns an error.
func ToUint8Slice(hexStr string, bitLength int) ([]uint8, error) {
	if hexStr == "" {
		return []uint8{}, nil
	}

	if !ValidHex(hexStr) {
		return nil, ErrInvalidHex
	}

	value := StripPrefix(hexStr)
	valLength := len(value) / 2
	var bufLength int
	if bitLength == -1 {
		bufLength = int(math.Ceil(float64(valLength)))
	} else {
		bufLength = int(math.Ceil(float64(bitLength) / float64(8)))
	}

	result := make([]uint8, bufLength)
	offset := int(math.Max(float64(0), float64(bufLength-valLength)))
	for index := 0; index < bufLength; index++ {
		n := (index * 2) + 2
		if n > len(value) {
			continue
		}
		s := value[index*2 : n]
		v, err := strconv.ParseInt(s, 16, 64)
		if err != nil {
			return nil, err
		}

		result[index+offset] = uint8(v)
	}

	return result, nil
}

// Reverse reverses a hex string
func Reverse(s string) string {
	s = StripPrefix(s)
	regex := regexp.MustCompile(`.{1,2}`)
	in := regex.FindAllString(s, -1)
	var out []string
	for i := range in {
		v := in[len(in)-1-i]
		out = append(out, v)
	}
	return strings.Join(out, "")
}
