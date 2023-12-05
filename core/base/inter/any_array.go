package inter

import (
	"encoding/json"
)

// AnyArray
// ### Usage example for SDK
//
//	type StringArray struct { AnyArray[string] }
//	func NewStringArray() *StringArray { return &StringArray{[]string{}} }
//
// ### Usage Done
type AnyArray[T any] []T

func (a AnyArray[T]) MarshalJSON() ([]byte, error) {
	var temp []T = a
	return json.Marshal(temp)
}

func (a *AnyArray[T]) UnmarshalJSON(data []byte) error {
	var out []T
	err := json.Unmarshal(data, &out)
	*a = out
	return err
}

func (a AnyArray[T]) JsonString() string {
	data, err := json.Marshal(a)
	if err != nil {
		return "null"
	}
	return string(data)
}

func (a AnyArray[T]) Count() int {
	return len(a)
}

func (a AnyArray[T]) ValueAt(index int) T {
	return a[index]
}

func (a *AnyArray[T]) Append(value T) {
	*a = append(*a, value)
}

func (a *AnyArray[T]) Remove(index int) T {
	r := (*a)[index]
	*a = append((*a)[:index], (*a)[index+1:]...)
	return r
}

func (a *AnyArray[T]) SetValue(value T, index int) {
	(*a)[index] = value
}

// FirstIndexOf
// 该方法的参数无法打包到 sdk, 因此从对象方法中移出为公共方法
// return -1 if not found
func FirstIndexOf[T any](arr []T, matcher func(elem T) bool) int {
	for idx, elem := range arr {
		if matcher(elem) {
			return idx
		}
	}
	return -1
}

// LastIndexOf
// 该方法的参数无法打包到 sdk, 因此从对象方法中移出为公共方法
// return -1 if not found
func LastIndexOf[T any](arr []T, matcher func(elem T) bool) int {
	for i := len(arr) - 1; i >= 0; i-- {
		if matcher(arr[i]) {
			return i
		}
	}
	return -1
}
