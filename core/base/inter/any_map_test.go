package inter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type stringMap struct {
	AnyMap[string, string]
}

func NewStringMap() *stringMap {
	return &stringMap{map[string]string{}}
}

func TestStringMap(t *testing.T) {
	m := NewStringMap()
	t.Log(m.JsonString())

	m.SetValue("saa", "name")
	m.SetValue("12", "age")
	require.Equal(t, m.Count(), 2)
	require.Equal(t, m.ValueOf("name"), "saa")
	require.Equal(t, m.ValueOf("age"), "12")
	require.Equal(t, m.ValueOf("invalidkey"), "")
	require.True(t, m.Contains("age"))

	res := m.Remove("age")
	require.Equal(t, res, "12")
	require.False(t, m.Contains("age"))

	m.SetValue("180", "height")
	t.Log(m.JsonString())
	t.Log(KeysOf(m.AnyMap))
	require.Equal(t, KeysOf(m.AnyMap), []string{"name", "height"})
}

func TestStringMap_Unmarshal(t *testing.T) {
	str := `{
		"height": "1000",
		"width": "200"
	}`

	var m stringMap
	err := json.Unmarshal([]byte(str), &m)
	require.Nil(t, err)

	t.Log(m.JsonString())
}
