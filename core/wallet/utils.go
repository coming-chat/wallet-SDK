package wallet

import "github.com/centrifuge/go-substrate-rpc-client/v4/types"

func ByteToHex(data []byte) string {
	return types.HexEncodeToString(data)
}

func HexToByte(hex string) ([]byte, error) {
	return types.HexDecodeString(hex)
}
