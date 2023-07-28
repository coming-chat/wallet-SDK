package btc

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/wallet-SDK/core/base"
	"regexp"
)

func ExtractPsbtToMsgTx(psbtTx string) (*wire.MsgTx, error) {
	packet, err := DecodePsbtTxToPacket(psbtTx)
	if err != nil {
		return nil, err
	}
	return PsbtPacketToMsgTx(packet)
}

func PsbtPacketToMsgTx(packet *psbt.Packet) (*wire.MsgTx, error) {
	if !packet.IsComplete() {
		for i := range packet.Inputs {
			err := psbt.Finalize(packet, i)
			if err != nil {
				return nil, err
			}
		}
	}
	return psbt.Extract(packet)

}

func DecodePsbtTxToPacket(encode string) (*psbt.Packet, error) {
	hexReg := regexp.MustCompile("^[a-fA-F0-9]+$")
	b64Reg := regexp.MustCompile("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$")
	var (
		encodeRaw []byte
		isB64     bool
		err       error
	)
	switch {
	case hexReg.MatchString(encode):
		encodeRaw, err = hex.DecodeString(encode)
		if err != nil {
			return nil, err
		}
	case b64Reg.MatchString(encode):
		encodeRaw = []byte(encode)
		isB64 = true
	default:
		return nil, ErrPsbtEncode
	}
	packet, err := psbt.NewFromRawBytes(bytes.NewReader(encodeRaw), isB64)
	if err != nil {
		return nil, err
	}
	return packet, nil
}

// SignPSBTTx just support segwit v0 & v1(Taproot)
func SignPSBTTx(tx *psbt.Packet, account *Account) error {
	updater, err := psbt.NewUpdater(tx)
	if err != nil {
		return err
	}
	prevOuts := make(map[wire.OutPoint]*wire.TxOut)
	for i, v := range tx.UnsignedTx.TxIn {
		prevOuts[v.PreviousOutPoint] = wire.NewTxOut(tx.Inputs[i].WitnessUtxo.Value, tx.Inputs[i].WitnessUtxo.PkScript)
	}
	privateKey, _ := account.PrivateKey()
	prk, puk := btcec.PrivKeyFromBytes(privateKey)
	for i := range tx.Inputs {
		sigHashes := txscript.NewTxSigHashes(tx.UnsignedTx, txscript.NewMultiPrevOutFetcher(prevOuts))
		scriptClass, _, _, err := txscript.ExtractPkScriptAddrs(tx.Inputs[i].WitnessUtxo.PkScript, account.chain)
		if err != nil {
			return err
		}
		switch scriptClass {
		case txscript.WitnessV1TaprootTy:
			updater.Upsbt.Inputs[i].TaprootInternalKey = schnorr.SerializePubKey(puk)
			sig, err := txscript.TaprootWitnessSignature(tx.UnsignedTx, sigHashes,
				i, tx.Inputs[i].WitnessUtxo.Value, tx.Inputs[i].WitnessUtxo.PkScript, txscript.SigHashDefault, prk)
			if err != nil {
				return err
			}

			updater.Upsbt.Inputs[i].TaprootKeySpendSig = sig[0]
		case txscript.WitnessV0PubKeyHashTy:
			sig, err := txscript.RawTxInWitnessSignature(tx.UnsignedTx, sigHashes, i,
				tx.Inputs[i].WitnessUtxo.Value, tx.Inputs[i].WitnessUtxo.PkScript,
				txscript.SigHashAll, prk)
			if err != nil {
				return err
			}
			success, err := updater.Sign(i, sig, puk.SerializeCompressed(), nil, nil)
			if err != nil {
				return err
			}
			if success != psbt.SignSuccesful {
				return err
			}
		default:
			return ErrPsbtUnsupportedAccountType
		}
		// NOTE: Do not finalize to support multi sign
	}
	return nil
}

// Broadcast transactions to the chain
// @return transaction hash
func sendRawTransaction(signedTx, chainnet string) (string, error) {
	client, err := rpcClientOf(chainnet)
	if err != nil {
		return "", err
	}

	tx, err := DecodeTx(signedTx)
	if err != nil {
		return "", err
	}

	hash, err := client.SendRawTransaction(tx, false)
	if err != nil {
		return "", base.MapAnyToBasicError(err)
	}

	return hash.String(), nil
}

func DecodeTx(txHex string) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	raw, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	if err = tx.Deserialize(bytes.NewReader(raw)); err != nil {
		return nil, err
	}
	return tx, nil
}

// Deprecated: SendRawTransaction is deprecated. Please Use Chain.SendRawTransaction() instead.
func SendRawTransaction(signedTx string, chainnet string) (string, error) {
	return sendRawTransaction(signedTx, chainnet)
}
