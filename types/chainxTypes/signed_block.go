package chainxTypes

import "github.com/centrifuge/go-substrate-rpc-client/v4/types"

type SignedBlock struct {
	Block         Block         `json:"block"`
	Justification Justification `json:"justification"`
}

// Block encoded with header and extrinsics
type Block struct {
	Header     types.Header
	Extrinsics []Extrinsic
}
