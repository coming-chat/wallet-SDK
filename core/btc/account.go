package btc

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type Account struct {
	privateKey *btcec.PrivateKey
	address    string
	chain      *chaincfg.Params

	// Default is `AddressTypeComingTaproot`
	addressType AddressType
}

func NewAccountWithMnemonic(mnemonic, chainnet string, addressType AddressType) (*Account, error) {
	chain, err := netParamsOf(chainnet)
	if err != nil {
		return nil, err
	}

	var privateKey *btcec.PrivateKey
	path := AddressTypeDerivePath(addressType)
	if path == "--" {
		addressType = AddressTypeComingTaproot
		privateKey, err = ComingPrivateKey(mnemonic)
	} else {
		privateKey, err = Derivation(mnemonic, path)
	}
	if err != nil {
		return nil, err
	}

	address, err := EncodePubKeyToAddress(privateKey.PubKey(), chain, addressType)
	if err != nil {
		return nil, err
	}

	return &Account{
		privateKey:  privateKey,
		address:     address,
		chain:       chain,
		addressType: addressType,
	}, nil
}

func AccountWithPrivateKey(prikey string, chainnet string, addressType AddressType) (*Account, error) {
	chain, err := netParamsOf(chainnet)
	if err != nil {
		return nil, err
	}

	var privateKey *btcec.PrivateKey
	wif, err := btcutil.DecodeWIF(prikey)
	if err != nil {
		seed, err := types.HexDecodeString(prikey)
		if err != nil {
			return nil, err
		}
		privateKey, _ = btcec.PrivKeyFromBytes(seed)
	} else {
		if !wif.IsForNet(chain) {
			return nil, fmt.Errorf("the specified chainnet does not match the wif private key")
		}
		privateKey = wif.PrivKey
	}

	address, err := EncodePubKeyToAddress(privateKey.PubKey(), chain, addressType)
	if err != nil {
		return nil, err
	}

	return &Account{
		privateKey:  privateKey,
		address:     address,
		chain:       chain,
		addressType: addressType,
	}, nil
}

func (a *Account) AddressType() AddressType {
	return a.addressType
}

func (a *Account) DeriveAccountAt(chainnet string) (*Account, error) {
	chain, err := netParamsOf(chainnet)
	if err != nil {
		return nil, err
	}
	address, err := EncodePubKeyToAddress(a.privateKey.PubKey(), chain, a.addressType)
	if err != nil {
		return nil, err
	}
	return &Account{
		privateKey:  a.privateKey,
		address:     address,
		chain:       chain,
		addressType: a.addressType,
	}, nil
}

func (a *Account) AddressTypeString() string {
	return AddressTypeDescription(a.addressType)
}
func (a *Account) DerivePath() string {
	return AddressTypeDerivePath(a.addressType)
}

func (a *Account) AddressWithType(addrType AddressType) (*base.OptionalString, error) {
	addr, err := EncodePubKeyToAddress(a.privateKey.PubKey(), a.chain, addrType)
	if err != nil {
		return nil, err
	}
	return base.NewOptionalString(addr), nil
}

func (a *Account) WIFPrivateKeyString() (*base.OptionalString, error) {
	str, err := PrivateKeyToWIF(a.privateKey, a.chain)
	if err != nil {
		return nil, err
	}
	return base.NewOptionalString(str), nil
}

// MARK - Implement the protocol Account

// @return privateKey data
func (a *Account) PrivateKey() ([]byte, error) {
	return a.privateKey.Serialize(), nil
}

// @return privateKey string that will start with 0x.
func (a *Account) PrivateKeyHex() (string, error) {
	return types.HexEncodeToString(a.privateKey.Serialize()), nil
}

// @return publicKey data
func (a *Account) PublicKey() []byte {
	return a.privateKey.PubKey().SerializeCompressed()
}

// @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	pub := a.privateKey.PubKey().SerializeCompressed()
	return types.HexEncodeToString(pub)
}

func (a *Account) MultiSignaturePubKey() string {
	pub := a.privateKey.PubKey().SerializeUncompressed()
	return types.HexEncodeToString(pub)
}

