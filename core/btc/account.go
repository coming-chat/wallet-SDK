package btc

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/tyler-smith/go-bip39"
)

type AddressType = base.SDKEnumInt

const (
	AddressTypeComingTaproot AddressType = 0
	AddressTypeNativeSegwit  AddressType = 1
	AddressTypeNestedSegwit  AddressType = 2
	AddressTypeTaproot       AddressType = 3
	AddressTypeLegacy        AddressType = 4
)

type Account struct {
	privateKey *btcec.PrivateKey
	address    *btcutil.AddressPubKey
	chain      *chaincfg.Params

	// Default is `AddressTypeComingTaproot`
	AddressType AddressType
}

func NewAccountWithMnemonic(mnemonic, chainnet string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	pri, pub := btcec.PrivKeyFromBytes(seed)
	chain, err := netParamsOf(chainnet)
	if err != nil {
		return nil, err
	}
	address, err := btcutil.NewAddressPubKey(pub.SerializeCompressed(), chain)
	if err != nil {
		return nil, err
	}

	return &Account{
		privateKey: pri,
		address:    address,
		chain:      chain,
	}, nil
}

func AccountWithPrivateKey(prikey string, chainnet string) (*Account, error) {
	var (
		pri     *btcec.PrivateKey
		pubData []byte
	)
	chain, err := netParamsOf(chainnet)
	if err != nil {
		return nil, err
	}
	wif, err := btcutil.DecodeWIF(prikey)
	if err != nil {
		seed, err := types.HexDecodeString(prikey)
		if err != nil {
			return nil, err
		}
		var pub *btcec.PublicKey
		pri, pub = btcec.PrivKeyFromBytes(seed)
		pubData = pub.SerializeCompressed()
	} else {
		if !wif.IsForNet(chain) {
			return nil, fmt.Errorf("the specified chainnet does not match the wif private key")
		}
		pri = wif.PrivKey
		pubData = wif.SerializePubKey()
	}

	address, err := btcutil.NewAddressPubKey(pubData, chain)
	if err != nil {
		return nil, err
	}
	return &Account{
		privateKey: pri,
		address:    address,
		chain:      chain,
	}, nil
}

// NativeSegwitAddress P2WPKH just for m/84'/
func (a *Account) NativeSegwitAddress() (*base.OptionalString, error) {
	address, err := btcutil.NewAddressWitnessPubKeyHash(a.address.AddressPubKeyHash().ScriptAddress(), a.chain)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: address.EncodeAddress()}, nil
}

// NestedSegwitAddress P2SH-P2WPKH just for m/49'/
func (a *Account) NestedSegwitAddress() (*base.OptionalString, error) {
	witAddr, err := btcutil.NewAddressWitnessPubKeyHash(a.address.AddressPubKeyHash().ScriptAddress(), a.chain)
	if err != nil {
		return nil, err
	}
	witnessProgram, err := txscript.PayToAddrScript(witAddr)
	if err != nil {
		return nil, err
	}
	address, err := btcutil.NewAddressScriptHash(witnessProgram, a.chain)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: address.EncodeAddress()}, nil
}

// TaprootAddress P2TR just for m/86'/
func (a *Account) TaprootAddress() (*base.OptionalString, error) {
	tapKey := txscript.ComputeTaprootKeyNoScript(a.address.PubKey())
	address, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(tapKey), a.chain,
	)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: address.EncodeAddress()}, nil
}

func (a *Account) ComingTaprootAddress() (*base.OptionalString, error) {
	taproot, err := btcutil.NewAddressTaproot(a.address.ScriptAddress()[1:33], a.chain)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: taproot.EncodeAddress()}, nil
}

// LegacyAddress P2PKH just for m/44'/
func (a *Account) LegacyAddress() (*base.OptionalString, error) {
	return &base.OptionalString{Value: a.address.AddressPubKeyHash().EncodeAddress()}, nil
}

func (a *Account) DeriveAccountAt(chainnet string) (*Account, error) {
	chain, err := netParamsOf(chainnet)
	if err != nil {
		return nil, err
	}
	address, err := btcutil.NewAddressPubKey(a.address.ScriptAddress(), chain)
	if err != nil {
		return nil, err
	}
	return &Account{
		privateKey: a.privateKey,
		address:    address,
		chain:      chain,
	}, nil
}

func (a *Account) AddressTypeString() string {
	switch a.AddressType {
	case AddressTypeComingTaproot:
		return "Coming Taproot"
	case AddressTypeNativeSegwit:
		return "Native Segwit (P2WPKH)"
	case AddressTypeNestedSegwit:
		return "Nested Segwit (P2SH-P2WPKH)"
	case AddressTypeTaproot:
		return "Taproot (P2TR)"
	case AddressTypeLegacy:
		return "Legacy (P2PKH)"
	}
	return "--"
}
func (a *Account) DerivePath() string {
	switch a.AddressType {
	case AddressTypeComingTaproot:
		return "--"
	case AddressTypeNativeSegwit:
		return "m/84'/0'/0/0"
	case AddressTypeNestedSegwit:
		return "m/49'/0'/0/0"
	case AddressTypeTaproot:
		return "m/86'/0'/0/0"
	case AddressTypeLegacy:
		return "m/44'/0'/0/0"
	}
	return "--"
}

func (a *Account) AddressWithType(addrType AddressType) (*base.OptionalString, error) {
	switch addrType {
	case AddressTypeComingTaproot:
		return a.ComingTaprootAddress()
	case AddressTypeNativeSegwit:
		return a.NativeSegwitAddress()
	case AddressTypeNestedSegwit:
		return a.NestedSegwitAddress()
	case AddressTypeTaproot:
		return a.TaprootAddress()
	case AddressTypeLegacy:
		return a.LegacyAddress()
	}
	return nil, fmt.Errorf("unknow address type `%v`", addrType)
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
	return a.address.ScriptAddress()
}

// @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	return types.HexEncodeToString(a.address.ScriptAddress())
}

func (a *Account) MultiSignaturePubKey() string {
	return types.HexEncodeToString(a.address.PubKey().SerializeUncompressed())
}

// @return default is the mainnet address
func (a *Account) Address() string {
	addr, err := a.AddressWithType(a.AddressType)
	if err != nil {
		return "--"
	}
	return addr.Value
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

	signbytes, err := ecdsa.SignCompact(a.privateKey, msgHash, true)
	if err != nil {
		return nil, err
	}
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
	if err != nil || ok == false {
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
	return EncodePublicKeyToAddress(publicKey, a.chain.Name)
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
