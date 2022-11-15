package wallet

import (
	"errors"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type WatchAccount struct {
	address string
}

func (a *WatchAccount) PrivateKey() ([]byte, error) {
	return nil, errors.New("The watch wallet cannot get private key")
}

func (a *WatchAccount) PrivateKeyHex() (string, error) {
	return "", errors.New("The watch wallet cannot get private key")
}

func (a *WatchAccount) PublicKey() []byte {
	return nil
}

func (a *WatchAccount) PublicKeyHex() string {
	return ""
}

func (a *WatchAccount) Address() string {
	return a.address
}

func (a *WatchAccount) Sign(message []byte, password string) ([]byte, error) {
	return nil, errors.New("The watch wallet cannot be signed")
}
func (a *WatchAccount) SignHex(messageHex string, password string) (*base.OptionalString, error) {
	return nil, errors.New("The watch wallet cannot be signed")
}