// @return default is the mainnet address
func (a *Account) Address() string {
	return a.address
}

// TODO: function not implement yet.
func (a *Account) Sign(message []byte, password string) ([]byte, error) {
	return nil, base.ErrUnsupportedFunction
}

// TODO: function not implement yet.
func (a *Account) SignHex(messageHex string, password string) (*base.OptionalString, error) {
	return nil, base.ErrUnsupportedFunction
}

// SignMessage
// https://developer.bitcoin.org/reference/rpc/signmessage.html
// @param msg The message to create a signature of.
// @return The signature of the message encoded in base64.
func (a *Account) SignMessage(msg string) (*base.OptionalString, error) {
	msgHash := messageHash(msg)

	signbytes := ecdsa.SignCompact(a.privateKey, msgHash, true)
	signature := base64.StdEncoding.EncodeToString(signbytes)
	return base.NewOptionalString(signature), nil
}

func VerifySignature(pubkey, message, signature string) bool {
	signBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}
	pubBytes, err := types.HexDecodeString(pubkey)
	if err != nil {
		return false
	}
	pub, err := btcec.ParsePubKey(pubBytes)
	if err != nil {
		return false
	}

	msgHash := messageHash(message)
	recoverPub, ok, err := ecdsa.RecoverCompact(signBytes, msgHash)
	if err != nil || !ok {
		return false
	}

	return pub.IsEqual(recoverPub)
}

func messageHash(msg string) []byte {
	var buf bytes.Buffer
	_ = wire.WriteVarString(&buf, 0, "Bitcoin Signed Message:\n")
	_ = wire.WriteVarString(&buf, 0, msg)
	return chainhash.DoubleHashB(buf.Bytes())
}

func (a *Account) SignPsbt(psbtHex string) (*SignedPsbtTransaction, error) {
	psbt, err := NewPsbtTransaction(psbtHex)
	if err != nil {
		return nil, err
	}
	signedTxn, err := psbt.SignedTransactionWithAccount(a)
	if err != nil {
		return nil, err
	}
	return signedTxn.(*SignedPsbtTransaction), nil
}

// MARK - Implement the protocol AddressUtil

// @param publicKey can start with 0x or not.
func (a *Account) EncodePublicKeyToAddress(publicKey string) (string, error) {
	return EncodePublicKeyToAddress(publicKey, a.chain.Name, AddressTypeComingTaproot)
}

// @return publicKey that will start with 0x.
func (a *Account) DecodeAddressToPublicKey(address string) (string, error) {
	return "", ErrDecodeAddress
}

func (a *Account) IsValidAddress(address string) bool {
	return IsValidAddress(address, a.chain.Name)
}

func AsBitcoinAccount(account base.Account) *Account {
	if r, ok := account.(*Account); ok {
		return r
	} else {
		return nil
	}
}

// PublicKeyTransform
// @param pubkey the original public key, can be uncompressed, compressed, or hybrid
// @param compress the transformed public key should be compressed or not
func PublicKeyTransform(pubkey string, compress bool) (string, error) {
	pubData, err := types.HexDecodeString(pubkey)
	if err != nil {
		return "", err
	}
	btcPubkey, err := btcutil.NewAddressPubKey(pubData, &chaincfg.MainNetParams)
	if err != nil {
		return "", err
	}
	isCompressed := btcPubkey.Format() == btcutil.PKFCompressed
	if isCompressed == compress {
		return types.HexEncodeToString(pubData), nil
	}
	if compress {
		return types.HexEncodeToString(btcPubkey.PubKey().SerializeCompressed()), nil
	} else {
		return types.HexEncodeToString(btcPubkey.PubKey().SerializeUncompressed()), nil
	}
}

func PrivateKeyToWIF(pri *btcec.PrivateKey, network *chaincfg.Params) (string, error) {
	wif, err := btcutil.NewWIF(pri, network, true)
	if err != nil {
		return "", err
	}
	return wif.String(), nil
}
