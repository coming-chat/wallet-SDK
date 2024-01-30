package eth

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func AbiCoder(abiTypes []string) (abi.Arguments, error) {
	if len(abiTypes) <= 0 {
		return nil, errors.New("no abi types")
	}
	args := make([]abi.Argument, len(abiTypes))
	for idx, typStr := range abiTypes {
		typAbi, err := abi.NewType(typStr, "", nil)
		if err != nil {
			return nil, err
		}
		args[idx] = abi.Argument{Type: typAbi}
	}
	return args, nil
}

// AbiCoderEncode
// usage like js ethers.AbiCoder
// https://docs.ethers.org/v5/api/utils/abi/coder/
func AbiCoderEncode(abiTypes []string, args ...any) ([]byte, error) {
	coder, err := AbiCoder(abiTypes)
	if err != nil {
		return nil, err
	}
	return coder.Pack(args...)
}
