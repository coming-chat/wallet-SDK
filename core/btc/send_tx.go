package btc

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/wallet-SDK/core/base"
)

// Broadcast transactions to the chain
// @return transaction hash
func sendRawTransaction(signedTx, chainnet string) (string, error) {
	client, err := rpcClientOf(chainnet)
	if err != nil {
		return "", err
	}

	tx, err := decodeTx(signedTx)
	if err != nil {
		return "", err
	}

	hash, err := client.SendRawTransaction(tx, false)
	if err != nil {
		return "", base.MapAnyToBasicError(err)
	}

	return hash.String(), nil
}

func decodeTx(txHex string) (*wire.MsgTx, error) {
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
