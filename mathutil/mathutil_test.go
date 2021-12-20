package mathutil

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestPow(t *testing.T) {
	type input struct {
		i *big.Int
		e *big.Int
	}
	for i, tt := range []struct {
		in  input
		out *big.Int
	}{
		{input{big.NewInt(16), big.NewInt(2)}, big.NewInt(256)},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := Pow(tt.in.i, tt.in.e)
			if result.String() != tt.out.String() {
				t.Fatalf("want %v; got %v", tt.out, result)
			}
		})
	}
}

func TestToUint8Slice(t *testing.T) {
	for i, tt := range []struct {
		val    *big.Int
		le     bool
		length int
		ret    []uint8
	}{
		{big.NewInt(-1234), false, -1, []uint8{4, 210}},
		{big.NewInt(-1234), true, -1, []uint8{210, 4}},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := ToUint8Slice(tt.val, tt.le, tt.length)
			if !reflect.DeepEqual(result, tt.ret) {
				t.Fatalf("want %v; got %v", tt.ret, result)
			}
		})
	}
}
