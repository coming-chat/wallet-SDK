package customscale

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type CallData struct {
	Mod     string
	Variant string
	Arg     []*CallArg
}

type CallArg struct {
	FieldName string
	Value     interface{}
}

func DecodeCall(metadata *types.Metadata, call *types.Call) (*CallData, error) {
	variant, mod := GetCallMethodFromMetadata(metadata, call)
	var (
		args  []types.Si1LookupTypeID
		names []string
	)

	for i, field := range variant.Fields {
		if field.HasName {
			names = append(names, string(field.Name))
		} else {
			names = append(names, fmt.Sprintf("unknown %d", i))
		}
		args = append(args, field.Type)
	}
	callArgList, err := ArgDecode(metadata, call.Args, args)
	if err != nil {
		return nil, err
	}
	var callArg []*CallArg
	for i, v := range callArgList {
		callArg = append(callArg, &CallArg{
			FieldName: names[i],
			Value:     v,
		})
	}
	return &CallData{
		Mod:     string(mod.Name),
		Variant: string(variant.Name),
		Arg:     callArg,
	}, nil
}
