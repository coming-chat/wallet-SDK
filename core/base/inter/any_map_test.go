package inter

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type StrMap struct {
	AnyMap[string]
}

func NewStrMap() *StrMap {
	return &StrMap{AnyMap: AnyMap[string]{map[string]string{}}}
}

func TestStrMap(t *testing.T) {
	m := NewStrMap()
	m.SetValue("saa", "name")
	m.SetValue("12", "age")
	t.Log(m.String())

	require.Equal(t, m.ValueOf("name"), "saa")
	require.Equal(t, m.ValueOf("height"), "")
	require.True(t, m.HasKey("age"))

	res := m.Remove("age")
	require.Equal(t, res, "12")
	require.False(t, m.HasKey("age"))

	t.Log(m.String())
}

func TestStrMap_json(t *testing.T) {
	str := `{
		"height": "1000",
		"width": "200"
	}`

	var m StrMap
	err := json.Unmarshal([]byte(str), &m)
	require.Nil(t, err)

	t.Log(m.String())
}
