package solanaswap

import (
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
)

type SoSwapPostCrossRequestArgs struct {
	so_data       []byte
	swap_data_src []byte
	wormhole_data []byte
	swap_data_dst []byte
}

func (args *SoSwapPostCrossRequestArgs) Serialize() []byte {
	sl := NewSolanaSerializer(1024)
	sl.AppendBytesWithLenth(args.so_data)
	sl.AppendBytesWithLenth(args.swap_data_src)
	sl.AppendBytesWithLenth(args.wormhole_data)
	sl.AppendBytesWithLenth(args.swap_data_dst)
	return sl.Bytes()
}

type SoSwapPostCrossRequestAccounts struct {
	payer            common.PublicKey
	config           common.PublicKey
	request          common.PublicKey
	fee_config       common.PublicKey
	foreign_contract common.PublicKey
	price_manager    common.PublicKey
	wormhole_bridge  common.PublicKey
}

func so_swap_post_cross_request(
	args SoSwapPostCrossRequestArgs,
	accounts SoSwapPostCrossRequestAccounts,
) types.Instruction {
	program_id := PROGRAM_ID

	keys := []types.AccountMeta{
		{PubKey: accounts.payer, IsSigner: true, IsWritable: true},
		{PubKey: accounts.config, IsSigner: false, IsWritable: true},
		{PubKey: accounts.request, IsSigner: false, IsWritable: true},
		{PubKey: accounts.fee_config, IsSigner: false, IsWritable: false},
		{PubKey: accounts.foreign_contract, IsSigner: false, IsWritable: false},
		{PubKey: accounts.price_manager, IsSigner: false, IsWritable: false},
		{PubKey: accounts.wormhole_bridge, IsSigner: false, IsWritable: false},
		{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
	}

	identifer := []byte("<\xc1\x84\x89\xc9\xe6\x17~")
	data := append(identifer, args.Serialize()...)

	return types.Instruction{
		ProgramID: program_id,
		Accounts:  keys,
		Data:      data,
	}
}
