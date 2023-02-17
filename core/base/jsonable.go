package base

import "encoding/json"

type Jsonable interface {
	JsonString() (*OptionalString, error)

	// You need to implement the following methods if your class name is Xxx
	// func NewXxxWithJsonString(str string) (*Xxx, error)

	// Maybe you can also implement this method
	// func NewXxxArrayWithJsonString(str string) (*AnyArray, error)

	// ====== template
	// func (j *Xxx) JsonString() (*base.OptionalString, error) {
	// 	return base.JsonString(s)
	// }
	// func NewXxxWithJsonString(str string) (*Xxx, error) {
	// 	var o Xxx
	// 	err := base.FromJsonString(str, &o)
	// 	return &o, err
	// }
	//	func NewXxxArrayWithJsonString(str string) (*base.AnyArray, error) {
	//		var o []*Xxx
	//		err := base.FromJsonString(str, &o)
	//		arr := make([]any, len(o))
	//		for i, v := range o {
	//			arr[i] = v
	//		}
	//		return &base.AnyArray{Values: arr}, err
	//	}
}

func JsonString(o interface{}) (*OptionalString, error) {
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return &OptionalString{Value: string(data)}, nil
}

func FromJsonString(jsonStr string, out interface{}) error {
	bytes := []byte(jsonStr)
	return json.Unmarshal(bytes, out)
}
