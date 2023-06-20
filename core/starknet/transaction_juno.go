package starknet

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/crypto"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/juno/utils"
)

func TransactionHash(transaction core.Transaction, n utils.Network) (*felt.Felt, error) {
	switch t := transaction.(type) {
	case *core.DeclareTransaction:
		return declareTransactionHash(t, n)
	case *core.InvokeTransaction:
		return invokeTransactionHash(t, n)
	case *core.DeployTransaction:
		// deploy transactions are deprecated after re-genesis therefore we don't verify
		// transaction hash
		return t.TransactionHash, nil
	case *core.L1HandlerTransaction:
		return l1HandlerTransactionHash(t, n)
	case *core.DeployAccountTransaction:
		return deployAccountTransactionHash(t, n)
	default:
		return nil, errors.New("unknown transaction")
	}
}

var (
	invokeFelt        = new(felt.Felt).SetBytes([]byte("invoke"))
	declareFelt       = new(felt.Felt).SetBytes([]byte("declare"))
	l1HandlerFelt     = new(felt.Felt).SetBytes([]byte("l1_handler"))
	deployAccountFelt = new(felt.Felt).SetBytes([]byte("deploy_account"))
)

func errInvalidTransactionVersion(t core.Transaction, version *felt.Felt) error {
	return fmt.Errorf("invalid Transaction (type: %v) version: %v", reflect.TypeOf(t), version.Text(felt.Base10))
}

func invokeTransactionHash(i *core.InvokeTransaction, n utils.Network) (*felt.Felt, error) {
	switch {
	case i.Version.IsZero():
		// Due to inconsistencies in version 0 hash calculation we don't verify the hash
		return i.TransactionHash, nil
	case i.Version.IsOne():
		return crypto.PedersenArray(
			invokeFelt,
			i.Version,
			i.SenderAddress,
			new(felt.Felt),
			crypto.PedersenArray(i.CallData...),
			i.MaxFee,
			n.ChainID(),
			i.Nonce,
		), nil
	default:
		return nil, errInvalidTransactionVersion(i, i.Version)
	}
}

func declareTransactionHash(d *core.DeclareTransaction, n utils.Network) (*felt.Felt, error) {
	switch {
	case d.Version.IsZero():
		// Due to inconsistencies in version 0 hash calculation we don't verify the hash
		return d.TransactionHash, nil
	case d.Version.IsOne():
		return crypto.PedersenArray(
			declareFelt,
			d.Version,
			d.SenderAddress,
			new(felt.Felt),
			crypto.PedersenArray(d.ClassHash),
			d.MaxFee,
			n.ChainID(),
			d.Nonce,
		), nil
	case d.Version.Equal(new(felt.Felt).SetUint64(2)):
		return crypto.PedersenArray(
			declareFelt,
			d.Version,
			d.SenderAddress,
			&felt.Zero,
			crypto.PedersenArray(d.ClassHash),
			d.MaxFee,
			n.ChainID(),
			d.Nonce,
			d.CompiledClassHash,
		), nil

	default:
		return nil, errInvalidTransactionVersion(d, d.Version)
	}
}

func l1HandlerTransactionHash(l *core.L1HandlerTransaction, n utils.Network) (*felt.Felt, error) {
	switch {
	case l.Version.IsZero():
		// There are some l1 handler transaction which do not return a nonce and for some random
		// transaction the following hash fails.
		if l.Nonce == nil {
			return l.TransactionHash, nil
		}
		return crypto.PedersenArray(
			l1HandlerFelt,
			l.Version,
			l.ContractAddress,
			l.EntryPointSelector,
			crypto.PedersenArray(l.CallData...),
			&felt.Zero,
			n.ChainID(),
			l.Nonce,
		), nil
	default:
		return nil, errInvalidTransactionVersion(l, l.Version)
	}
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
