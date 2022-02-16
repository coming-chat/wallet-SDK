package customscale

import "github.com/centrifuge/go-substrate-rpc-client/v4/types"

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

func GetSi1TypeFromMetadata(metadata *types.Metadata, typeId types.Si1LookupTypeID) *types.Si1Type {
	for _, lookupType := range metadata.AsMetadataV14.Lookup.Types {
		if types.Eq(lookupType.ID, typeId) {
			return &lookupType.Type
		}
	}
	return nil
}
