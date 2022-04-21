package polka

import (
	"encoding/json"
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/wallet"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/sr25519"
)

type RootAccount struct {
	keypair  *signature.KeyringPair
	keystore *wallet.Keystore
	rootUtil *Util
}

func NewRootAccountWithMnemonic(mnemonic string) (*RootAccount, error) {
	network := 44
	if len(mnemonic) == 0 {
		return nil, wallet.ErrSeedOrPhrase
	}
	keyringPair, err := signature.KeyringPairFromSecret(mnemonic, uint8(network))
	if err != nil {
		return nil, err
	}

	util := NewUtilWithNetwork(44)
	return &RootAccount{
		keypair:  &keyringPair,
		rootUtil: util,
	}, nil
}

func NewRootAccountWithKeystore(keystoreString, password string) (*RootAccount, error) {
	var keyStore wallet.Keystore
	err := json.Unmarshal([]byte(keystoreString), &keyStore)
	if err != nil {
		return nil, err
	}

	if err = keyStore.CheckPassword(password); err != nil {
		return nil, err
	}

	util := NewUtilWithNetwork(44)
	return &RootAccount{
		keystore: &keyStore,
		rootUtil: util,
	}, nil
}

func newRootAccountWithKeystoreObj(keystore *wallet.Keystore) (*RootAccount, error) {
	return &RootAccount{keystore: keystore}, nil
}

// MARK - Implement the protocol wallet.Account

// @return privateKey data
func (a *RootAccount) PrivateKeyData() ([]byte, error) {
	if a.keypair == nil {
		return nil, wallet.ErrNilKey
	}

	scheme := sr25519.Scheme{}
	kyr, err := subkey.DeriveKeyPair(scheme, a.keypair.URI)
	if err != nil {
		return nil, err
	}
	return kyr.Seed(), nil
}

// @return privateKey string that will start with 0x.
func (a *RootAccount) PrivateKey() (string, error) {
	data, err := a.PrivateKeyData()
	if err != nil {
		return "", err
	}
	return types.HexEncodeToString(data), nil
}

// @return publicKey string that will start with 0x.
func (a *RootAccount) PublicKey() string {
	if a.keypair != nil {
		return types.HexEncodeToString(a.keypair.PublicKey)
	} else if a.keystore != nil {
		pub, err := a.rootUtil.DecodeAddressToPublicKey(a.keystore.Address)
		if err != nil {
			return ""
		}
		return pub
	}
	return ""
}

// @return address string
func (a *RootAccount) Address() string {
	address, err := a.rootUtil.EncodePublicKeyToAddress(a.PublicKey())
	if err != nil {
		return ""
	}
	return address
}

// TODO: function not implement yet.
func (a *RootAccount) SignData(data []byte, password string) (string, error) {
	return "", errors.New("TODO: function not implement yet.")
}

// TODO: function not implement yet.
func (a *RootAccount) SignHexData(hex string, password string) (string, error) {
	return "", errors.New("TODO: function not implement yet.")
}
