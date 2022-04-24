package wallet

import (
	"github.com/tyler-smith/go-bip39"
)

func GenMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	return mnemonic, err
}

func IsValidMnemonic(mnemonic string) bool {
	_, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	return err == nil
}
