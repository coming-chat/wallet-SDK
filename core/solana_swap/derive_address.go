package solanaswap

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/blocto/solana-go-sdk/common"
)

// 优化: 方法接收 string 或 common.PublicKey

type SendNativeTransferAccounts struct {
	send_config                   common.PublicKey
	fee_config                    common.PublicKey
	price_manager                 common.PublicKey
	foreign_contract              common.PublicKey
	tmp_token_account             common.PublicKey
	token_bridge_config           common.PublicKey
	token_bridge_custody          common.PublicKey
	token_bridge_authority_signer common.PublicKey
	token_bridge_custody_signer   common.PublicKey
	wormhole_bridge               common.PublicKey
	token_bridge_emitter          common.PublicKey
	token_bridge_sequence         common.PublicKey
	wormhole_fee_collector        common.PublicKey
	wormhole_program              common.PublicKey
	token_bridge_program          common.PublicKey
}

func getSendNativeTransferAccounts(
	token_bridge_program_id string,
	wormhole_program_id string,
	omniswap_program_id string,
	recipient_chain uint16,
	native_mint_key common.PublicKey) (acc *SendNativeTransferAccounts, err error) {
	acc = &SendNativeTransferAccounts{}

	if acc.send_config, err = deriveSenderConfigKey(omniswap_program_id); err != nil {
		return
	}
	if acc.fee_config, err = deriveSoFeeConfigKey(omniswap_program_id); err != nil {
		return
	}
	if acc.price_manager, err = derivePriceManagerKey(omniswap_program_id, recipient_chain); err != nil {
		return
	}
	if acc.foreign_contract, err = deriveForeignContractKey(omniswap_program_id, recipient_chain); err != nil {
		return
	}
	if acc.tmp_token_account, err = deriveTmpTokenAccountKey(omniswap_program_id, native_mint_key); err != nil {
		return
	}
	if acc.token_bridge_config, err = deriveTokenBridgeConfigKey(token_bridge_program_id); err != nil {
		return
	}
	if acc.token_bridge_custody, err = deriveCustodyKey(token_bridge_program_id, native_mint_key); err != nil {
		return
	}
	if acc.token_bridge_authority_signer, err = deriveAuthoritySignerKey(token_bridge_program_id); err != nil {
		return
	}
	if acc.token_bridge_custody_signer, err = deriveCustodySignerKey(token_bridge_program_id); err != nil {
		return
	}
	if acc.wormhole_bridge, err = deriveWormholeBridgeDataKey(wormhole_program_id); err != nil {
		return
	}
	if acc.token_bridge_emitter, err = deriveWormholeEmitterKey(token_bridge_program_id); err != nil {
		return
	}
	if acc.token_bridge_sequence, err = deriveEmitterSequenceKey(wormhole_program_id, acc.token_bridge_emitter); err != nil {
		return
	}
	if acc.wormhole_fee_collector, err = deriveFeeCollectorKey(wormhole_program_id); err != nil {
		return
	}
	acc.wormhole_program = common.PublicKeyFromString(wormhole_program_id)
	acc.token_bridge_program = common.PublicKeyFromString(token_bridge_program_id)
	return
}

func deriveWormholeEmitterKey(programId string) (common.PublicKey, error) {
	return findProgramAddress(programId, "emitter")
}

func deriveEmitterSequenceKey(programId string, emitter common.PublicKey) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("Sequence"),
		emitter[:],
	})
}

func deriveWormholeBridgeDataKey(programId string) (common.PublicKey, error) {
	return findProgramAddress(programId, "Bridge")
}

func deriveFeeCollectorKey(programId string) (common.PublicKey, error) {
	return findProgramAddress(programId, "fee_collector")
}

func deriveSoFeeConfigKey(programId string) (common.PublicKey, error) {
	return findProgramAddress(programId, "so_fee")
}

func deriveSenderConfigKey(programId string) (common.PublicKey, error) {
	return findProgramAddress(programId, "sender")
}

func deriveRedeemerConfigKey(programId string) (common.PublicKey, error) {
	return findProgramAddress(programId, "redeemer")
}

func deriveForeignContractKey(programId string, chainId uint16) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("foreign_contract"),
		binary.LittleEndian.AppendUint16(make([]byte, 0, 2), chainId),
	})
}

func derivePriceManagerKey(programId string, chainId uint16) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("foreign_contract"),
		binary.LittleEndian.AppendUint16(make([]byte, 0, 2), chainId),
		[]byte("price_manager"),
	})
}

