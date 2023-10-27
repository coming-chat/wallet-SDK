package solanaswap

import (
	"context"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/coming-chat/wallet-SDK/core/solana"
	"github.com/stretchr/testify/require"
)

func RandomBytes32() []byte {
	radomId, _, _ := ed25519.GenerateKey(nil)
	return radomId
}

func TestOmniswap_estimate_relayer_fee(t *testing.T) {
	bsc_wormhole_chain_id := 4

	receiver, _ := hex.DecodeString("cAF084133CBdBE27490d3afB0Da220a40C32E307")
	receivingAssetId, _ := hex.DecodeString("51a3cc54eA30Da607974C5D07B8502599801AC08")
	sendingAssId := common.PublicKeyFromString("4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU")
	soData := SoData{
		TransactionId:      RandomBytes32(),
		Receiver:           receiver,
		SourceChainId:      1,
		SendingAssetId:     sendingAssId[:],
		DestinationChainId: 4,
		ReceivingAssetId:   receivingAssetId,
		Amount:             big.NewInt(0).SetInt64(1e6),
	}
	diamond, _ := hex.DecodeString("84B7cA95aC91f8903aCb08B27F5b41A4dE2Dc0fc")
	wormholeData := WormholeData{
		DstWormholeChainId:            4,
		DstMaxGasPriceInWeiForRelayer: 5e9,
		WormholeFee:                   1e3,
		DstSoDiamond:                  diamond,
	}

	soDataBytes, err := soData.Serialize()
	require.Nil(t, err)
	wormholeBytes, err := wormholeData.Serialize()
	require.Nil(t, err)
	fee, err := Omniswap_estimate_relayer_fee(devChain,
		uint16(bsc_wormhole_chain_id),
		wormholeBytes,
		soDataBytes,
		[]byte{})
	require.Nil(t, err)

	t.Log(fee)
}

type RelayerFee struct {
	SrcFee       uint16
	ConsumeValue uint16
	DstMaxGas    uint64
}

func Omniswap_estimate_relayer_fee(chain *solana.Chain,
	dstWormholeChainId uint16,
	wormholeData []byte,
	soData []byte,
	swapDataDst []byte) (*RelayerFee, error) {
	payer := account
	omniswap_program_id := config.Program.SoDiamond
	wormhole_program_id := config.Program.Wormhole

	fee_config_key, err := deriveSoFeeConfigKey(omniswap_program_id)
	if err != nil {
		return nil, err
	}
	foreign_contract_key, err := deriveForeignContractKey(omniswap_program_id, dstWormholeChainId)
	if err != nil {
		return nil, err
	}
	price_mannager_key, err := derivePriceManagerKey(omniswap_program_id, dstWormholeChainId)
	if err != nil {
		return nil, err
	}
	wormhole_bridge_data_key, err := deriveWormholeBridgeDataKey(wormhole_program_id)
	if err != nil {
		return nil, err
	}

	cli := chain.Client()
	recentHash, err := cli.GetLatestBlockhash(context.Background())
	if err != nil {
		return nil, err
	}
	message := types.NewMessage(types.NewMessageParam{
		FeePayer:        payer.Account().PublicKey,
		RecentBlockhash: recentHash.Blockhash,
		Instructions: []types.Instruction{
			EstimateRelayerFee(
				EstimateRelayerFeeArgs{
					ChainId:      dstWormholeChainId,
					SoData:       soData,
					WormholeData: wormholeData,
					SwapDataDst:  swapDataDst,
				},
				EstimateRelayerFeeAccounts{
					FeeConfig:       fee_config_key,
					ForeignContract: foreign_contract_key,
					PriceManager:    price_mannager_key,
					WormholeBridge:  wormhole_bridge_data_key,
				}),
		},
	})

	signedTxn, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: []types.Account{*payer.Account()},
	})

	txnResult, err := cli.SimulateTransaction(context.Background(), signedTxn)
	if err != nil {
		return nil, err
	}
	if txnResult.ReturnData == nil {
		return nil, errors.New("simulate transaction error: return data is nil")
	}

	returnData := txnResult.ReturnData.Data

	srcFee := binary.LittleEndian.Uint16(returnData[:8])
	consumeValue := binary.LittleEndian.Uint16(returnData[8:16])
	dstMaxGas := binary.LittleEndian.Uint64(returnData[16:])

	return &RelayerFee{
		SrcFee:       srcFee,
		ConsumeValue: consumeValue,
		DstMaxGas:    dstMaxGas,
	}, nil
}
