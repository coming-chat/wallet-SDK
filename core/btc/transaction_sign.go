package btc

import (
	"bytes"
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/util/hexutil"
)

type SignedTransaction struct {
	msgTx *wire.MsgTx
}

func (t *Transaction) SignWithAccount(account base.Account) (signedTxn *base.OptionalString, err error) {
	txn, err := t.SignedTransactionWithAccount(account)
	if err != nil {
		return nil, err
	}
	return txn.HexString()
}

func (t *Transaction) SignedTransactionWithAccount(account base.Account) (signedTxn base.SignedTransaction, err error) {
	if len(t.inputs) == 0 || len(t.outputs) == 0 {
		return nil, errors.New("invalid inputs or outputs")
	}

	btcAcc, ok := account.(*Account)
	if !ok {
		return nil, base.ErrInvalidAccountType
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)
	for _, input := range t.inputs {
		txIn := wire.NewTxIn(input.outPoint, nil, nil)
		tx.TxIn = append(tx.TxIn, txIn)
		prevOutFetcher.AddPrevOut(*input.outPoint, input.prevOut)
	}
	for _, output := range t.outputs {
		tx.TxOut = append(tx.TxOut, output)
	}

	privateKey := btcAcc.privateKey
	isComing := btcAcc.addressType == AddressTypeComingTaproot
	err = Sign(tx, privateKey, prevOutFetcher, isComing)
	if err != nil {
		return nil, err
	}
	return &SignedTransaction{
		msgTx: tx,
	}, nil
}

func (t *SignedTransaction) HexString() (res *base.OptionalString, err error) {
	var buf bytes.Buffer
	if err := t.msgTx.Serialize(&buf); err != nil {
		return nil, err
	}
	str := hexutil.HexEncodeToString(buf.Bytes())
	return base.NewOptionalString(str), nil
}

func Sign(tx *wire.MsgTx, privKey *btcec.PrivateKey, prevOutFetcher *txscript.MultiPrevOutFetcher, isComing bool) error {
	for i, in := range tx.TxIn {
		prevOut := prevOutFetcher.FetchPrevOutput(in.PreviousOutPoint)
		txSigHashes := txscript.NewTxSigHashes(tx, prevOutFetcher)
		if isComing {
			witness, err := coming_TaprootWitnessSignature(tx, txSigHashes, i, prevOut.Value, prevOut.PkScript, txscript.SigHashDefault, privKey)
			if err != nil {
				return err
			}
			in.Witness = witness
		} else if txscript.IsPayToTaproot(prevOut.PkScript) {
			witness, err := txscript.TaprootWitnessSignature(tx, txSigHashes, i, prevOut.Value, prevOut.PkScript, txscript.SigHashDefault, privKey)
			if err != nil {
				return err
			}
			in.Witness = witness
		} else if txscript.IsPayToPubKeyHash(prevOut.PkScript) {
			sigScript, err := txscript.SignatureScript(tx, i, prevOut.PkScript, txscript.SigHashAll, privKey, true)
			if err != nil {
				return err
			}
			in.SignatureScript = sigScript
		} else {
			pubKeyBytes := privKey.PubKey().SerializeCompressed()
			script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
			if err != nil {
				return err
			}
			amount := prevOut.Value
			witness, err := txscript.WitnessSignature(tx, txSigHashes, i, amount, script, txscript.SigHashAll, privKey, true)
			if err != nil {
				return err
			}
			in.Witness = witness

			if txscript.IsPayToScriptHash(prevOut.PkScript) {
				redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
				if err != nil {
					return err
				}
				in.SignatureScript = append([]byte{byte(len(redeemScript))}, redeemScript...)
			}
		}
	}

	return nil
}

func PayToPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(pubKeyHash).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
		Script()
}

func PayToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}

// custom from txscript.TaprootWitnessSignature(...)
func coming_TaprootWitnessSignature(tx *wire.MsgTx, sigHashes *txscript.TxSigHashes, idx int,
	amt int64, pkScript []byte, hashType txscript.SigHashType,
	key *btcec.PrivateKey) (wire.TxWitness, error) {

	// As we're assuming this was a BIP 86 key, we use an empty root hash
	// which means output key commits to just the public key.
	fakeTapscriptRootHash := []byte{}

	sig, err := coming_RawTxInTaprootSignature(
		tx, sigHashes, idx, amt, pkScript, fakeTapscriptRootHash,
		hashType, key,
	)
	if err != nil {
		return nil, err
	}

	// The witness script to spend a taproot input using the key-spend path
	// is just the signature itself, given the public key is
	// embedded in the previous output script.
	return wire.TxWitness{sig}, nil
}

// custom from txscript.RawTxInTaprootSignature(...)
func coming_RawTxInTaprootSignature(tx *wire.MsgTx, sigHashes *txscript.TxSigHashes, idx int,
	amt int64, pkScript []byte, tapScriptRootHash []byte, hashType txscript.SigHashType,
	key *btcec.PrivateKey) ([]byte, error) {

	// First, we'll start by compute the top-level taproot sighash.
	sigHash, err := txscript.CalcTaprootSignatureHash(
		sigHashes, hashType, tx, idx,
		txscript.NewCannedPrevOutputFetcher(pkScript, amt),
	)
	if err != nil {
		return nil, err
	}

	// Before we sign the sighash, we'll need to apply the taptweak to the
	// private key based on the tapScriptRootHash.
	privKeyTweak := key

	// With the sighash constructed, we can sign it with the specified
	// private key.
	signature, err := schnorr.Sign(privKeyTweak, sigHash)
	if err != nil {
		return nil, err
	}

	sig := signature.Serialize()

	// If this is sighash default, then we can just return the signature
	// directly.
	if hashType == txscript.SigHashDefault {
		return sig, nil
	}

	// Otherwise, append the sighash type to the final sig.
	return append(sig, byte(hashType)), nil
}
