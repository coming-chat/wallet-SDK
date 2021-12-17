package wallet

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/sr25519"
)

type Wallet struct {
	key *signature.KeyringPair
}

func NewWallet(seedOrPhrase string, network int) (*Wallet, error) {
	keyringPair, err := signature.KeyringPairFromSecret(seedOrPhrase, uint8(network))
	if err != nil {
		return nil, err
	}
	return &Wallet{
		key: &keyringPair,
	}, nil
}

func NewWalletByKeyStore(keyStoreJson string) (*Wallet, error) {

}

func (w *Wallet) Sign(message []byte) ([]byte, error) {
	if w.key == nil {
		return nil, ErrNilKey
	}
	return signature.Sign(message, w.key.URI)
}

func (w *Wallet) GetPublicKey() ([]byte, error) {
	if w.key == nil {
		return nil, ErrNilKey
	}
	return w.key.PublicKey, nil
}

func (w *Wallet) GetPublicKeyHex() (string, error) {
	if w.key == nil {
		return "", ErrNilKey
	}
	return types.HexEncodeToString(w.key.PublicKey), nil
}

func (w *Wallet) GetAddress() (string, error) {
	if w.key == nil {
		return "", ErrNilKey
	}
	return w.key.Address, nil
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
