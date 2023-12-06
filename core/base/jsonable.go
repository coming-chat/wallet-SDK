package base

import "encoding/json"

type Jsonable interface {
	JsonString() (*OptionalString, error)

	// You need to implement the following methods if your class name is Xxx
	// func NewXxxWithJsonString(str string) (*Xxx, error)

	// ====== template
	// func (j *Xxx) JsonString() (*base.OptionalString, error) {
	// 	return base.JsonString(j)
	// }
	// func NewXxxWithJsonString(str string) (*Xxx, error) {
	// 	var o Xxx
	// 	err := base.FromJsonString(str, &o)
	// 	return &o, err
	// }
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
