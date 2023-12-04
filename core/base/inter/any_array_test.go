package inter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type stringArray struct {
	AnyArray[string]
}

func TestStrArray(t *testing.T) {
	arr2 := stringArray{AnyArray: []string{}}
	require.Equal(t, arr2.JsonString(), "[]")

	arr := stringArray{}
	require.Equal(t, arr.JsonString(), "null")

	arr.Append("AA")
	arr.Append("BB")
	require.Equal(t, arr.Count(), 2)
	require.Equal(t, arr.ValueAt(1), "BB")

	arr.Append("CC")
	require.Equal(t, arr.Count(), 3)

	require.Equal(t, arr.Remove(1), "BB")
	require.Equal(t, arr.Count(), 2)

	arr.Append("DD")
	require.Equal(t, arr.JsonString(), `["AA","CC","DD"]`)

	arr.Append("CC")
	idx1 := arr.FirstIndexOf(func(elem string) bool { return elem == "CC" })
	require.Equal(t, idx1, 1)
	idx2 := arr.LastIndexOf(func(elem string) bool { return elem == "CC" })
	require.Equal(t, idx2, 3)
	idx3 := arr.FirstIndexOf(func(elem string) bool { return elem == "EE" })
	require.Equal(t, idx3, -1)
}

type intArray struct {
	AnyArray[int]
}

func TestIntArray(t *testing.T) {
	arr := intArray{}
	for i := 0; i < 10; i++ {
		arr.Append(i)
	}
	require.Equal(t, arr.Count(), 10)
	require.Equal(t, arr.JsonString(), `[0,1,2,3,4,5,6,7,8,9]`)

	idx := arr.LastIndexOf(func(elem int) bool { return elem == 3 })
	require.Equal(t, idx, 3)
}
