package solanaswap

import (
	"context"
	"encoding/binary"
	"errors"
	"time"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
)

func post_cross_requset(
	client *client.Client,
	dst_wormhole_chain_id uint16,
	args SoSwapPostCrossRequestArgs) (key common.PublicKey, err error) {
	payer := account

	omniswap_program_id := config.Program.SoDiamond
	wormhole_program_id := config.Program.Wormhole

	send_config_key, err := deriveSenderConfigKey(omniswap_program_id)
	if err != nil {
		return
	}
	sender_config, err := client.GetAccountInfo(context.Background(), send_config_key.ToBase58())
	if err != nil {
		return
	}
	request_seq := binary.LittleEndian.Uint64(sender_config.Data[len(sender_config.Data)-8:])
	println("request_seq =", request_seq)

	request_key, err := deriveCrossRequestKey(omniswap_program_id, request_seq)
	if err != nil {
		return
	}
	println("request_key =", request_key.ToBase58())
	fee_config_key, err := deriveSoFeeConfigKey(omniswap_program_id)
	if err != nil {
		return
	}
	foreign_contract_key, err := deriveForeignContractKey(omniswap_program_id, dst_wormhole_chain_id)
	if err != nil {
		return
	}
	price_manager_key, err := derivePriceManagerKey(omniswap_program_id, dst_wormhole_chain_id)
	if err != nil {
		return
	}
	wormhole_bridge_data_key, err := deriveWormholeBridgeDataKey(wormhole_program_id)
	if err != nil {
		return
	}

	ix := so_swap_post_cross_request(
		args,
		SoSwapPostCrossRequestAccounts{
			payer:            payer.Account().PublicKey,
			config:           send_config_key,
			request:          request_key,
			fee_config:       fee_config_key,
			foreign_contract: foreign_contract_key,
			price_manager:    price_manager_key,
			wormhole_bridge:  wormhole_bridge_data_key,
		},
	)

	blockHash, err := client.GetLatestBlockhash(context.Background())
	if err != nil {
		return
	}
	message := types.NewMessage(types.NewMessageParam{
		FeePayer:        payer.Account().PublicKey,
		RecentBlockhash: blockHash.Blockhash,
		Instructions:    []types.Instruction{ix},
	})

	txn, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: []types.Account{*payer.Account()},
	})
	if err != nil {
		return
	}
	// simulate, err := client.SimulateTransaction(context.Background(), txn)
	// if err != nil {
	// 	return
	// }
	// println(simulate.Logs)
	hash, err := client.SendTransaction(context.Background(), txn)
	if err != nil {
		return
	}
	for num := 0; num < 10; num++ {
		print("Transaction not confirmed yet. Waiting...")
		time.Sleep(5 * time.Second)
		resp, _ := client.GetTransaction(context.Background(), hash)
		if resp != nil && resp.BlockTime != nil {
			return request_key, nil
		}
	}
	return common.PublicKey{}, errors.New("post request timeout")
}
