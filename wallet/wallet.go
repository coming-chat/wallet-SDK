package wallet

import (
	"encoding/json"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/sr25519"
)

type Wallet struct {
	key      *signature.KeyringPair
	keystore *keystore
}

func NewWallet(seedOrPhrase string) (*Wallet, error) {
	network := 44
	if len(seedOrPhrase) == 0 {
		return nil, ErrSeedOrPhrase
	}
	keyringPair, err := signature.KeyringPairFromSecret(seedOrPhrase, uint8(network))
	if err != nil {
		return nil, err
	}
	return &Wallet{
		key: &keyringPair,
	}, nil
}

func NewWalletFromKeyStore(keyStoreJson string, password string) (*Wallet, error) {
	var keyStore keystore
	err := json.Unmarshal([]byte(keyStoreJson), &keyStore)
	if err != nil {
		return nil, err
	}

	if !keyStore.checkPassword(password) {
		return nil, ErrPassword
	}

	return &Wallet{
		keystore: &keyStore,
	}, nil
}

func (w *Wallet) CheckPassword(password string) (bool, error) {
	if w.keystore == nil {
		return false, ErrNilKeystore
	}
	return w.keystore.checkPassword(password), nil
}

func (w *Wallet) Sign(message []byte, password string) ([]byte, error) {
	if w.key != nil {
		return signature.Sign(message, w.key.URI)
	} else if w.keystore != nil {
		return w.keystore.Sign(password, message)
	}
	return nil, ErrNilWallet
}

func (w *Wallet) GetPublicKey() ([]byte, error) {
	if w.key != nil {
		return w.key.PublicKey, nil
	} else if w.keystore != nil {
		publicKey, err := AddressToPublicKey(w.keystore.Address)
		if err != nil {
			return nil, err
		}
		return types.HexDecodeString(publicKey)
	}

	return nil, ErrNilWallet
}

func (w *Wallet) GetPublicKeyHex() (string, error) {
	if w.key != nil {
		return types.HexEncodeToString(w.key.PublicKey), nil
	} else if w.keystore != nil {
		return AddressToPublicKey(w.keystore.Address)
	}

	return "", ErrNilWallet
}

func (w *Wallet) GetAddress(network int) (string, error) {
	if w.key != nil {
		return PublicKeyToAddress(ByteToHex(w.key.PublicKey), network)
	} else if w.keystore != nil {
		publicKey, err := AddressToPublicKey(w.keystore.Address)
		if err != nil {
			return "", err
		}
		return PublicKeyToAddress(publicKey, network)
	}
	return "", ErrNilWallet
}

func (w *Wallet) GetPrivateKeyHex() (string, error) {
	if w.key == nil {
		return "", ErrNilKey
	}

	scheme := sr25519.Scheme{}
	kyr, err := subkey.DeriveKeyPair(scheme, w.key.URI)
	if err != nil {
		return "", err
	}
	return types.HexEncodeToString(kyr.Seed()), nil
}
