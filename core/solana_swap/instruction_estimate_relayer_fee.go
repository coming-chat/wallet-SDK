package solanaswap

import (
	"encoding/binary"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
)

type EstimateRelayerFeeArgs struct {
	ChainId      uint16
	SoData       []byte
	WormholeData []byte
	SwapDataDst  []byte
}

func (args *EstimateRelayerFeeArgs) Serialize() []byte {
	data := make([]byte, 0, 1024)
	buf := &data
	*buf = binary.LittleEndian.AppendUint16(*buf, args.ChainId)
	LittleEndianSerializeBytesWithLength(buf, args.SoData)
	LittleEndianSerializeBytesWithLength(buf, args.WormholeData)
	LittleEndianSerializeBytesWithLength(buf, args.SwapDataDst)
	return data
}

type EstimateRelayerFeeAccounts struct {
	FeeConfig       common.PublicKey
	ForeignContract common.PublicKey
	PriceManager    common.PublicKey
	WormholeBridge  common.PublicKey
}

func EstimateRelayerFee(args EstimateRelayerFeeArgs, accounts EstimateRelayerFeeAccounts) types.Instruction {
	data := args.Serialize()
	identifer := []byte("!\xe9(\x129\xfa:\x85")
	data = append(identifer, data...)

	return types.Instruction{
		Accounts: []types.AccountMeta{
			{PubKey: accounts.FeeConfig},
			{PubKey: accounts.ForeignContract},
			{PubKey: accounts.PriceManager},
			{PubKey: accounts.WormholeBridge},
		},
		ProgramID: PROGRAM_ID,
		Data:      data,
	}
}
