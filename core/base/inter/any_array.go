package inter

import (
	"encoding/json"
)

type AnyArray[T any] []T

// `AnyArray` only support Marshal
func (a *AnyArray[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(*a)
}

func (a *AnyArray[T]) JsonString() string {
	data, err := json.Marshal(*a)
	if err != nil {
		return "null"
	}
	return string(data)
}

func (a *AnyArray[T]) Count() int {
	return len(*a)
}

func (a *AnyArray[T]) ValueAt(index int) T {
	return (*a)[index]
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

// return -1 if not found
func (a *AnyArray[T]) FirstIndexOf(matcher func(elem T) bool) int {
	for idx, elem := range *a {
		if matcher(elem) {
			return idx
		}
	}
	return -1
}

// return -1 if not found
func (a *AnyArray[T]) LastIndexOf(matcher func(elem T) bool) int {
	for i := len(*a) - 1; i >= 0; i-- {
		if matcher((*a)[i]) {
			return i
		}
	}
	return -1
}
