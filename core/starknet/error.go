package starknet

import "regexp"

func IsNotDeployedError(err error) bool {
	if err == nil {
		return false
	}
	regInsufficientGas := regexp.MustCompile(`.*UNINITIALIZED_CONTRACT.*contract address 0x[0-9a-zA-Z]+ is not deployed.*`)
	match := regInsufficientGas.FindAllStringSubmatch(err.Error(), -1)
	return len(match) > 0
}
