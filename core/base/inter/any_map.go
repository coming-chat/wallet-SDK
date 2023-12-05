package inter

import (
	"encoding/json"
)

// AnyMap
// ### Usage example for SDK
//
//	type StringMap struct { AnyMap[string, string] }
//	func NewStringMap() *StringMap { return &StringMap{map[string]string{}} }
//
// ### Usage Done
type AnyMap[K comparable, V any] map[K]V

func (a AnyMap[K, V]) MarshalJSON() ([]byte, error) {
	var temp map[K]V = a
	return json.Marshal(temp)
}

func (a *AnyMap[K, V]) UnmarshalJSON(data []byte) error {
	var out map[K]V
	err := json.Unmarshal(data, &out)
	*a = out
	return err
}

func (a AnyMap[K, V]) JsonString() string {
	data, err := json.Marshal(a)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func (a AnyMap[K, V]) Count() int {
	return len(a)
}

func (a AnyMap[K, V]) ValueOf(key K) V {
	return a[key]
}

func (a *AnyMap[K, V]) SetValue(value V, key K) {
	(*a)[key] = value
}

func (a *AnyMap[K, V]) Remove(key K) V {
	if v, ok := (*a)[key]; ok {
		delete((*a), key)
		return v
	}
	return (*a)[key]
}

// Deprecated: Use Contains(key) instead.
func (a AnyMap[K, V]) HasKey(key K) bool {
	_, ok := a[key]
	return ok
}

func (a AnyMap[K, V]) Contains(key K) bool {
	_, ok := a[key]
	return ok
}

// Keys
// 该方法的返回值无法打包到 sdk, 因此从对象方法中移出为公共方法
func KeysOf[K comparable, V any](m map[K]V) []K {
	keys := make([]K, len(m))
	idx := 0
	for k := range m {
		keys[idx] = k
		idx++
	}
	return keys
}
