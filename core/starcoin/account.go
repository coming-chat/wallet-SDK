package starcoin

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/go-aptos/crypto/derivation"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/eth"
	"github.com/cosmos/go-bip39"
	starTypes "github.com/starcoinorg/starcoin-go/types"
	"golang.org/x/crypto/sha3"
)

const (
	addressLength = 16
)

type Account struct {
	privateKey ed25519.PrivateKey
	AuthKey    []byte
	address    string
}

func NewAccountWithMnemonic(mnemonic string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	key, err := derivation.DeriveForPath("m/44'/101010'/0'/0'/0'", seed)
	if err != nil {
		return nil, err
	}
	return accountWithKey(key.Key), nil
}

func AccountWithPrivateKey(privateKey string) (*Account, error) {
	key, err := types.HexDecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	return accountWithKey(key), nil
}

func accountWithKey(key []byte) *Account {
	privateKey := ed25519.NewKeyFromSeed(key)

	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKey = append(publicKey, 0x00)
	authKey := sha3.Sum256(publicKey)
	address := hex.EncodeToString(authKey[len(authKey)-addressLength:])
	address = eth.TransformEIP55Address(address)
	return &Account{
		privateKey: privateKey,
		AuthKey:    authKey[:],
		address:    address,
	}
}

// MARK - Implement the protocol Account

// @return privateKey data
func (a *Account) PrivateKey() ([]byte, error) {
	return a.privateKey[:ed25519.SeedSize], nil
}

// @return privateKey string that will start with 0x.
func (a *Account) PrivateKeyHex() (string, error) {
	return types.HexEncodeToString(a.privateKey[:ed25519.SeedSize]), nil
}

// @return publicKey data
func (a *Account) PublicKey() []byte {
	return a.privateKey.Public().(ed25519.PublicKey)
}

// @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	return types.HexEncodeToString(a.PublicKey())
}

func (a *Account) Address() string {
	return a.address
}

func (a *Account) Sign(message []byte, password string) ([]byte, error) {
	return ed25519.Sign(a.privateKey, message), nil
}

func (a *Account) SignHex(messageHex string, password string) (*base.OptionalString, error) {
	msg, err := types.HexDecodeString(messageHex)
	if err != nil {
		return nil, errors.New("Invalid message hex string")
	}
	signature := ed25519.Sign(a.privateKey, msg)
	signString := hex.EncodeToString(signature)
	return &base.OptionalString{Value: signString}, nil
}

// MARK - Implement the protocol AddressUtil

// @param publicKey can start with 0x or not.
func (a *Account) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return EncodePublicKeyToAddress(publicKey)
}

// @return publicKey that will start with 0x.
func (a *Account) DecodeAddressToPublicKey(address string) (string, error) {
	return DecodeAddressToPublicKey(address)
}

func (a *Account) IsValidAddress(address string) bool {
	return IsValidAddress(address)
}

func (a *Account) AccountAddress() starTypes.AccountAddress {
	addr := starTypes.AccountAddress{}
	copy(addr[:], a.AuthKey[len(a.AuthKey)-addressLength:])
	return addr
}

func (a *Account) StarcoinPrivateKey() starTypes.Ed25519PrivateKey {
	return starTypes.Ed25519PrivateKey(a.privateKey[:ed25519.SeedSize])
}

func (a *Account) StarcoinPublicKey() starTypes.Ed25519PublicKey {
	return starTypes.Ed25519PublicKey(a.PublicKey())
}

func AsStarcoinAccount(account base.Account) *Account {
	if r, ok := account.(*Account); ok {
		return r
	} else {
		return nil
	}
}
