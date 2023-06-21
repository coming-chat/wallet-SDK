package starknet

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/utils"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/dontpanicdao/caigo/gateway"
	"github.com/dontpanicdao/caigo/types"
)

type DeployAccountTransaction struct {
	txn     *core.DeployAccountTransaction
	network utils.Network
}

func (t *DeployAccountTransaction) SignWithAccount(account base.Account) (signedTxn *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

func (t *DeployAccountTransaction) SignedTransactionWithAccount(account base.Account) (signedTx base.SignedTransaction, err error) {
	starknetAccount := AsStarknetAccount(account)
	if starknetAccount == nil {
		return nil, base.ErrInvalidAccountType
	}
	txnHash, err := deployAccountTransactionHash(t.txn, t.network)
	if err != nil {
		return nil, err
	}
	s1, s2, err := starknetAccount.SignHash(txnHash.BigInt(&big.Int{}))
	if err != nil {
		return nil, err
	}

	txnReq := *parseDeployAccountTransaction(t.txn)
	txnReq.Signature = types.Signature{s1, s2}
	return &SignedTransaction{
		Account:   starknetAccount,
		depolyTxn: &txnReq,
	}, nil
}

var (
	deployAccountFelt = new(felt.Felt).SetBytes([]byte("deploy_account"))
)

func errInvalidTransactionVersion(t core.Transaction, version *felt.Felt) error {
	return fmt.Errorf("invalid Transaction (type: %v) version: %v", reflect.TypeOf(t), version.Text(felt.Base10))
}

func deployAccountTransactionHash(d *core.DeployAccountTransaction, n utils.Network) (*felt.Felt, error) {
	callData := []*felt.Felt{d.ClassHash, d.ContractAddressSalt}
	callData = append(callData, d.ConstructorCallData...)
	// There is no version 0 for deploy account
	if d.Version.IsOne() {
		return crypto.PedersenArray(
			deployAccountFelt,
			d.Version,
			d.ContractAddress,
			&felt.Zero,
			crypto.PedersenArray(callData...),
			d.MaxFee,
			n.ChainID(),
			d.Nonce,
		), nil
	}
	return nil, errInvalidTransactionVersion(d, d.Version)
}

func deployAccountTxnForArgentX(pubKey string) (*core.DeployAccountTransaction, error) {
	pubData, err := new(felt.Felt).SetString(pubKey)
	if err != nil {
		return nil, base.ErrInvalidPublicKey
	}
	classHash, _ := new(felt.Felt).SetString("0x25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918")
	data1, _ := new(felt.Felt).SetString("0x33434ad846cdd5f23eb73ff09fe6fddd568284a0fb7d1be20ee482f044dabe2")
	data2, _ := new(felt.Felt).SetString("0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463") // getSelectorFromName("initialize")
	data3, _ := new(felt.Felt).SetString("0x2")
	data4 := pubData
	data5, _ := new(felt.Felt).SetString("0x0")
	txn := &core.DeployAccountTransaction{
		DeployTransaction: core.DeployTransaction{
			ClassHash:           classHash,
			ContractAddressSalt: pubData,
			ConstructorCallData: []*felt.Felt{
				data1, data2, data3, data4, data5,
			},
			Version: new(felt.Felt).SetUint64(1),
		},
		MaxFee: new(felt.Felt).SetBytes(big.NewInt(1e13).Bytes()),
		Nonce:  new(felt.Felt).SetUint64(0),
	}

	callerAddress, _ := new(felt.Felt).SetString("0x0000000000000000000000000000000000000000")
	address := core.ContractAddress(callerAddress, txn.ClassHash, txn.ContractAddressSalt, txn.ConstructorCallData)
	txn.ContractAddress = address

	return txn, nil
}

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
	// 	req.Nonce = nil // can not set nonce
	// }
	return &req
}
