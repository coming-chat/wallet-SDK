package secp256k1

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

func IsValidSignature(publicKey, msg, signature []byte) bool {
	srSignature, err := schnorr.ParseSignature(signature)
	if err != nil {
		return false
	}
	pubKey, err := btcec.ParsePubKey(publicKey)
	if err != nil {
		return false
	}
	return srSignature.Verify(msg, pubKey)
}
