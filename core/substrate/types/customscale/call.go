package customscale

import (
	"bytes"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type CallData struct {
	Method string
	Arg    []*CallArg
}

type CallArg struct {
	FieldName string
	Value     any
}

func DecodeCall(metadata *types.Metadata, call *types.Call) (*CallData, error) {
	variant, mod := GetCallMethodFromMetadata(metadata, call)

	arg := NewArgDecoder(bytes.NewReader(call.Args))
	callArg, err := ArgDecode(metadata, arg, variant.Fields)
	if err != nil {
		return &CallData{
			Method: fmt.Sprintf("%s.%s", string(mod.Name), string(variant.Name)),
			Arg:    nil,
		}, err
	}
	return &CallData{
		Method: fmt.Sprintf("%s.%s", string(mod.Name), string(variant.Name)),
		Arg:    callArg,
	}, nil
}
