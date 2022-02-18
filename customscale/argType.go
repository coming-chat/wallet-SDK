package customscale

import (
	"errors"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"io"
	"math"
	"reflect"
)

type ArgDecoder struct {
	*scale.Decoder
}

type Variants struct {
	MethodName string
	Value      []interface{}
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

func PrimitiveDecode(arg *ArgDecoder, typeDefPrimitive types.Si1TypeDefPrimitive) (interface{}, error) {
	switch typeDefPrimitive.Si0TypeDefPrimitive {
	case types.IsBool:
		var filed bool
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsChar:
		var filed byte
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsStr:
		var filed string
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsU8:
		var filed types.U8
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsU16:
		var filed types.U16
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsU32:
		var filed types.U32
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsU64:
		var filed types.U64
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsU128:
		var filed types.U128
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsU256:
		var filed types.U256
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsI8:
		var filed types.I8
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsI16:
		var filed types.I16
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsI32:
		var filed types.I32
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsI64:
		var filed types.I64
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsI128:
		var filed types.I128
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case types.IsI256:
		var filed types.I256
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	default:
		return nil, errors.New("not support filed")
	}
}

func DecodeByTypeID(metadata *types.Metadata, arg *ArgDecoder, typeId types.Si1LookupTypeID) (interface{}, error) {
	si1Type := GetSi1TypeFromMetadata(metadata, typeId)
	switch {
	case si1Type.Def.IsPrimitive:
		filed, err := PrimitiveDecode(arg, si1Type.Def.Primitive)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case si1Type.Def.IsArray:
		var filed []interface{}
		for i := 0; i < int(si1Type.Def.Array.Len); i++ {
			singleData, err := DecodeByTypeID(metadata, arg, si1Type.Def.Array.Type)
			if err != nil {
				return nil, err
			}
			filed = append(filed, singleData)
		}
		return filed, nil
	case si1Type.Def.IsBitSequence:
		var filed types.Si1TypeDefBitSequence
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case si1Type.Def.IsCompact:
		var filed types.UCompact
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case si1Type.Def.IsComposite:
		var values []interface{}
		filedList, err := ArgDecode(metadata, arg, si1Type.Def.Composite.Fields)
		if err != nil {
			return nil, err
		}
		for _, v := range filedList {
			values = append(values, v.Value)
		}
		return values, nil
	case si1Type.Def.IsTuple:
		var filed types.Si1TypeDefTuple
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case si1Type.Def.IsHistoricMetaCompat:
		var filed types.Type
		err := arg.Decode(&filed)
		if err != nil {
			return nil, err
		}
		return filed, nil
	case si1Type.Def.IsVariant:
		variants := Variants{}
		oneByte, err := arg.ReadOneByte()
		if err != nil {
			return nil, err
		}
		index := types.NewU8(oneByte)
		var internalVariant *types.Si1Variant
		for _, variant := range si1Type.Def.Variant.Variants {
			if variant.Index == index {
				internalVariant = &variant
				break
			}
		}
		if internalVariant == nil {
			return nil, errors.New("cannot find index")
		}
		variants.MethodName = string(internalVariant.Name)
		filedList, err := ArgDecode(metadata, arg, internalVariant.Fields)
		if err != nil {
			return nil, err
		}
		for _, v := range filedList {
			variants.Value = append(variants.Value, v.Value)
		}
		return variants, nil
	case si1Type.Def.IsSequence:
		codedLenUint, err := arg.DecodeUintCompact()
		if err != nil {
			return nil, err
		}
		if codedLenUint.Uint64() > math.MaxUint32 {
			return nil, errors.New("encoded array length is higher than allowed by the protocol (32-bit unsigned integer)")
		}
		if codedLenUint.Uint64() > uint64(int(^uint(0)>>1)) {
			return nil, errors.New("encoded array length is higher than allowed by the platform")
		}
		codedLen := int(codedLenUint.Uint64())
		var filed []interface{}
		for i := 0; i < codedLen; i++ {
			singleData, err := DecodeByTypeID(metadata, arg, si1Type.Def.Sequence.Type)
			if err != nil {
				return nil, err
			}
			filed = append(filed, singleData)
		}
		return filed, nil
	default:
		return nil, errors.New("typeId not support")
	}
}

func ArgDecode(metadata *types.Metadata, arg *ArgDecoder, si1Field []types.Si1Field) ([]*CallArg, error) {
	var callArg []*CallArg
	for i, field := range si1Field {
		data := &CallArg{}
		if field.HasName {
			data.FieldName = string(field.Name)
		} else {
			data.FieldName = fmt.Sprintf("unknown %d", i)
		}
		filed, err := DecodeByTypeID(metadata, arg, field.Type)
		if err != nil {
			return nil, err
		}
		data.Value = filed
		callArg = append(callArg, data)
	}
	return callArg, nil
}
