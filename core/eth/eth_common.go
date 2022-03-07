package eth

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

func (e *EthChain) UnpackParams(out interface{}, inputs abi.Arguments, paramsStr string) error {
	paramsStr = strings.TrimPrefix(paramsStr, "0x")
	data, err := hex.DecodeString(paramsStr)
	if err != nil {
		return err
	}
	a, err := inputs.Unpack(data)
	if err != nil {
		return err
	}
	err = inputs.Copy(out, a)
	if err != nil {
		return err
	}
	return nil
}

// 不带 0x 前缀
func (e *EthChain) PackParams(inputs abi.Arguments, args ...interface{}) (string, error) {
	bytes_, err := inputs.Pack(args...)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes_), nil
}

func (e *EthChain) MethodIdFromMethodStr(methodStr string) string {
	return hex.EncodeToString(crypto.Keccak256([]byte(methodStr))[:4])
}

func (e *EthChain) MethodFromPayload(abiStr string, payloadStr string) (*abi.Method, error) {
	if len(payloadStr) < 8 {
		return nil, errors.New("payloadStr error")
	}

	payloadStr = strings.TrimPrefix(payloadStr, "0x")

	parsedAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return nil, err
	}
	data, err := hex.DecodeString(payloadStr)
	if err != nil {
		return nil, err
	}
	method, err := parsedAbi.MethodById(data[:4])
	if err != nil {
		return nil, err
	}
	return method, err
}
