package solanaswap

import (
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
)

type SoSwapNativeWithoutSwapAccounts struct {
	payer                         common.PublicKey
	request                       common.PublicKey
	config                        common.PublicKey
	fee_config                    common.PublicKey
	price_manager                 common.PublicKey
	beneficiary_account           common.PublicKey
	foreign_contract              common.PublicKey
	mint                          common.PublicKey
	from_token_account            common.PublicKey
	tmp_token_account             common.PublicKey
	wormhole_program              common.PublicKey
	token_bridge_program          common.PublicKey
	token_bridge_config           common.PublicKey
	token_bridge_custody          common.PublicKey
	token_bridge_custody_signer   common.PublicKey
	token_bridge_authority_signer common.PublicKey
	wormhole_bridge               common.PublicKey
	wormhole_message              common.PublicKey
	token_bridge_emitter          common.PublicKey
	token_bridge_sequence         common.PublicKey
	wormhole_fee_collector        common.PublicKey
}

func so_swap_native_without_swap(
	accounts SoSwapNativeWithoutSwapAccounts) types.Instruction {
	program_id := PROGRAM_ID

	keys := []types.AccountMeta{
		{PubKey: accounts.payer, IsSigner: true, IsWritable: true},
		{PubKey: accounts.request, IsSigner: false, IsWritable: true},
		{PubKey: accounts.config, IsSigner: false, IsWritable: false},
		{PubKey: accounts.fee_config, IsSigner: false, IsWritable: false},
		{PubKey: accounts.price_manager, IsSigner: false, IsWritable: false},
		{PubKey: accounts.beneficiary_account, IsSigner: false, IsWritable: true},
		{PubKey: accounts.foreign_contract, IsSigner: false, IsWritable: false},
		{PubKey: accounts.mint, IsSigner: false, IsWritable: true},
		{PubKey: accounts.from_token_account, IsSigner: false, IsWritable: true},
		{PubKey: accounts.tmp_token_account, IsSigner: false, IsWritable: true},
		{PubKey: accounts.wormhole_program, IsSigner: false, IsWritable: false},
		{PubKey: accounts.token_bridge_program, IsSigner: false, IsWritable: false},
		{PubKey: accounts.token_bridge_config, IsSigner: false, IsWritable: false},
		{PubKey: accounts.token_bridge_custody, IsSigner: false, IsWritable: true},
		{PubKey: accounts.token_bridge_custody_signer, IsSigner: false, IsWritable: false},
		{PubKey: accounts.token_bridge_authority_signer, IsSigner: false, IsWritable: false},
		{PubKey: accounts.wormhole_bridge, IsSigner: false, IsWritable: true},
		{PubKey: accounts.wormhole_message, IsSigner: false, IsWritable: true},
		{PubKey: accounts.token_bridge_emitter, IsSigner: false, IsWritable: true},
		{PubKey: accounts.token_bridge_sequence, IsSigner: false, IsWritable: true},
		{PubKey: accounts.wormhole_fee_collector, IsSigner: false, IsWritable: true},
		{PubKey: common.SystemProgramID, IsSigner: false, IsWritable: false},
		{PubKey: common.TokenProgramID, IsSigner: false, IsWritable: false},
		{PubKey: common.SPLAssociatedTokenAccountProgramID, IsSigner: false, IsWritable: false},
		{PubKey: common.SysVarClockPubkey, IsSigner: false, IsWritable: false},
		{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
	}

	identifier := []byte("d\x91\x88I_T]\xb1")
	data := identifier

	return types.Instruction{
		ProgramID: program_id,
		Accounts:  keys,
		Data:      data,
	}
}