func deriveTokenTransferMessageKey(programId string, nextSeq uint64) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("bridged"),
		binary.LittleEndian.AppendUint64(make([]byte, 0, 8), nextSeq),
	})
}

func deriveCrossRequestKey(programId string, sequence uint64) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("request"),
		binary.LittleEndian.AppendUint64(make([]byte, 0, 8), sequence),
	})
}

func deriveForeignEndPointKey(programId string, chainId uint16, foreignContract common.PublicKey) (common.PublicKey, error) {
	if chainId == 1 {
		return common.PublicKey{}, errors.New("emitterChain == CHAIN_ID_SOLANA cannot exist as foreign token bridge emitter")
	}
	return findProgramAddressSeed(programId, [][]byte{
		binary.BigEndian.AppendUint16(make([]byte, 0, 2), chainId),
		foreignContract[:],
	})
}

func deriveTokenBridgeConfigKey(programId string) (common.PublicKey, error) {
	return findProgramAddress(programId, "config")
}

func deriveAuthoritySignerKey(programId string) (common.PublicKey, error) {
	return findProgramAddress(programId, "authority_signer")
}

func deriveCustodyKey(programId string, nativeMint common.PublicKey) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		nativeMint[:],
	})
}

func deriveCustodySignerKey(programId string) (common.PublicKey, error) {
	return findProgramAddress(programId, "custody_signer")
}

func deriveMintAuthorityKey(programId string) (common.PublicKey, error) {
	return findProgramAddress(programId, "mint_signer")
}

func deriveWrappedMintKey(programId string, tokenChain uint16, tokenAddress []byte) (common.PublicKey, error) {
	if tokenChain == 1 {
		return common.PublicKey{}, errors.New("tokenChain == CHAIN_ID_SOLANA does not have wrapped mint key")
	}
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("wrapped"),
		binary.BigEndian.AppendUint16(make([]byte, 0, 2), tokenChain),
		tokenAddress,
	})
}

func derivePostedVaaKey(programId string, hash []byte) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("PostedVAA"),
		hash,
	})
}

func deriveGuardianSetKey(programId string, index uint32) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("GuardianSet"),
		binary.BigEndian.AppendUint32(make([]byte, 0, 4), index),
	})
}

func deriveClaimKey(programId string, emitterAddress [32]byte, emitterChain uint16, sequence uint64) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		emitterAddress[:],
		binary.BigEndian.AppendUint16(make([]byte, 0, 2), emitterChain),
		binary.BigEndian.AppendUint64(make([]byte, 0, 8), sequence),
	})
}

func deriveWrappedMetaKey(programId string, mintKey common.PublicKey) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("meta"),
		mintKey[:],
	})
}

func deriveTmpTokenAccountKey(programId string, wrappedMint common.PublicKey) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("tmp"),
		wrappedMint[:],
	})
}

func deriveWhirlpoolOracleKey(programId string, whirlpool common.PublicKey) (common.PublicKey, error) {
	return findProgramAddressSeed(programId, [][]byte{
		[]byte("oracle"),
		whirlpool[:],
	})
}

// MARK - helper

func findProgramAddress(programId string, key string) (common.PublicKey, error) {
	id := common.PublicKeyFromString(programId)
	addr, _, err := common.FindProgramAddress([][]byte{[]byte(key)}, id)
	return addr, err
}

func findProgramAddressSeed(programId string, seed [][]byte) (common.PublicKey, error) {
	id := common.PublicKeyFromString(programId)
	addr, _, err := common.FindProgramAddress(seed, id)
	return addr, err
}

func decode_address_look_up_table(data []byte) []common.PublicKey {
	// https://github.com/solana-labs/solana-web3.js/blob/c7ef49cc49ee61422a4777d439a814160f6d7ce4/packages/library-legacy/src/programs/address-lookup-table/state.ts#L23
	LOOKUP_TABLE_META_SIZE := 56

	data_len := len(data)
	valid := data_len > LOOKUP_TABLE_META_SIZE && (data_len-LOOKUP_TABLE_META_SIZE)%32 == 0
	if !valid {
		err := fmt.Sprintf("data lenth error: %v", data_len)
		panic(err)
	}

	keys := make([]common.PublicKey, 0, LOOKUP_TABLE_META_SIZE)
	for i := LOOKUP_TABLE_META_SIZE; i < data_len; i += 32 {
		key := common.PublicKeyFromBytes(data[i : i+32])
		keys = append(keys, key)
	}

	return keys
}
