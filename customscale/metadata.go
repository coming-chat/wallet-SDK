package customscale

import (
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

func GetCallMethodFromMetadata(metadata *types.Metadata, call *types.Call) (variant types.Si1Variant, mod types.PalletMetadataV14) {
	switch metadata.Version {
	case 14:
		for _, mod = range metadata.AsMetadataV14.Pallets {
			if mod.Index != types.NewU8(call.CallIndex.SectionIndex) {
				continue
			}
			if typ, ok := metadata.AsMetadataV14.EfficientLookup[mod.Calls.Type.Int64()]; !ok || len(typ.Def.Variant.Variants) < 0 {
				break
			} else {
				for _, variant = range typ.Def.Variant.Variants {
					if variant.Index == types.NewU8(call.CallIndex.MethodIndex) {
						break
					}
				}
				break
			}

		}

	}
	return
}

func FindEventNamesForEventID(metadata *types.Metadata, eventID types.EventID) (types.Text, types.Text, []types.Si1Field, error) {
	switch metadata.Version {
	case 14:
		for _, mod := range metadata.AsMetadataV14.Pallets {
			if !mod.HasEvents {
				continue
			}
			if mod.Index != types.NewU8(eventID[0]) {
				continue
			}
			eventType := mod.Events.Type.Int64()
			typ, ok := metadata.AsMetadataV14.EfficientLookup[eventType]
			if !ok {
				continue
			}
			if len(typ.Def.Variant.Variants) <= 0 {
				continue
			}
			for _, vars := range typ.Def.Variant.Variants {
				if uint8(vars.Index) == eventID[1] {
					return mod.Name, vars.Name, vars.Fields, nil
				}
			}
		}
	case 12:
		for _, mod := range metadata.AsMetadataV12.Modules {
			if !mod.HasEvents {
				continue
			}
			if mod.Index != eventID[0] {
				continue
			}
			if int(eventID[1]) >= len(mod.Events) {
				continue
			}

			event := mod.Events[int(eventID[1])]
			return mod.Name, event.Name, nil, nil
		}
	default:
		return "", "", nil, fmt.Errorf("module index %v out of range", eventID[0])
	}
	return "", "", nil, fmt.Errorf("module index %v out of range", eventID[0])
}

func GetSi1TypeFromMetadata(metadata *types.Metadata, typeId types.Si1LookupTypeID) *types.Si1Type {
	for _, lookupType := range metadata.AsMetadataV14.Lookup.Types {
		if types.Eq(lookupType.ID, typeId) {
			return &lookupType.Type
		}
	}
	return nil
}
