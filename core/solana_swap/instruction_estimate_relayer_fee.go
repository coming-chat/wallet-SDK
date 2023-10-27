package solanaswap

import (
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/bincode"
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
	data = append(data, bincode.MustSerializeData(args.ChainId)...)
	data = append(data, SerializeBytesWithLength(args.SoData)...)
	data = append(data, SerializeBytesWithLength(args.WormholeData)...)
	data = append(data, SerializeBytesWithLength(args.SwapDataDst)...)
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
