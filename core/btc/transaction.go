package btc

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/wire"
)

func ExtractPSBTHexToMsgTx(psbtTxHex string) (*wire.MsgTx, error) {
	psbtTx, err := hex.DecodeString(psbtTxHex)
	if err != nil {
		return nil, err
	}
	return ExtractPSBTPacketToMsgTx(psbtTx, false)
}

func ExtractPSBTBse64ToMsgTx(psbtTxBase64 string) (*wire.MsgTx, error) {
	return ExtractPSBTPacketToMsgTx([]byte(psbtTxBase64), true)
}

func ExtractPSBTPacketToMsgTx(psbtTx []byte, isBase64 bool) (*wire.MsgTx, error) {
	psbtPacket, err := psbt.NewFromRawBytes(bytes.NewReader(psbtTx), isBase64)
	if err != nil {
		return nil, err
	}
	if !psbtPacket.IsComplete() {
		for i := range psbtPacket.Inputs {
			err = psbt.Finalize(psbtPacket, i)
			if err != nil {
				return nil, err
			}
		}
	}
	return psbt.Extract(psbtPacket)
}
