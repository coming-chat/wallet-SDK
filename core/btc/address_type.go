package btc

import (
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type AddressType = base.SDKEnumInt

const (
	AddressTypeComingTaproot AddressType = 0
	AddressTypeNativeSegwit  AddressType = 1
	AddressTypeNestedSegwit  AddressType = 2
	AddressTypeTaproot       AddressType = 3
	AddressTypeLegacy        AddressType = 4
)

func AddressTypeDescription(t AddressType) string {
	switch t {
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

func AddressTypeDerivePath(t AddressType) string {
	switch t {
	case AddressTypeComingTaproot:
		return "--"
	case AddressTypeNativeSegwit:
		return "m/84'/0'/0'/0/0"
	case AddressTypeNestedSegwit:
		return "m/49'/0'/0'/0/0"
	case AddressTypeTaproot:
		return "m/86'/0'/0'/0/0"
	case AddressTypeLegacy:
		return "m/44'/0'/0'/0/0"
	}
	return "--"
}

func EncodePubKeyToAddress(pubkey *btcec.PublicKey, chain *chaincfg.Params, addressType AddressType) (string, error) {
	switch addressType {
	case AddressTypeComingTaproot:
		return EncodeAddressComingTaproot(pubkey, chain)
	case AddressTypeNativeSegwit:
		return EncodeAddressNativeSegwit(pubkey, chain)
	case AddressTypeNestedSegwit:
		return EncodeAddressNestedSegwit(pubkey, chain)
	case AddressTypeTaproot:
		return EncodeAddressTaproot(pubkey, chain)
	case AddressTypeLegacy:
		return EncodeAddressLegacy(pubkey, chain)
	default:
		return "", errors.New("unsupported address type")
	}
}

func EncodeAddressComingTaproot(pubkey *btcec.PublicKey, chain *chaincfg.Params) (res string, err error) {
	address, err := btcutil.NewAddressTaproot(pubkey.SerializeCompressed()[1:33], chain)
	if err != nil {
		return
	}
	return address.EncodeAddress(), nil
}

func EncodeAddressNativeSegwit(pubkey *btcec.PublicKey, chain *chaincfg.Params) (res string, err error) {
	addrPubKey, err := btcutil.NewAddressPubKey(pubkey.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return
	}
	address, err := btcutil.NewAddressWitnessPubKeyHash(addrPubKey.AddressPubKeyHash().ScriptAddress(), chain)
	if err != nil {
		return
	}
	return address.EncodeAddress(), nil
}

func EncodeAddressNestedSegwit(pubkey *btcec.PublicKey, chain *chaincfg.Params) (res string, err error) {
	addrPubKey, err := btcutil.NewAddressPubKey(pubkey.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		return
	}
	witAddr, err := btcutil.NewAddressWitnessPubKeyHash(addrPubKey.AddressPubKeyHash().ScriptAddress(), chain)
	if err != nil {
		return
	}
	witProgram, err := txscript.PayToAddrScript(witAddr)
	if err != nil {
		return
	}
	address, err := btcutil.NewAddressScriptHash(witProgram, chain)
	if err != nil {
		return
	}
	return address.EncodeAddress(), nil
}

func EncodeAddressTaproot(pubkey *btcec.PublicKey, chain *chaincfg.Params) (res string, err error) {
	tapKey := txscript.ComputeTaprootKeyNoScript(pubkey)
	address, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(tapKey), chain)
	if err != nil {
		return
	}
	return address.EncodeAddress(), nil
}

func EncodeAddressLegacy(pubkey *btcec.PublicKey, chain *chaincfg.Params) (res string, err error) {
	address, err := btcutil.NewAddressPubKey(pubkey.SerializeCompressed(), chain)
	if err != nil {
		return
	}
	return address.EncodeAddress(), nil
}
