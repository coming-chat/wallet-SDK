package starknet

import (
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/xiang-xx/starknet.go/account"
	"github.com/xiang-xx/starknet.go/rpc"
	"github.com/xiang-xx/starknet.go/utils"
)

type DeployAccountTransaction struct {
	*rpc.DeployAccountTxn
	TransactionHash *felt.Felt
}

func (txn *DeployAccountTransaction) SignWithAccount(account base.Account) (signedTxn *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

func (txn *DeployAccountTransaction) SignedTransactionWithAccount(account base.Account) (signedTx base.SignedTransaction, err error) {
	starknetAccount := AsStarknetAccount(account)
	if starknetAccount == nil {
		return nil, base.ErrInvalidAccountType
	}
	txn.DeployAccountTxn.Signature, err = starknetAccount.SignHash(txn.TransactionHash)
	if err != nil {
		return nil, err
	}

	return &SignedTransaction{
		depolyTxn: txn.DeployAccountTxn,
	}, nil
}

func NewDeployAccountTransaction(pubkey string, maxFee *big.Int, acc *account.Account) (*DeployAccountTransaction, error) {
	pubFelt, err := utils.HexToFelt(pubkey)
	if err != nil {
		return nil, err
	}
	param := defaultDeployParam(*pubFelt)
	deployTxn := rpc.DeployAccountTxn{
		ContractAddressSalt: &param.Pubkey,
		ClassHash:           param.ClassHash,
		ConstructorCalldata: param.CallData,

		MaxFee:    utils.BigIntToFelt(maxFee),
		Version:   rpc.TransactionV1,
		Type:      rpc.TransactionType_DeployAccount,
		Nonce:     &felt.Zero,
		Signature: nil,
	}

	contractAddress, err := param.ComputeContractAddress()
	if err != nil {
		return nil, err
	}
	txnHash, err := acc.TransactionHashDeployAccount(deployTxn, contractAddress)
	if err != nil {
		return nil, err
	}
	return &DeployAccountTransaction{
		DeployAccountTxn: &deployTxn,
		TransactionHash:  txnHash,
	}, nil
}
