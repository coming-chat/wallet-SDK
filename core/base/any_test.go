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
	require.Equal(t, arr.ValueOf(0), a1)
	require.Equal(t, arr.ValueOf(1), a2)

	a3 := &Any{"abc"}
	arr.SetValue(a3, 1)
	require.Equal(t, arr.ValueOf(1), a3)
	require.Equal(t, arr.String(), `[true,"abc"]`)
	require.Equal(t, arr.Count(), 2)

	a4 := &Any{uint16(456)}
	arr.Append(a4)
	require.Equal(t, arr.Count(), 3)
	arr.Remove(0)
	require.Equal(t, arr.Count(), 2)
	require.Equal(t, arr.String(), `["abc",456]`)

	ap := &AnyPerson{Name: "GGG", Age: 22}
	arr.Append(ap.AsAny())
	t.Log(arr.String())
	require.NotNil(t, AsAnyPerson(arr.ValueOf(2)))
	require.Nil(t, AsAnyPerson(arr.ValueOf(0)))
}

func TestAnyMap(t *testing.T) {
	mp := NewAnyMap()

	require.Equal(t, mp.String(), "{}")

	a1 := &Any{true}
	mp.SetValue(a1, "bbb")
	require.Equal(t, mp.Keys().Count(), 1)
	require.Equal(t, mp.Keys().String(), `["bbb"]`)
	require.Equal(t, mp.ValueOf("bbb"), a1)

	a2 := &Any{"abcd"}
	mp.SetValue(a2, "alpha")
	t.Log(mp.String())
	require.Equal(t, mp.HasKey("alpha"), true)
	require.Equal(t, mp.HasKey("beta"), false)

	require.Nil(t, mp.Remove("notkey"))
	require.Equal(t, mp.Remove("bbb"), a1)

	ap := &AnyPerson{Name: "GGG"}
	mp.SetValue(ap.AsAny(), "person")
	t.Log(mp.String())
	p := AsAnyPerson(mp.ValueOf("person"))
	require.Equal(t, p.Name, "GGG")
}
