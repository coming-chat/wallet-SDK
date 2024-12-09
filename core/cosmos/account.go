package cosmos

import (
	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

type Account struct {
	Cointype      int64
	privKey       types.PrivKey
	AddressPrefix string
}

func NewAccountWithMnemonic(mnemonic string, cointype int64, addressPrefix string) (*Account, error) {
	password := ""
	path := hd.CreateHDPath(uint32(cointype), 0, 0).String()
	derivedPriv, err := hd.Secp256k1.Derive()(mnemonic, password, path)
	if err != nil {
		return nil, err
	}

	privKey := hd.Secp256k1.Generate()(derivedPriv)

	return &Account{
		privKey:       privKey,
		Cointype:      cointype,
		AddressPrefix: addressPrefix,
	}, nil
}

func AccountWithPrivateKey(privatekey string, cointype int64, addressPrefix string) (*Account, error) {
	priData, err := hexTypes.HexDecodeString(privatekey)
	if err != nil {
		return nil, err
	}
	privKey := hd.Secp256k1.Generate()(priData)
	return &Account{
		privKey:       privKey,
		Cointype:      cointype,
		AddressPrefix: addressPrefix,
	}, nil
}

// return NewAccountWithMnemonic(mnemonic, 118, "cosmos")
func NewCosmosAccountWithMnemonic(mnemonic string) (*Account, error) {
	return NewAccountWithMnemonic(mnemonic, sdk.CoinType, sdk.Bech32MainPrefix)
}

// return NewAccountWithMnemonic(mnemonic, 330, "terra")
func NewTerraAccountWithMnemonic(mnemonic string) (*Account, error) {
	return NewAccountWithMnemonic(mnemonic, 330, "terra")
}

// MARK - Implement the protocol wallet.Account

// @return privateKey data
func (a *Account) PrivateKey() ([]byte, error) {
	return a.privKey.Bytes(), nil
}

// @return privateKey string that will start with 0x.
func (a *Account) PrivateKeyHex() (string, error) {
	return hexTypes.HexEncodeToString(a.privKey.Bytes()), nil
}

// Is decode from address
// @return publicKey data
func (a *Account) PublicKey() []byte {
	return a.privKey.PubKey().Bytes()
}

// The ethereum public key is same as address in coming
// @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	return hexTypes.HexEncodeToString(a.privKey.PubKey().Bytes())
}

// The ethereum address is same as public key in coming
func (a *Account) Address() string {
	addr, _ := bech32.ConvertAndEncode(a.AddressPrefix, a.privKey.PubKey().Address().Bytes())
	return addr
}

// TODO: Need Test.
func (a *Account) Sign(message []byte, password string) ([]byte, error) {
	return a.privKey.Sign(message)
}

// TODO: Need Test.
func (a *Account) SignHex(messageHex string, password string) (*base.OptionalString, error) {
	bytes, err := hexTypes.HexDecodeString(messageHex)
	if err != nil {
		return nil, err
	}
	signed, err := a.privKey.Sign(bytes)
	if err != nil {
		return nil, err
	}
	signedString := hexTypes.HexEncodeToString(signed)
	return &base.OptionalString{Value: signedString}, nil
}

func AsCosmosAccount(account base.Account) *Account {
	if r, ok := account.(*Account); ok {
		return r
	} else {
		return nil
	}
}
