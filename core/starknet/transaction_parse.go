package starknet

import (
	"math/big"

	"github.com/NethermindEth/juno/core"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
)

func parseDeployAccountTransaction(txn *core.DeployAccountTransaction) *types.DeployAccountRequest {
	req := types.DeployAccountRequest{
		Type: gateway.DEPLOY_ACCOUNT,
	}
	if txn.ContractAddressSalt != nil {
		req.ContractAddressSalt = txn.ContractAddressSalt.String()
	}
	if txn.ClassHash != nil {
		req.ClassHash = txn.ClassHash.String()
	}
	if txn.ConstructorCallData != nil {
		callData := make([]string, len(txn.ConstructorCallData))
		for i, felt := range txn.ConstructorCallData {
			callData[i] = felt.String()
		}
		req.ConstructorCalldata = callData
	}
	if txn.Version != nil {
		req.Version = txn.Version.BigInt(&big.Int{}).Uint64()
	}
	if txn.MaxFee != nil {
		req.MaxFee = txn.MaxFee.BigInt(&big.Int{})
	}
	if txn.TransactionSignature != nil {
		signatures := make([]*big.Int, len(txn.TransactionSignature))
		for i, felt := range txn.TransactionSignature {
			signatures[i] = felt.BigInt(&big.Int{})
		}
		req.Signature = types.Signature(signatures)
	}
	// if txn.Nonce != nil {
	// 	req.Nonce = nil // no need nonce
	// }
	return &req
}
