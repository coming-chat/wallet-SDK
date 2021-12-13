package wallet

import (
	"encoding/hex"
	"errors"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/sr25519"
	"wallet-SDK/chainxTypes"
)

var (
	ErrNilKey      = errors.New("no mnemonic or private key")
	ErrNilMetadata = errors.New("no metadata")
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

func (w *Wallet) sign(message []byte) ([]byte, error) {
	if w.key == nil {
		return nil, ErrNilKey
	}
	return signature.Sign(message, w.key.URI)
}

func (w *Wallet) Sign(message string) (string, error) {
	signed, err := w.sign([]byte(message))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signed), nil
}

func (w *Wallet) GetPublicKey() (string, error) {
	if w.key == nil {
		return "", ErrNilKey
	}
	return "0x" + hex.EncodeToString(w.key.PublicKey), nil
}

func (w *Wallet) GetAddress() (string, error) {
	if w.key == nil {
		return "", ErrNilKey
	}
	return w.key.Address, nil
}

func (w *Wallet) GetPrivateKey() (string, error) {
	if w.key == nil {
		return "", ErrNilKey
	}

	scheme := sr25519.Scheme{}
	kyr, err := subkey.DeriveKeyPair(scheme, w.key.URI)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(kyr.Seed()), nil
}

func (w *Wallet) SignAndGetSendTx(transaction *Transaction) (string, error) {
	var tx string
	signDataByte, err := types.EncodeToBytes(transaction.payload)
	if err != nil {
		return "", err
	}

	signedDataByte, err := w.sign(signDataByte)
	if err != nil {
		return "", err
	}

	signatureData := types.NewSignature(signedDataByte)

	if transaction.extrinsicV3 != nil {
		extSigV3 := chainxTypes.ExtrinsicSignatureV3{
			Signer:    types.NewAddressFromAccountID(w.key.PublicKey),
			Signature: types.MultiSignature{IsSr25519: true, AsSr25519: signatureData},
			Era:       *transaction.extrinsicEra,
			Nonce:     transaction.signatureOptions.Nonce,
			Tip:       transaction.signatureOptions.Tip,
		}
		transaction.extrinsicV3.Signature = extSigV3

		// mark the extrinsic as signed
		transaction.extrinsicV3.Version |= types.ExtrinsicBitSigned
		tx, err = types.EncodeToHexString(transaction.extrinsicV3)
	} else {
		extSigV4 := types.ExtrinsicSignatureV4{
			Signer:    types.NewMultiAddressFromAccountID(w.key.PublicKey),
			Signature: types.MultiSignature{IsSr25519: true, AsSr25519: signatureData},
			Era:       *transaction.extrinsicEra,
			Nonce:     transaction.signatureOptions.Nonce,
			Tip:       transaction.signatureOptions.Tip,
		}
		transaction.extrinsicV4.Signature = extSigV4

		// mark the extrinsic as signed
		transaction.extrinsicV4.Version |= types.ExtrinsicBitSigned
		tx, err = types.EncodeToHexString(transaction.extrinsicV4)
	}

	if err != nil {
		return "", err
	}
	return tx, nil
}
