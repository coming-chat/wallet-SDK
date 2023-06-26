package starknet

import (
	"errors"
	"math/big"

	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/dontpanicdao/caigo"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
)

var (
	DeployAccountHash = big.NewInt(0).SetBytes([]byte("deploy_account"))

	ErrInvalidDeployTransactionVersion = errors.New("invalid deploy transaction version, only support version 1")
)

type DeployAccountTransaction struct {
	TransactionHash *big.Int
	// A random number used to distinguish between different instances of the contract.
	ContractAddressSalt *big.Int
	// The address of the contract.
	ContractAddress *big.Int
	// The hash of the class which defines the contract’s functionality.
	ClassHash *big.Int
	// The arguments passed to the constructor during deployment.
	ConstructorCallData []*big.Int
	// The transaction’s version. Possible values are 1 or 0.
	//
	// When the fields that comprise a transaction change,
	// either with the addition of a new field or the removal of an existing field,
	// then the transaction version increases.
	// Transaction version 0 is deprecated and will be removed in a future version of Starknet.
	Version *big.Int

	// The maximum fee that the sender is willing to pay for the transaction.
	MaxFee *big.Int
	// Additional information given by the sender, used to validate the transaction.
	TransactionSignature []*big.Int
	// The transaction nonce.
	Nonce *big.Int

	Network Network
}

func (txn *DeployAccountTransaction) SignWithAccount(account base.Account) (signedTxn *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

func (txn *DeployAccountTransaction) SignedTransactionWithAccount(account base.Account) (signedTx base.SignedTransaction, err error) {
	starknetAccount := AsStarknetAccount(account)
	if starknetAccount == nil {
		return nil, base.ErrInvalidAccountType
	}
	txn.TransactionHash, err = txn.TransactionHashAt(txn.Network)
	if err != nil {
		return nil, err
	}
	s1, s2, err := starknetAccount.SignHash(txn.TransactionHash)
	if err != nil {
		return nil, err
	}
	txn.TransactionSignature = []*big.Int{s1, s2}

	return &SignedTransaction{
		Account:   starknetAccount,
		depolyTxn: txn,
	}, nil
}

func newDeployAccountTransactionForArgentX(pubkey string, network Network) (*DeployAccountTransaction, error) {
	pubData, err := hexTypes.HexDecodeString(pubkey)
	if err != nil {
		return nil, base.ErrInvalidPublicKey
	}
	pubkeyInt := big.NewInt(0).SetBytes(pubData)

	txn := DeployAccountTransaction{
		ClassHash:           types.HexToBN("0x25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918"),
		ContractAddressSalt: pubkeyInt,
		ConstructorCallData: []*big.Int{
			types.HexToBN("0x33434ad846cdd5f23eb73ff09fe6fddd568284a0fb7d1be20ee482f044dabe2"),
			types.HexToBN("0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463"),
			types.HexToBN("0x2"),
			pubkeyInt,
			types.HexToBN("0x0"),
		},
		Version: big.NewInt(1),
		MaxFee:  big.NewInt(1e13),
		Nonce:   big.NewInt(0),

		Network: network,
	}

	callerAddress := big.NewInt(0)
	txn.ContractAddress, err = txn.ComputeContractAddress(callerAddress)
	if err != nil {
		return nil, err
	}
	return &txn, nil
}

// ContractAddress computes the address of a Starknet contract.
func (txn *DeployAccountTransaction) ComputeContractAddress(callerAddress *big.Int) (*big.Int, error) {
	prefix := big.NewInt(0).SetBytes([]byte(caigo.CONTRACT_ADDRESS_PREFIX))
	callDataHash, err := caigo.Curve.ComputeHashOnElements(txn.ConstructorCallData)
	if err != nil {
		return nil, err
	}

	// https://docs.starknet.io/documentation/architecture_and_concepts/Contracts/contract-address
	return caigo.Curve.ComputeHashOnElements([]*big.Int{
		prefix,
		callerAddress,
		txn.ContractAddressSalt,
		txn.ClassHash,
		callDataHash,
	})
}

func (txn *DeployAccountTransaction) TransactionHashAt(network Network) (*big.Int, error) {
	callData := []*big.Int{txn.ClassHash, txn.ContractAddressSalt}
	callData = append(callData, txn.ConstructorCallData...)
	callDataHash, err := caigo.Curve.ComputeHashOnElements(callData)
	if err != nil {
		return nil, err
	}
	chainHash, err := NetworkChainID(network)
	if err != nil {
		return nil, err
	}

	if big.NewInt(1).Cmp(txn.Version) == 0 {
		return caigo.Curve.ComputeHashOnElements([]*big.Int{
			DeployAccountHash,
			txn.Version,
			txn.ContractAddress,
			big.NewInt(0),
			callDataHash,
			txn.MaxFee,
			chainHash,
			txn.Nonce,
		})
	}
	return nil, ErrInvalidDeployTransactionVersion
}

func (txn *DeployAccountTransaction) CaigoDeployAccountRequest() *types.DeployAccountRequest {
	req := types.DeployAccountRequest{
		Type: gateway.DEPLOY_ACCOUNT,
	}
	if txn.ContractAddressSalt != nil {
		req.ContractAddressSalt = types.BigToHex(txn.ContractAddressSalt)
	}
	if txn.ClassHash != nil {
		req.ClassHash = types.BigToHex(txn.ClassHash)
	}
	if txn.ConstructorCallData != nil {
		callData := make([]string, len(txn.ConstructorCallData))
		for i, data := range txn.ConstructorCallData {
			callData[i] = types.BigToHex(data)
		}
		req.ConstructorCalldata = callData
	}
	if txn.Version != nil {
		req.Version = txn.Version.Uint64()
	}
	if txn.MaxFee != nil {
		req.MaxFee = txn.MaxFee
	}
	if txn.TransactionSignature != nil {
		req.Signature = txn.TransactionSignature
	}
	// if txn.Nonce != nil {
	// 	req.Nonce = nil // can not set nonce
	// }
	return &req
}
