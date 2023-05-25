package btc

import (
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/tyler-smith/go-bip39"
)

type Account struct {
	privateKey *btcec.PrivateKey
	address    *btcutil.AddressPubKey
	chain      *chaincfg.Params
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
		pub     *btcec.PublicKey
		pubData []byte
		chain   *chaincfg.Params
	)
	wif, err := btcutil.DecodeWIF(prikey)
	if err != nil {
		seed, err := types.HexDecodeString(prikey)
		if err != nil {
			return nil, err
		}
		pri, pub = btcec.PrivKeyFromBytes(seed)
		chain, err = netParamsOf(chainnet)
		if err != nil {
			return nil, err
		}
		pubData = pub.SerializeCompressed()
		chain, err = netParamsOf(chainnet)
		if err != nil {
			return nil, err
		}
	} else {
		pri = wif.PrivKey
		pub = wif.PrivKey.PubKey()
		if wif.IsForNet(&chaincfg.SigNetParams) {
			chain = &chaincfg.SigNetParams
		} else if wif.IsForNet(&chaincfg.MainNetParams) {
			chain = &chaincfg.MainNetParams
		} else {
			return nil, ErrUnsupportedChain
		}
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
func (a *Account) NativeSegwitAddress() (string, error) {
	address, err := btcutil.NewAddressWitnessPubKeyHash(a.address.AddressPubKeyHash().ScriptAddress(), a.chain)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

// NestedSegwitAddress P2SH-P2WPKH just for m/49'/
func (a *Account) NestedSegwitAddress() (string, error) {
	witAddr, err := btcutil.NewAddressWitnessPubKeyHash(a.address.AddressPubKeyHash().ScriptAddress(), a.chain)
	if err != nil {
		return "", err
	}
	witnessProgram, err := txscript.PayToAddrScript(witAddr)
	if err != nil {
		return "", err
	}
	address, err := btcutil.NewAddressScriptHash(witnessProgram, a.chain)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

// TaprootAddress P2TR just for m/86'/
func (a *Account) TaprootAddress() (string, error) {
	tapKey := txscript.ComputeTaprootKeyNoScript(a.address.PubKey())
	address, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(tapKey), a.chain,
	)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}

// LegacyAddress P2PKH just for m/44'/
func (a *Account) LegacyAddress() (string, error) {
	return a.address.AddressPubKeyHash().EncodeAddress(), nil
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

func (a *Account) AddressType() string {
	return "Taproot (P2TR)"
}
func (a *Account) DerivePath() string {
	return "--"
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

// @return default is the mainnet address
func (a *Account) Address() string {
	return a.address.EncodeAddress()
}

// TODO: function not implement yet.
func (a *Account) Sign(message []byte, password string) ([]byte, error) {
	return nil, errors.New("TODO: function not implement yet.")
}

// TODO: function not implement yet.
func (a *Account) SignHex(messageHex string, password string) (*base.OptionalString, error) {
	return nil, errors.New("TODO: function not implement yet.")
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
