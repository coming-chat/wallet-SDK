package customscale

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"io"
	"reflect"
)

type ArgDecoder struct {
	*scale.Decoder
}

func NewArgDecoder(reader io.Reader) *ArgDecoder {
	return &ArgDecoder{scale.NewDecoder(reader)}
}

func (pd ArgDecoder) Decode(target ...interface{}) error {
	for _, tg := range target {
		t0 := reflect.TypeOf(tg)
		if t0.Kind() != reflect.Ptr {
			return errors.New("target must be a pointer, but was " + fmt.Sprint(t0))
		}
		val := reflect.ValueOf(tg)
		if val.IsNil() {
			return errors.New("target is a nil pointer")
		}
		err := pd.DecodeIntoReflectValue(val.Elem())
		if err != nil {
			return err
		}
	}
	return nil
}

func ArgDecode(metadata *types.Metadata, argBytes []byte, args []types.Si1LookupTypeID) (callArg []interface{}, err error) {
	var (
		arg = NewArgDecoder(bytes.NewReader(argBytes))
	)

	for _, typeId := range args {
		si1Type := GetSi1TypeFromMetadata(metadata, typeId)
		switch {
		case si1Type.Def.IsPrimitive:
			switch si1Type.Def.Primitive.Si0TypeDefPrimitive {
			case types.IsBool:
				var filed bool
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsChar:
				var filed byte
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsStr:
				var filed string
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsU8:
				var filed types.U8
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsU16:
				var filed types.U16
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsU32:
				var filed types.U32
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsU64:
				var filed types.U64
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsU128:
				var filed types.U128
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsU256:
				var filed types.U256
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsI8:
				var filed types.I8
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsI16:
				var filed types.I16
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsI32:
				var filed types.I32
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsI64:
				var filed types.I64
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsI128:
				var filed types.I128
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			case types.IsI256:
				var filed types.I256
				err = arg.Decode(&filed)
				if err != nil {
					return nil, err
				}
				callArg = append(callArg, filed)
			}
		case si1Type.Def.IsArray:
		case si1Type.Def.IsBitSequence:
		case si1Type.Def.IsCompact:
		case si1Type.Def.IsComposite:
		case si1Type.Def.IsTuple:
		case si1Type.Def.IsHistoricMetaCompat:
		case si1Type.Def.IsVariant:
		case si1Type.Def.IsSequence:

		}

	}
	return
}
