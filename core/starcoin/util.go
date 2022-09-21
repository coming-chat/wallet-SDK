package starcoin

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
	"github.com/starcoinorg/starcoin-go/types"
)

func NewStructTag(tag string) (*types.StructTag, error) {
	if strings.Contains(tag, "<") {
		return nil, errors.New("Not implemented")
	}

	parts := strings.Split(tag, "::")
	if len(parts) != 3 {
		return nil, errors.New("Invalid struct tag string literal.")
	}
	addr, err := NewAccountAddressFromHex(parts[0])
	if err != nil {
		return nil, err
	}
	return &types.StructTag{
		Address: *addr,
		Module:  types.Identifier(parts[1]),
		Name:    types.Identifier(parts[2]),
	}, nil
}

func NewAccountAddressFromHex(addr string) (*types.AccountAddress, error) {
	if strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X") {
		addr = addr[2:]
	}
	if len(addr)%2 != 0 {
		addr = "0" + addr
	}

	bytes, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	if len(bytes) > addressLength {
		return nil, fmt.Errorf("Hex string is too long. Address's length is %v bytes.", addressLength)
	}

	res := types.AccountAddress{}
	copy(res[addressLength-len(bytes):], bytes[:])
	return &res, nil
}

func StructTagToString(t types.StructTag) string {
	return fmt.Sprintf("0x%v::%v::%v", hex.EncodeToString(t.Address[:]), t.Module, t.Name)
}

func NewU128FromString(number string) (*serde.Uint128, error) {
	n, ok := big.NewInt(0).SetString(number, 10)
	if !ok {
		return nil, errors.New("Invalid U128: not a number")
	}
	if n.Sign() < 0 {
		return nil, errors.New("Invalid U128: negative number")
	}
	bitLen := n.BitLen()
	if bitLen > 128 {
		return nil, errors.New("Invalid U128: too large number")
	}
	res := &serde.Uint128{}
	res.Low = n.Uint64()
	if bitLen > 64 {
		n = n.Rsh(n, 64)
		res.High = n.Uint64()
	}
	return res, nil
}
