package inter

import (
	"testing"
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

	t.Log(p.TotalCount())
	t.Log(p.CurrentCount())
}

func TestFlowInterface(t *testing.T) {
	vv := PersonPage{
		&SdkPageable[*Person]{},
	}

	t.Log(vv.JsonString())
}
