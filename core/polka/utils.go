package polka

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/itering/subscan/util/base58"
	"github.com/itering/subscan/util/ss58"
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

func Verify(publicKey [32]byte, msg []byte, signature []byte) bool {
	var sigs [64]byte
	copy(sigs[:], signature)
	sig := new(schnorrkel.Signature)
	if err := sig.Decode(sigs); err != nil {
		return false
	}
	publicKeyD := schnorrkel.NewPublicKey(publicKey)
	return publicKeyD.Verify(sig, schnorrkel.NewSigningContext([]byte("substrate"), msg))
}

func VerifyWithPublicHex(publicKey string, msg []byte, signature []byte) bool {
	publicData, err := types.HexDecodeString(publicKey)
	if err != nil {
		return false
	}

	var public32 [32]byte
	copy(public32[:], publicData)
	return Verify(public32, msg, signature)
}
