package types

import (
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	k, v *types.Si1Type
)

type Map struct {
	Data map[interface{}]interface{}
}

func NewMap(kd, vd *types.Si1Type) Map {
	k = kd
	v = vd
	return Map{}
}

func (m *Map) Decode(decoder scale.Decoder) error {
	if k == nil || v == nil {
		return errors.New("please new map")
	}
	m.Data = map[interface{}]interface{}{}
	length, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}
	length = length >> 2
	if !k.Def.IsVariant {
		return errors.New("not supported type")
	}
	for i := 0; uint8(i) < length; i++ {
		kIndex, err := decoder.ReadOneByte()
		if err != nil {
			return err
		}
		for _, variant := range k.Def.Variant.Variants {
			if kIndex != byte(variant.Index) {
				continue
			}
			var v types.U128
			err = decoder.Decode(&v)
			if err != nil {
				return err
			}
			m.Data[string(variant.Name)] = v
			break
		}
	}
	return nil
}
