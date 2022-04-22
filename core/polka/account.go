package polka

import (
	"encoding/json"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/wallet"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/sr25519"
)

type Account struct {
	keypair  *signature.KeyringPair
	keystore *wallet.Keystore

	publicKey string
	address   string

	Network int
}

func NewAccountWithMnemonic(mnemonic string, network int) (*Account, error) {
	if len(mnemonic) == 0 {
		return nil, wallet.ErrSeedOrPhrase
	}
	keyringPair, err := signature.KeyringPairFromSecret(mnemonic, uint8(network))
	if err != nil {
		return nil, err
	}

	publicKey := types.HexEncodeToString(keyringPair.PublicKey)
	address, err := EncodePublicKeyToAddress(publicKey, network)
	if err != nil {
		return nil, err
	}

	return &Account{
		keypair:   &keyringPair,
		publicKey: publicKey,
		address:   address,
		Network:   network,
	}, nil
}

func NewAccountWithKeystore(keystoreString, password string, network int) (*Account, error) {
	var keyStore wallet.Keystore
	err := json.Unmarshal([]byte(keystoreString), &keyStore)
	if err != nil {
		return nil, err
	}
	if err = keyStore.CheckPassword(password); err != nil {
		return nil, err
	}

	publicKey, err := DecodeAddressToPublicKey(keyStore.Address)
	if err != nil {
		return nil, err
	}
	address, err := EncodePublicKeyToAddress(publicKey, network)
	if err != nil {
		return nil, err
	}

	return &Account{
		keystore:  &keyStore,
		publicKey: publicKey,
		address:   address,
		Network:   network,
	}, nil
}

func (a *Account) DeriveAccountAt(network int) (*Account, error) {
	address, err := EncodePublicKeyToAddress(a.publicKey, network)
	if err != nil {
		return nil, err
	}
	return &Account{
		keypair:   a.keypair,
		keystore:  a.keystore,
		publicKey: a.publicKey,
		address:   address,
		Network:   network,
	}, nil
}

// MARK - Implement the protocol wallet.Account

// @return privateKey data
func (a *Account) PrivateKeyData() ([]byte, error) {
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
func (a *Account) PrivateKey() (string, error) {
	data, err := a.PrivateKeyData()
	if err != nil {
		return "", err
	}
	return types.HexEncodeToString(data), nil
}

// @return publicKey string that will start with 0x.
func (a *Account) PublicKey() string {
	return a.publicKey
}

// @return address string
func (a *Account) Address() string {
	return a.address
}

func (a *Account) Sign(message []byte, password string) (data []byte, err error) {
	defer func() {
		errPanic := recover()
		if errPanic != nil {
			err = wallet.ErrSign
			fmt.Println(errPanic)
			return
		}
	}()
	if a.keypair != nil {
		data, err := signature.Sign(message, a.keypair.URI)
		return data, err // Must be separate to ensure that err can catch panic
	} else if a.keystore != nil {
		data, err := a.keystore.Sign(message, password)
		return data, err
	}
	return nil, wallet.ErrNilWallet
}

func (a *Account) SignHex(messageHex string, password string) ([]byte, error) {
	message, err := types.HexDecodeString(messageHex)
	if err != nil {
		return nil, err
	}
	return a.Sign(message, password)
}
