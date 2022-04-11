package types

import (
	"encoding/hex"
	"testing"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/itering/subscan/util/base58"
	"github.com/itering/subscan/util/ss58"
)

func TestGetXAssertsBalance(t *testing.T) {
	api, err := gsrpc.NewSubstrateAPI("wss://testnet3.chainx.org")
	if err != nil {
	}

	metadata, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return
	}

	ss58Format := base58.Decode("5QUEnWNMDFqsbUGpvvtgWGUgiiojnEpLf7581ELLAQyQ1xnT")
	publicKey, err := hex.DecodeString(ss58.Decode("5QUEnWNMDFqsbUGpvvtgWGUgiiojnEpLf7581ELLAQyQ1xnT", int(ss58Format[0])))
	if err != nil {
		return
	}

	assetId, err := types.EncodeToBytes(uint32(1))
	if err != nil {
		return
	}

	call, err := types.CreateStorageKey(metadata, "XAssets", "AssetBalance", publicKey, assetId)
	if err != nil {
		return
	}
	entryMetadata, err := metadata.FindStorageEntryMetadata("XAssets", "AssetBalance")
	if err != nil {
		return
	}
	i := entryMetadata.(types.StorageEntryMetadataV14).Type.AsMap.Value
	kIndex := metadata.AsMetadataV14.EfficientLookup[i.Int64()].Params[0].Type.Int64()
	vValue := metadata.AsMetadataV14.EfficientLookup[i.Int64()].Params[1].Type.Int64()
	data := NewMap(metadata.AsMetadataV14.EfficientLookup[kIndex], metadata.AsMetadataV14.EfficientLookup[vValue])
	//var res string
	//err = client.CallWithBlockHash(api.Client, &res, "state_getStorage", nil, call.Hex())
	//if err != nil {
	//	return
	//}
	ok, err := api.RPC.State.GetStorageLatest(call, &data)
	if err != nil {
		t.Log(err)
	}
	t.Log(ok)
	t.Log(data)
	t.Log(data.Data["Usable"])
}
