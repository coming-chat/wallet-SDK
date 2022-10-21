package polka

import (
	"crypto/ed25519"
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/itering/subscan/util/base58"
	"github.com/itering/subscan/util/ss58"
	"golang.org/x/crypto/blake2b"
)

func addressStringToMultiAddress(dest string) (types.MultiAddress, error) {
	ss58Format := base58.Decode(dest)
	if len(ss58Format) == 0 {
		return types.MultiAddress{}, ErrAddress
	}
	destPublicKey := ss58.Decode(dest, int(ss58Format[0]))
	if len(destPublicKey) == 0 {
		return types.MultiAddress{}, ErrAddress
	}
	return types.NewMultiAddressFromHexAccountID(destPublicKey)
}

func addressStringToAddress(dest string) (types.Address, error) {
	ss58Format := base58.Decode(dest)
	if len(ss58Format) == 0 {
		return types.Address{}, ErrAddress
	}
	destPublicKey := ss58.Decode(dest, int(ss58Format[0]))
	if len(destPublicKey) == 0 {
		return types.Address{}, ErrAddress
	}
	return types.NewAddressFromHexAccountID(destPublicKey)
}

func ByteToHex(data []byte) string {
	return types.HexEncodeToString(data)
}

func HexToByte(hex string) ([]byte, error) {
	return types.HexDecodeString(hex)
}

type Sr25519Util struct{}

func (s Sr25519Util) IsValidSignature(publicKey, msg, signature []byte) bool {
	if len(msg) > 256 {
		h := blake2b.Sum256(msg)
		msg = h[:]
	}
	var (
		sigs        [64]byte
		fixedPubKey [32]byte
		sig         = new(schnorrkel.Signature)
	)
	copy(fixedPubKey[:], publicKey[:32])
	copy(sigs[:], signature[:64])
	pubKey := schnorrkel.NewPublicKey(fixedPubKey)
	if err := sig.Decode(sigs); err != nil {
		return false
	}
	return pubKey.Verify(sig, schnorrkel.NewSigningContext([]byte("substrate"), msg))
}

type Ed25519 struct{}

func (e *Ed25519) IsValidSignature(publicKey, msg, signature []byte) bool {
	return ed25519.Verify(publicKey, msg, signature)
}
