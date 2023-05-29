package inter

import (
	"encoding/json"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type AnyMap[T any] struct {
	Values map[string]T
}

// `AnyMap` only support Marshal
func (a AnyMap[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Values)
}
func (a *AnyMap[T]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &a.Values)
	return err
}

func (a *AnyMap[T]) JsonString() (*base.OptionalString, error) {
	return base.JsonString(a)
}

func (a *AnyMap[T]) ValueOf(key string) T {
	return a.Values[key]
}

func (a *AnyMap[T]) SetValue(value T, key string) {
	a.Values[key] = value
}

func (a *AnyMap[T]) Remove(key string) T {
	if v, ok := a.Values[key]; ok {
		delete(a.Values, key)
		return v
	}
	return a.Values[key]
}

func (a *AnyMap[T]) HasKey(key string) bool {
	_, ok := a.Values[key]
	return ok
}

func (a *AnyMap[T]) Keys() *base.StringArray {
	keys := make([]string, len(a.Values))
	i := 0
	for k := range a.Values {
		keys[i] = k
		i++
	}
	return &base.StringArray{Values: keys}
}

func (a *AnyMap[T]) String() string {
	data, err := json.Marshal(a.Values)
	if err != nil {
		return "{}"
	}
	return string(data)
}
