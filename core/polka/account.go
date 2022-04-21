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

type Account struct {
	keypair  *signature.KeyringPair
	keystore *wallet.Keystore
	util     *Util
}

func NewAccountWithMnemonic(mnemonic string, network int) (*Account, error) {
	util := NewUtilWithNetwork(network)

	if len(mnemonic) == 0 {
		return nil, wallet.ErrSeedOrPhrase
	}
	keyringPair, err := signature.KeyringPairFromSecret(mnemonic, uint8(network))
	if err != nil {
		return nil, err
	}
	return &Account{
		keypair: &keyringPair,
		util:    util,
	}, nil
}

func NewAccountWithKeystore(keystoreString, password string, network int) (*Account, error) {
	util := NewUtilWithNetwork(network)

	var keyStore wallet.Keystore
	err := json.Unmarshal([]byte(keystoreString), &keyStore)
	if err != nil {
		return nil, err
	}

	if err = keyStore.CheckPassword(password); err != nil {
		return nil, err
	}

	return &Account{
		keystore: &keyStore,
		util:     util,
	}, nil
}

// MARK - Implement the protocol wallet.Account

// If the account generated using keystore, it will return empty
// @return privateKey that will start with 0x.
func (a *Account) PrivateKey() string {
	if a.keypair == nil {
		return ""
	}

	scheme := sr25519.Scheme{}
	kyr, err := subkey.DeriveKeyPair(scheme, a.keypair.URI)
	if err != nil {
		return ""
	}
	return types.HexEncodeToString(kyr.Seed())
}

// @return publicKey that will start with 0x.
func (a *Account) PublicKey() string {
	if a.keypair != nil {
		return types.HexEncodeToString(a.keypair.PublicKey)
	} else if a.keystore != nil {
		pub, err := a.util.DecodeAddressToPublicKey(a.keystore.Address)
		if err != nil {
			return ""
		}
		return pub
	}
	return ""
}

func (a *Account) Address() string {
	address, err := a.util.EncodePublicKeyToAddress(a.PublicKey())
	if err != nil {
		return ""
	}
	return address
}

// TODO: function not implement yet.
func (a *Account) SignData(data []byte, password string) (string, error) {
	return "", errors.New("TODO: function not implement yet.")
}

// TODO: function not implement yet.
func (a *Account) SignHexData(hex string, password string) (string, error) {
	return "", errors.New("TODO: function not implement yet.")
}

// Only available to accounts generated with keystore.
// @return If the password is correct, will return nil
func (a *Account) CheckPassword(password string) error {
	if a.keypair != nil {
		return nil
	}
	if a.keystore == nil {
		return wallet.ErrNilKeystore
	}
	if err := a.keystore.CheckPassword(password); err != nil {
		return err
	}
	return nil
}

// MARK - Implement the protocol wallet.Util

// @param publicKey can start with 0x or not.
func (a *Account) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return a.util.EncodePublicKeyToAddress(publicKey)
}

// Warning: Btc cannot support decode address to public key
func (a *Account) DecodeAddressToPublicKey(address string) (string, error) {
	return a.util.DecodeAddressToPublicKey(address)
}

func (a *Account) IsValidAddress(address string) bool {
	return a.util.IsValidAddress(address)
}
