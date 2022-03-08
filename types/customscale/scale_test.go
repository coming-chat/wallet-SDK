package customscale

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"log"
	"testing"
)

func TestDecodeEvent(t *testing.T) {
	//eventString := "0x1c00000000000000585f8f090000000002000000010000000a066c707b1690a6b0e01b5dea252fe1887930a5afc0ec203f96705331749c37ae4a000064a7b3b6e00d000000000000000000000100000035016c707b1690a6b0e01b5dea252fe1887930a5afc0ec203f96705331749c37ae4a000001000000000000e1f5050000000000000000020000000c0303000000502334ba4a30b12b38ba5f8e1fa719ebb6420fdb360abf915d0d4b3656ae214140420f000000000000000000000000000000020000003504502334ba4a30b12b38ba5f8e1fa719ebb6420fdb360abf915d0d4b3656ae214140420f000000000000000000000000000303000000000002000000000000e1f50500000000000000"
	api, err := gsrpc.NewSubstrateAPI("wss://mainnet.chainx.org/ws")

	if err != nil {
		log.Printf("get metadata error: %v", err)
	}

	blockHash, err := api.RPC.Chain.GetBlockHash(2809362)
	if err != nil {
		return
	}

	var res string
	err = client.CallWithBlockHash(api.Client, &res, "state_getMetadata", &blockHash)
	if err != nil {
		log.Println(err)
	}

	var metadata types.Metadata

	err = types.DecodeFromHexString(res, &metadata)
	if err != nil {
		return
	}

	call, err := types.CreateStorageKey(&metadata, "System", "Events")
	if err != nil {
		return
	}

	rawData, err := api.RPC.State.GetStorageRaw(call, blockHash)
	if err != nil {
		return
	}
	eventRecord := EventRaw(*rawData)
	eventData, err := eventRecord.DecodeRaw(&metadata)
	if err != nil {
		return
	}
	t.Log(eventData)
}
