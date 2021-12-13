package wallet

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/itering/subscan/util/base58"
	"github.com/itering/subscan/util/ss58"
)

func addressStringToMultiAddress(dest string) (types.MultiAddress, error) {
	ss58Format := base58.Decode(dest)
	destPublicKey := ss58.Decode(dest, int(ss58Format[0]))
	return types.NewMultiAddressFromHexAccountID(destPublicKey)
}

func addressStringToAddress(dest string) (types.Address, error) {
	ss58Format := base58.Decode(dest)
	destPublicKey := ss58.Decode(dest, int(ss58Format[0]))
	return types.NewAddressFromHexAccountID(destPublicKey)
}

func AddressToPublicKey(address string) string {
	ss58Format := base58.Decode(address)
	return "0x" + ss58.Decode(address, int(ss58Format[0]))
}

func PublicKeyToAddress(publicKey string, network int) string {
	return ss58.Encode(publicKey, network)
}
