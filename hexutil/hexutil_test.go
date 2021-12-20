package hexutil

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestHasPrefix(t *testing.T) {
	for i, tt := range []struct {
		in  string
		out bool
	}{
		{"12", false},
		{"0x12", true},
		{"0x", true},
		{"0", false},
		{"", false},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := HasPrefix(tt.in)
			if result != tt.out {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestValidHex(t *testing.T) {
	for i, tt := range []struct {
		in  string
		out bool
	}{
		{"12", true},
		{"0x12", true},
		{"0x", true},
		{"0", true},
		{"", false},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := ValidHex(tt.in)
			if result != tt.out {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestAddPrefix(t *testing.T) {
	for i, tt := range []struct {
		in  string
		out string
	}{
		{"12", "0x12"},
		{"123", "0x0123"},
		{"0x123", "0x123"},
		{"", "0x"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := AddPrefix(tt.in)
			if result != tt.out {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestStripPrefix(t *testing.T) {
	for i, tt := range []struct {
		in  string
		out string
	}{
		{"0x123", "123"},
		{"123", "123"},
		{"0x", ""},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := StripPrefix(tt.in)
			if result != tt.out {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestHexFixLength(t *testing.T) {
	type input struct {
		hexStr      string
		bitLength   int
		withPadding bool
	}

	for i, tt := range []struct {
		in  input
		out string
	}{
		{input{"0x12", 16, false}, "0x12"},
		{input{"0x12", 16, true}, "0x0012"},
		{input{"0x0012", 8, false}, "0x12"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := HexFixLength(tt.in.hexStr, tt.in.bitLength, tt.in.withPadding)
			if result != tt.out {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestToBN(t *testing.T) {
	type input struct {
		hexStr       string
		littleEndian bool
		isNegative   bool
	}

	for i, tt := range []struct {
		in  input
		out *big.Int
	}{
		{input{"0x14", false, false}, big.NewInt(20)},
		{input{"14", false, false}, big.NewInt(20)},
		{input{"0x14", true, false}, big.NewInt(20)},
		{input{"14", true, false}, big.NewInt(20)},
		{input{"81", false, false}, big.NewInt(129)},
		{input{"0", false, false}, big.NewInt(0)},
		{input{"", false, false}, big.NewInt(0)},
		{input{"0x", true, false}, big.NewInt(0)},
		{input{"0x4500000000000000", true, false}, big.NewInt(69)},
		{input{"0x0000000000000100", false, false}, big.NewInt(256)},
		{input{"0x0001000000000000", true, false}, big.NewInt(256)},
		{input{"0x2efb", true, true}, big.NewInt(-1234)},
		{input{"0xfb2e", false, true}, big.NewInt(-1234)},
		{input{"0x0100", false, false}, big.NewInt(256)},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result, err := ToBN(tt.in.hexStr, tt.in.littleEndian, tt.in.isNegative)
			if err != nil {
				t.Fatal(err)
			}
			if result.String() != tt.out.String() {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestToUint8Slice(t *testing.T) {
	type input struct {
		hexStr    string
		bitLength int
	}

	for i, tt := range []struct {
		in  input
		out []uint8
	}{
		{input{"0x80001f", -1}, []uint8{0x80, 0x00, 0x1f}},
		{input{"0x80001f", 32}, []uint8{0x00, 0x80, 0x00, 0x1f}},
		{input{"0x80000a", -1}, []uint8{128, 0, 10}},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result, err := ToUint8Slice(tt.in.hexStr, tt.in.bitLength)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result, tt.out) {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}
