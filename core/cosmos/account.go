package cosmos

import (
	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

type Account struct {
	privKey       types.PrivKey
	Cointype      int64
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

// Is deocde from address
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
func (a *Account) SignHex(messageHex string, password string) ([]byte, error) {
	message, err := hexTypes.HexDecodeString(messageHex)
	if err != nil {
		return nil, err
	}
	return a.privKey.Sign(message)
}
