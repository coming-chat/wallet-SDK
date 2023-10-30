package solanaswap

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/compute_budget"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestTransfer_SendNativeToken(t *testing.T) {
	chain := devChain
	cli := chain.Client()

	txn, err := OmniswapSendNativeTokenTransaction(cli)
	require.Nil(t, err)

	simulate, err := cli.SimulateTransaction(context.Background(), *txn)
	require.Nil(t, err)
	t.Log(simulate)

	hash, err := cli.SendTransaction(context.Background(), *txn)
	require.Nil(t, err)
	t.Log(hash)
}

func OmniswapSendNativeTokenTransaction(cli *client.Client) (txn *types.Transaction, err error) {
	payer := account
	println("payer address =", payer.Address())

	omniswap_program_id := config.Program.SoDiamond
	wormhole_program_id := config.Program.Wormhole
	token_bridge_program_id := config.Program.TokenBridge
	lookup_table_key := common.PublicKeyFromString(config.Lookup_Table.Key)

	omnibtc_chainid_src := config.OmnibtcChainid
	omnibtc_chainid_dst := config.Wormhole.Dst_chain.OmnibtcChainid
	wormhole_dst_chain := config.Wormhole.Dst_chain.Chainid

	usdc_mint := common.PublicKeyFromString("4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU")
	usdc_token_on_solana := usdc_mint[:]
	// usdc account
	// usdc_account := common.PublicKeyFromString("68DjnBuZ6UtM6dGoTGhu2rqV5ZSowsPGgv2AWD1xuGB4")
	usdc_account := common.PublicKeyFromString("9LcCRussWRZqX7jXU6iB1KsXHHpbkRj2GrbH6wE6dyNs")
	// send 1 usdc
	amount := 1e6
	// recipient
	recipient_address, _ := hex.DecodeString("cAF084133CBdBE27490d3afB0Da220a40C32E307")
	usdc_token_on_bsc, _ := hex.DecodeString("51a3cc54eA30Da607974C5D07B8502599801AC08")
	dst_so_diamond_padding, _ := hex.DecodeString("00000000000000000000000084b7ca95ac91f8903acb08b27f5b41a4de2dc0fc")
	// collect relayer gas fee
	beneficiary_account := common.PublicKeyFromString(config.Beneficiary)

	send_native_accounts, err := getSendNativeTransferAccounts(
		token_bridge_program_id,
		wormhole_program_id,
		omniswap_program_id,
		wormhole_dst_chain,
		usdc_mint,
	)
	if err != nil {
		return
	}

	current_seq_bytes, err := cli.GetAccountInfo(context.Background(), send_native_accounts.token_bridge_sequence.ToBase58())
	if err != nil {
		return
	}
	current_seq := binary.LittleEndian.Uint64(current_seq_bytes.Data)
	println("current_seq =", current_seq)

	next_seq := current_seq + 1
	wormhole_message, err := deriveTokenTransferMessageKey(omniswap_program_id, next_seq)
	if err != nil {
		return
	}

	soDataTxnId := RandomBytes32()
	println("sodata txn id =", hex.EncodeToString(soDataTxnId))
	soData := SoData{
		TransactionId:      soDataTxnId,
		Receiver:           recipient_address,
		SourceChainId:      omnibtc_chainid_src,
		SendingAssetId:     usdc_token_on_solana,
		DestinationChainId: omnibtc_chainid_dst,
		ReceivingAssetId:   usdc_token_on_bsc,
		Amount:             big.NewInt(int64(amount)),
	}
	soDataBytes, err := soData.Serialize()
	if err != nil {
		return
	}
	wormholeData := WormholeData{
		DstWormholeChainId:            wormhole_dst_chain,
		DstMaxGasPriceInWeiForRelayer: 1e10,
		WormholeFee:                   0,
		DstSoDiamond:                  dst_so_diamond_padding,
	}
	wormholeDataBytes, err := wormholeData.Serialize()
	if err != nil {
		return
	}
	request_key, err := post_cross_requset(
		cli,
		wormhole_dst_chain,
		SoSwapPostCrossRequestArgs{
			so_data:       soDataBytes,
			wormhole_data: wormholeDataBytes,
		},
	)
	if err != nil {
		return
	}

	// ExceededMaxInstructions
	// devnet_limit=200_000, real=212433
	ix0 := compute_budget.SetComputeUnitLimit(compute_budget.SetComputeUnitLimitParam{
		Units: 300_000,
	})

	ix1 := so_swap_native_without_swap(
		SoSwapNativeWithoutSwapAccounts{
			payer:                         payer.Account().PublicKey,
			request:                       request_key,
			config:                        send_native_accounts.send_config,
			fee_config:                    send_native_accounts.fee_config,
			price_manager:                 send_native_accounts.price_manager,
			beneficiary_account:           beneficiary_account,
			foreign_contract:              send_native_accounts.foreign_contract,
			mint:                          usdc_mint,
			from_token_account:            usdc_account,
			tmp_token_account:             send_native_accounts.tmp_token_account,
			wormhole_program:              send_native_accounts.wormhole_program,
			token_bridge_program:          send_native_accounts.token_bridge_program,
			token_bridge_config:           send_native_accounts.token_bridge_config,
			token_bridge_custody:          send_native_accounts.token_bridge_custody,
			token_bridge_authority_signer: send_native_accounts.token_bridge_authority_signer,
			token_bridge_custody_signer:   send_native_accounts.token_bridge_custody_signer,
			wormhole_bridge:               send_native_accounts.wormhole_bridge,
			wormhole_message:              wormhole_message,
			token_bridge_emitter:          send_native_accounts.token_bridge_emitter,
			token_bridge_sequence:         send_native_accounts.token_bridge_sequence,
			wormhole_fee_collector:        send_native_accounts.wormhole_fee_collector,
		},
	)

	blockhash, err := cli.GetLatestBlockhash(context.Background())
	if err != nil {
		return
	}

	lookup_table_data, err := cli.GetAccountInfo(context.Background(), lookup_table_key.ToBase58())
	if err != nil {
		return
	}
	lookup_table_addresses := decode_address_look_up_table(lookup_table_data.Data)
	lookup_table := types.AddressLookupTableAccount{
		Key:       lookup_table_key,
		Addresses: lookup_table_addresses,
	}

	message := types.NewMessage(types.NewMessageParam{
		FeePayer:                   payer.Account().PublicKey,
		Instructions:               []types.Instruction{ix0, ix1},
		RecentBlockhash:            blockhash.Blockhash,
		AddressLookupTableAccounts: []types.AddressLookupTableAccount{lookup_table},
	})
	print(message.Version)

	txnn, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: []types.Account{*payer.Account()},
	})
	if err != nil {
		return
	}
	return &txnn, nil
}
