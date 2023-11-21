package inter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type stringArray struct {
	AnyArray[string]
}

func TestStrArray(t *testing.T) {
	arr := stringArray{}

	arr.Append("AA")
	arr.Append("BB")
	require.Equal(t, arr.Count(), 2)
	require.Equal(t, arr.ValueAt(1), "BB")

	arr.Append("CC")
	require.Equal(t, arr.Count(), 3)

	require.Equal(t, arr.Remove(1), "BB")
	require.Equal(t, arr.Count(), 2)

	arr.Append("DD")
	jjj, err := arr.JsonString()
	require.Nil(t, err)
	require.Equal(t, jjj.Value, `["AA","CC","DD"]`)
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
	jjj, err := arr.JsonString()
	require.Nil(t, err)
	require.Equal(t, jjj.Value, `[0,1,2,3,4,5,6,7,8,9]`)
}
