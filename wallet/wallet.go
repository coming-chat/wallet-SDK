package wallet

import (
	"encoding/hex"
	"errors"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/vedhavyas/go-subkey"
	"github.com/vedhavyas/go-subkey/sr25519"
)

var (
	ErrNilKey      = errors.New("no mnemonic or private key")
	ErrNilMetadata = errors.New("no metadata")
)

type Wallet struct {
	metadata *types.Metadata
	key      *signature.KeyringPair
}

func NewWallet(metadataString, seedOrPhrase string, network int) (*Wallet, error) {
	var metadata types.Metadata
	if err := types.DecodeFromHexString(metadataString, &metadata); err != nil {
		return nil, err
	}

	keyringPair, err := signature.KeyringPairFromSecret(seedOrPhrase, uint8(network))
	if err != nil {
		return nil, err
	}
	return &Wallet{
		metadata: &metadata,
		key:      &keyringPair,
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
	signDataByte, err := types.EncodeToBytes(transaction.payload)
	if err != nil {
		return "", err
	}

	signerPubKey := types.NewMultiAddressFromAccountID(w.key.PublicKey)
	signedDataByte, err := w.sign(signDataByte)
	if err != nil {
		return "", err
	}

	signatureData := types.NewSignature(signedDataByte)
	extSig := types.ExtrinsicSignatureV4{
		Signer:    signerPubKey,
		Signature: types.MultiSignature{IsSr25519: true, AsSr25519: signatureData},
		Era:       *transaction.extrinsicEra,
		Nonce:     transaction.signatureOptions.Nonce,
		Tip:       transaction.signatureOptions.Tip,
	}

	transaction.extrinsic.Signature = extSig

	// mark the extrinsic as signed
	transaction.extrinsic.Version |= types.ExtrinsicBitSigned
	tx, err := types.EncodeToHexString(transaction.extrinsic)
	if err != nil {
		return "", err
	}
	return tx, nil
}
