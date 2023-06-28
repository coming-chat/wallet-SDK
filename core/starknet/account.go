package starknet

import (
	"crypto/sha256"
	"errors"
	"math/big"
	"strings"

	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/go-aptos/crypto/derivation"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/dontpanicdao/caigo"
	"github.com/tyler-smith/go-bip39"
)

type Account struct {
	privateKey *big.Int
}

func grindKey(seed []byte) (*big.Int, error) {
	// order := caigo.Curve.N
	// max := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil)
	// limit := big.NewInt(0).Sub(max, big.NewInt(0).Mod(max, order))
	limit := caigo.Curve.N

	for i := 0; i < 100000; i++ {
		bb := append(seed, big.NewInt(int64(i)).Bytes()...)
		key := sha256.Sum256(bb)
		kb := big.NewInt(0).SetBytes(key[:])
		if kb.Cmp(limit) == -1 {
			return kb, nil
		}
	}
	return nil, errors.New("grindKey is broken: tried 100k vals")
}

func IsValidPrivateKey(key string) bool {
	_, err := AccountWithPrivateKey(key)
	return err == nil
}

func NewAccountWithMnemonic(mnemonic string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	key, err := derivation.DeriveForPath("m/44'/9004'/0'/0", seed)
	if err != nil {
		return nil, err
	}
	prikey, err := grindKey(key.Key)
	if err != nil {
		return nil, err
	}
	return &Account{
		privateKey: prikey,
	}, nil
}

func AccountWithPrivateKey(privatekey string) (*Account, error) {
	var priKey *big.Int
	if strings.HasPrefix(privatekey, "0x") {
		priData, err := hexTypes.HexDecodeString(privatekey)
		if err != nil {
			return nil, err
		}
		priKey = big.NewInt(0).SetBytes(priData)
	} else {
		var ok bool
		priKey, ok = big.NewInt(0).SetString(privatekey, 10)
		if !ok {
			return nil, base.ErrInvalidPrivateKey
		}
	}

	_, _, err := caigo.Curve.PrivateToPoint(priKey)
	if err != nil {
		return nil, err
	}
	return &Account{
		privateKey: priKey,
	}, nil
}

// MARK - Implement the protocol wallet.Account

// @return privateKey data
func (a *Account) PrivateKey() ([]byte, error) {
	return a.privateKey.Bytes(), nil
}

// @return privateKey string that will start with 0x.
func (a *Account) PrivateKeyHex() (string, error) {
	return hexTypes.HexEncodeToString(a.privateKey.Bytes()), nil
}

// Is deocde from address
// @return publicKey data
func (a *Account) PublicKey() []byte {
	pub, _, err := caigo.Curve.PrivateToPoint(a.privateKey)
	if err != nil {
		return nil
	}
	return pub.Bytes()
}

// The ethereum public key is same as address in coming
// @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	pub, _, err := caigo.Curve.PrivateToPoint(a.privateKey)
	if err != nil {
		return ""
	}
	return hexTypes.HexEncodeToString(pub.Bytes())
}

// The ethereum address is same as public key in coming
func (a *Account) Address() string {
	addr, _ := EncodePublicKeyToAddress(a.PublicKeyHex())
	return addr
}

func (a *Account) Sign(message []byte, password string) ([]byte, error) {
	return nil, base.ErrUnsupportedFunction
}

func (a *Account) SignHex(messageHex string, password string) (*base.OptionalString, error) {
	return nil, base.ErrUnsupportedFunction
}

func AsStarknetAccount(account base.Account) *Account {
	if r, ok := account.(*Account); ok {
		return r
	} else {
		return nil
	}
}

func (a *Account) SignHash(msgHash *big.Int) (*big.Int, *big.Int, error) {
	return caigo.Curve.Sign(msgHash, a.privateKey)
}
