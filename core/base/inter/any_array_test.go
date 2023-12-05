package inter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// MARK - StringArray

type stringArray struct {
	AnyArray[string]
}

func NewStringArray() *stringArray {
	return &stringArray{[]string{}}
}

func NewStringArrayWithItem(elem string) *stringArray {
	return &stringArray{[]string{elem}}
}

func (a stringArray) Contains(value string) bool {
	idx := FirstIndexOf(a.AnyArray, func(elem string) bool { return elem == value })
	return idx != -1
}

func TestStringArray(t *testing.T) {
	arr2 := NewStringArray()
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
	idx1 := FirstIndexOf(arr.AnyArray, func(elem string) bool { return elem == "CC" })
	require.Equal(t, idx1, 1)
	idx2 := LastIndexOf(arr.AnyArray, func(elem string) bool { return elem == "CC" })
	require.Equal(t, idx2, 3)
	idx3 := FirstIndexOf(arr.AnyArray, func(elem string) bool { return elem == "EE" })
	require.Equal(t, idx3, -1)
	require.False(t, arr.Contains("EE"))
}

func TestAnyArray_Unmarshal(t *testing.T) {
	jsonStr := `
	["aa", "bb","cc",   "", "d d"]   `

	var arr stringArray
	err := json.Unmarshal([]byte(jsonStr), &arr)
	require.Nil(t, err)

	require.Equal(t, arr.JsonString(), `["aa","bb","cc","","d d"]`)
	require.True(t, arr.Contains("aa"))
	require.True(t, arr.Contains("d d"))
	require.False(t, arr.Contains("xx"))
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

	idx := LastIndexOf(arr.AnyArray, func(elem int) bool { return elem == 3 })
	require.Equal(t, idx, 3)
}
