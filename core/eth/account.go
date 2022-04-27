package eth

import (
	"crypto/ecdsa"
	"errors"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
)

type Account struct {
	*Util
	privateKey []byte
	address    string
}

func NewAccountWithMnemonic(mnemonic string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		return nil, err
	}

	key := masterKey
	for _, n := range path {
		key, err = key.DeriveNonStandard(n)
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	privateKeyECDSA := privateKey.ToECDSA()
	privateKeyData := privateKey.Serialize()

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("General public key failed.")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	addressHex := address.Hex()

	return &Account{
		Util:       NewUtil(),
		privateKey: privateKeyData,
		address:    addressHex,
	}, nil
}

// MARK - Implement the protocol wallet.Account

// @return privateKey data
func (a *Account) PrivateKey() ([]byte, error) {
	return a.privateKey, nil
}

// @return privateKey string that will start with 0x.
func (a *Account) PrivateKeyHex() (string, error) {
	return types.HexEncodeToString(a.privateKey), nil
}

// Is deocde from address
// @return publicKey data
func (a *Account) PublicKey() []byte {
	pub, _ := types.HexDecodeString(a.address)
	return pub
}

// The ethereum public key is same as address in coming
// @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	return a.address
}

// The ethereum address is same as public key in coming
func (a *Account) Address() string {
	return a.address
}

// TODO: function not implement yet.
func (a *Account) Sign(message []byte, password string) ([]byte, error) {
	return nil, errors.New("TODO: function not implement yet.")
}

// TODO: function not implement yet.
func (a *Account) SignHex(messageHex string, password string) ([]byte, error) {
	return nil, errors.New("TODO: function not implement yet.")
}
