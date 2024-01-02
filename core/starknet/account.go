package starknet

import (
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"github.com/xiang-xx/starknet.go/curve"
	"github.com/xiang-xx/starknet.go/utils"
)

type Account struct {
	privateKey *big.Int
}

func grindKey(keySeed []byte) (*big.Int, error) {
	keyValueLimit := curve.Curve.N
	sha256EcMaxDigest := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil)
	maxAllowedVal := big.NewInt(0).Sub(sha256EcMaxDigest, big.NewInt(0).Mod(sha256EcMaxDigest, keyValueLimit))

	for i := 0; i < 100000; i++ {
		key := hashKeyWithIndex(keySeed, i)
		if key.Cmp(maxAllowedVal) == -1 {
			return big.NewInt(0).Mod(key, keyValueLimit), nil
		}
	}
	return nil, errors.New("grindKey is broken: tried 100k vals")
}

func hashKeyWithIndex(seed []byte, i int) *big.Int {
	var payload []byte
	if i == 0 {
		payload = append(seed, 0)
	} else {
		payload = append(seed, big.NewInt(int64(i)).Bytes()...)
	}
	hash := sha256.Sum256(payload)
	return big.NewInt(0).SetBytes(hash[:])
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
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	path, err := accounts.ParseDerivationPath("m/44'/9004'/0'/0/0")
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
	prikey, err := grindKey(crypto.FromECDSA(privateKey.ToECDSA()))
	if err != nil {
		return nil, err
	}
	return &Account{
		privateKey: prikey,
	}, nil
}

func AccountWithPrivateKey(privatekey string) (*Account, error) {
	priKey, err := base.ParseNumber(privatekey)
	if err != nil {
		return nil, base.ErrInvalidPrivateKey
	}

	_, _, err = curve.Curve.PrivateToPoint(priKey)
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
	pub, _, err := curve.Curve.PrivateToPoint(a.privateKey)
	if err != nil {
		return nil
	}
	return pub.Bytes()
}

// The ethereum public key is same as address in coming
// @return publicKey string that will start with 0x.
func (a *Account) PublicKeyHex() string {
	pub, _, err := curve.Curve.PrivateToPoint(a.privateKey)
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

func (a *Account) SignHash(msgHash *felt.Felt) ([]*felt.Felt, error) {
	if msgHash == nil {
		return nil, base.ErrInvalidTransactionData
	}
	privateKeyFelt := utils.BigIntToFelt(a.privateKey)
	s1, s2, err := curve.Curve.SignFelt(msgHash, privateKeyFelt)
	if err != nil {
		return nil, err
	}
	return []*felt.Felt{s1, s2}, nil
}
