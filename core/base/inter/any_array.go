package inter

import (
	"encoding/json"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type AnyArray[T any] []T

// `AnyArray` only support Marshal
func (a *AnyArray[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(*a)
}

func (j *AnyArray[T]) JsonString() (*base.OptionalString, error) {
	return base.JsonString(j)
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
