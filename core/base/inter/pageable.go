package inter

import "github.com/coming-chat/wallet-SDK/core/base"

// `SdkPageable` implemented wallet-SDK/base's interface `Jsonable`
// If you new class `Xxx` extends it, you should implement `NewXxxWithJsonString` by your self.
type SdkPageable[T any] struct {
	TotalCount_    int    `json:"totalCount"`
	CurrentCount_  int    `json:"currentCount"`
	CurrentCursor_ string `json:"currentCursor"`
	HasNextPage_   bool   `json:"hasNextPage"`

	Items    []T `json:"items"`
	anyArray *base.AnyArray
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

func (p *SdkPageable[T]) JsonString() (*base.OptionalString, error) {
	return base.JsonString(p)
}

func (p *SdkPageable[T]) ItemArray() *base.AnyArray {
	if p.anyArray == nil {
		a := make([]any, len(p.Items))
		for idx, n := range p.Items {
			a[idx] = n
		}
		p.anyArray = &base.AnyArray{Values: a}
	}
	return p.anyArray
}

// It's will crash when index out of range
func (p *SdkPageable[T]) ItemAt(index int) T {
	return p.Items[index]
}
