package u8util

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestConcat(t *testing.T) {
	for i, tt := range []struct {
		in  [][]uint8
		out []uint8
	}{
		{
			[][]uint8{
				{0x1, 0x2, 0x3},
				{0x4, 0x5},
			},
			[]uint8{0x1, 0x2, 0x3, 0x4, 0x5},
		},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := Concat(tt.in...)
			if !reflect.DeepEqual(result, tt.out) {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestFixLength(t *testing.T) {
	type input struct {
		value     []uint8
		bitLength int
		atStart   bool
	}
	for i, tt := range []struct {
		in  input
		out []uint8
	}{
		{
			input{[]uint8{0x12, 0x34, 0x56, 0x78}, -1, false},
			[]uint8{0x12, 0x34, 0x56, 0x78},
		},
		{
			input{[]uint8{0x12, 0x34, 0x56, 0x78}, 32, false},
			[]uint8{0x12, 0x34, 0x56, 0x78},
		},
		{
			input{[]uint8{0x12, 0x34, 0x56, 0x78}, 16, false},
			[]uint8{0x12, 0x34},
		},
		{
			input{[]uint8{0x12, 0x34}, 32, false},
			[]uint8{0, 0, 0x12, 0x34},
		},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := FixLength(tt.in.value, tt.in.bitLength, tt.in.atStart)
			if !reflect.DeepEqual(result, tt.out) {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestToString(t *testing.T) {
	for i, tt := range []struct {
		in  []uint8
		out string
	}{
		{[]uint8{0x68, 0x65, 0x6c, 0x6c, 0x6f}, "hello"},
		{nil, ""},
		{[]uint8{208, 159, 209, 128, 208, 184, 208, 178, 208, 181, 209, 130, 44, 32, 208, 188, 208, 184, 209, 128, 33}, "Привет, мир!"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := ToString(tt.in)
			if result != tt.out {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestToHex(t *testing.T) {
	type input struct {
		value      []uint8
		bitLength  int
		isPrefixed bool
	}
	for i, tt := range []struct {
		in  input
		out string
	}{
		{input{[]uint8{0x68, 0x65, 0x6c, 0x6c, 0xf}, -1, true}, "0x68656c6c0f"},
		{input{[]uint8{}, -1, true}, "0x"},
		{input{[]uint8{}, -1, false}, ""},
		{input{[]uint8{0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, -1, true}, "0x0001000000000000"},
		{input{[]uint8{128, 0, 10, 11, 12, 13}, 32, true}, "0x8000…0c0d"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := ToHex(tt.in.value, tt.in.bitLength, tt.in.isPrefixed)
			if result != tt.out {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestFromHex(t *testing.T) {
	for i, tt := range []struct {
		in  string
		out []uint8
	}{
		{"0x68656c6c0f", []uint8{0x68, 0x65, 0x6c, 0x6c, 0xf}},
		{"68656c6c0f", []uint8{0x68, 0x65, 0x6c, 0x6c, 0xf}},
		{"", []uint8{}},
		{"0x", []uint8{}},
		{"0001000000000000", []uint8{0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
		{"0x0001000000000000", []uint8{0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := FromHex(tt.in)
			if !reflect.DeepEqual(result, tt.out) {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestToBN(t *testing.T) {
	type input struct {
		value          []uint8
		isLittleEndian bool
	}
	for i, tt := range []struct {
		in  input
		out *big.Int
	}{
		{input{[]uint8{0x12, 0x34}, false}, big.NewInt(4660)},
		{input{[]uint8{0x12, 0x34}, true}, big.NewInt(13330)},
		//{input{[]uint8{0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, true}, big.NewInt(256)},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := ToBN(tt.in.value, tt.in.isLittleEndian)
			if result.String() != tt.out.String() {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}
