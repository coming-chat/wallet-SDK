package base

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type AnyPerson struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (o *AnyPerson) AsAny() *Any {
	return &Any{o}
}

func AsAnyPerson(a *Any) *AnyPerson {
	if res, ok := a.Value.(*AnyPerson); ok {
		return res
	}
	return nil
}

func TestAny(t *testing.T) {
	a := Any{}

	a.SetBool(true)
	require.Equal(t, a.GetBool(), true)

	a.SetInt(1234)
	require.Equal(t, a.GetInt(), 1234)

	a.SetString("qwer")
	require.Equal(t, a.GetString(), "qwer")
}

func TestAnyArray(t *testing.T) {
	arr := AnyArray{}

	a1 := &Any{true}
	arr.Append(a1)
	require.Equal(t, arr.Count(), 1)

	a2 := &Any{int(123)}
	arr.Append(a2)
	require.Equal(t, arr.ValueAt(0), a1)
	require.Equal(t, arr.ValueAt(1), a2)

	a3 := &Any{"abc"}
	arr.SetValue(a3, 1)
	require.Equal(t, arr.ValueAt(1), a3)
	require.Equal(t, arr.JsonString(), `[true,"abc"]`)
	require.Equal(t, arr.Count(), 2)

	a4 := &Any{uint16(456)}
	arr.Append(a4)
	require.Equal(t, arr.Count(), 3)
	arr.Remove(0)
	require.Equal(t, arr.Count(), 2)
	require.Equal(t, arr.JsonString(), `["abc",456]`)

	ap := &AnyPerson{Name: "GGG", Age: 22}
	arr.Append(ap.AsAny())
	t.Log(arr.JsonString())
	require.NotNil(t, AsAnyPerson(arr.ValueAt(2)))
	require.Nil(t, AsAnyPerson(arr.ValueAt(0)))
}

func TestAnyMap(t *testing.T) {
	mp := NewAnyMap()

	require.Equal(t, mp.JsonString(), "{}")

	a1 := &Any{true}
	mp.SetValue(a1, "bbb")
	require.Equal(t, mp.Keys().Count(), 1)
	require.Equal(t, mp.Keys().JsonString(), `["bbb"]`)
	require.Equal(t, mp.ValueOf("bbb"), a1)

	a2 := &Any{"abcd"}
	mp.SetValue(a2, "alpha")
	t.Log(mp.JsonString())
	require.Equal(t, mp.Contains("alpha"), true)
	require.Equal(t, mp.Contains("beta"), false)

	require.Nil(t, mp.Remove("notkey"))
	require.Equal(t, mp.Remove("bbb"), a1)

	ap := &AnyPerson{Name: "GGG"}
	mp.SetValue(ap.AsAny(), "person")
	t.Log(mp.JsonString())
	p := AsAnyPerson(mp.ValueOf("person"))
	require.Equal(t, p.Name, "GGG")
}

func TestBigInt(t *testing.T) {
	number := "99999999999999999999999999999999999999999999999999999999999999999999999999999"
	a := NewAny()
	a.SetBigInt(NewBigIntFromString(number, 10))
	require.Equal(t, number, a.GetBigInt().String())
}
