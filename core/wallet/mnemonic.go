package wallet

import (
	"errors"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
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

// ExtendMasterKey derives a master key from the given mnemonic and chain network identifier.
//
// Parameters:
//   - mnemonic: A string representing the mnemonic phrase used to generate the seed.
//   - chainnet: The blockchain network, which must be either
//     "mainnet", "testnet", "signet", "simnet" or "regtest".
func ExtendMasterKey(mnemonic string, chainnet string) (string, error) {
	var net chaincfg.Params
	switch chainnet {
	case "mainnet", "bitcoin":
		net = chaincfg.MainNetParams
	case "testnet", "testnet3":
		net = chaincfg.TestNet3Params
	case "signet":
		net = chaincfg.SigNetParams
	case "simnet":
		net = chaincfg.SimNetParams
	case "regtest":
		net = chaincfg.RegressionNetParams
	default:
		return "", errors.New("invalid chainnet")
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return "", err
	}
	masterKey, err := hdkeychain.NewMaster(seed, &net)
	if err != nil {
		return "", err
	}
	return masterKey.String(), nil
}
