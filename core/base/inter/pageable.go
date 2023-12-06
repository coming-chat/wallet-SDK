package inter

import (
	"encoding/json"
)

// `SdkPageable` implemented wallet-SDK/base's interface `Jsonable`
// If you new class `Xxx` extends it, you should implement `NewXxxWithJsonString` by your self.
type SdkPageable[T any] struct {
	TotalCount_    int    `json:"totalCount"`
	CurrentCount_  int    `json:"currentCount"`
	CurrentCursor_ string `json:"currentCursor"`
	HasNextPage_   bool   `json:"hasNextPage"`

	Items []T `json:"items"`
}

func (p *SdkPageable[T]) TotalCount() int {
	return p.TotalCount_
}

func (p *SdkPageable[T]) CurrentCount() int {
	p.CurrentCount_ = len(p.Items)
	return p.CurrentCount_
}

func (p *SdkPageable[T]) CurrentCursor() string {
	return p.CurrentCursor_
}

func (p *SdkPageable[T]) HasNextPage() bool {
	return p.HasNextPage_
}

func (p *SdkPageable[T]) JsonString() string {
	data, err := json.Marshal(p)
	if err != nil {
		return "null"
	}
	return string(data)
}

// It's will crash when index out of range
func (p *SdkPageable[T]) ItemAt(index int) T {
	return p.Items[index]
}
