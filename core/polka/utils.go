package polka

import (
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
