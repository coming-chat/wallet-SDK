package base

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type JsonObj struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (o *JsonObj) JsonString() (*OptionalString, error) {
	return JsonString(o)
}

func NewJsonObjWithJsonString(str string) (*JsonObj, error) {
	var o JsonObj
	err := FromJsonString(str, &o)
	return &o, err
}

func NewJsonObjArrayWithJsonString(str string) (*AnyArray, error) {
	var o []*JsonObj
	err := FromJsonString(str, &o)
	arr := make([]any, len(o))
	for i, v := range o {
		arr[i] = v
	}
	return &AnyArray{Values: arr}, err
}

func (o *JsonObj) AsAny() *Any {
	return &Any{o}
}

func TestJsonObject(t *testing.T) {
	o1 := JsonObj{Name: "Zhi", Age: 20}
	jsonStr, err := o1.JsonString()
	require.Nil(t, err)
	t.Log(jsonStr.Value)

	o2, err := NewJsonObjWithJsonString(jsonStr.Value)
	require.Nil(t, err)
	t.Log(o2)
}

func TestJsonForAny(t *testing.T) {
	o1 := JsonObj{Name: "Zhi", Age: 20}
	a1 := o1.AsAny()

	jsonStrO1, err := o1.JsonString()
	require.Nil(t, err)
	jsonStrA1, err := a1.JsonString()
	require.Nil(t, err)

	require.Equal(t, jsonStrO1, jsonStrA1)
	t.Log(jsonStrO1.Value)

	// ======================
	o2 := JsonObj{Name: "A22", Age: 17}

	arr1 := AnyArray{Values: []any{a1, o2}} // a1 is Any, o2 is JsonObj
	jsonStrArr1, err := arr1.JsonString()
	require.Nil(t, err)
	arr2 := []JsonObj{o1, o2}
	jsonStrArr2, err := JsonString(arr2)
	require.Nil(t, err)

	require.Equal(t, jsonStrArr1, jsonStrArr2)
	t.Log(jsonStrArr1.Value)

	// ======================= new array
	objArray, err := NewJsonObjArrayWithJsonString(jsonStrArr1.Value)
	require.Nil(t, err)
	t.Log(objArray.Values...)
}

func TestJsonForNestAny(t *testing.T) {
	o1 := JsonObj{Name: "Zhi", Age: 20}

	a1 := Any{Value: o1}
	a2 := Any{Value: a1}
	a3 := Any{Value: a2}

	jsonStr, err := a3.JsonString()
	require.Nil(t, err)
	t.Log(jsonStr.Value)

	// ===============
	arr1 := AnyArray{Values: []any{a1}}
	arr2 := AnyArray{Values: []any{arr1}}
	arr3 := AnyArray{Values: []any{arr2}}
	jsonStrArr, err := arr3.JsonString()
	require.Nil(t, err)
	t.Log(jsonStrArr.Value)
}
