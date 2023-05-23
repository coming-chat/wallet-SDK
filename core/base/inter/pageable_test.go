package inter

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type Person struct {
	Name string
	Age  int
}

type PersonPage struct {
	*SdkPageable[*Person]
}

func (a *PersonPage) SecondObject() *Person {
	return a.Items[1]
}

func TestSdkPageable(t *testing.T) {
	p := PersonPage{
		&SdkPageable[*Person]{
			Items: []*Person{
				{
					Name: "aa",
					Age:  123,
				},
				{
					Name: "bb",
					Age:  999,
				},
			},
		},
	}
	t.Log(p.ItemAt(0))
	t.Log(p.SecondObject())

	t.Log(p.ItemArray().Values...)
	t.Log(p.TotalCount())
	t.Log(p.CurrentCount())
}

func TestFlowInterface(t *testing.T) {
	vv := PersonPage{
		&SdkPageable[*Person]{},
	}

	var jj base.Jsonable = vv
	t.Log(jj.JsonString())

	var pp base.Pageable = vv
	t.Log(pp.ItemArray())
}
