package btc

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/tyler-smith/go-bip39"
)

func Derivation(mnemonic string, path string) (*btcec.PrivateKey, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	dPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}
	key := masterKey
	for _, n := range dPath {
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}
	return key.ECPrivKey()
}

func ComingPrivateKey(mnemonic string) (*btcec.PrivateKey, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	pri, _ := btcec.PrivKeyFromBytes(seed)
	return pri, nil
}
