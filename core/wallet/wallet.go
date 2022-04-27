package wallet

import (
	"github.com/ChainSafe/go-schnorrkel"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/btc"
	"github.com/coming-chat/wallet-SDK/core/eth"
	"github.com/coming-chat/wallet-SDK/core/polka"
)

type Wallet struct {
	Mnemonic string

	Keystore string
	password string

	// cache
	polkaAccounts   map[int]*polka.Account
	bitcoinAccounts map[string]*btc.Account
	ethereumAccount *eth.Account
}

func NewWalletFromMnemonic(mnemonic string) (*Wallet, error) {
	if !IsValidMnemonic(mnemonic) {
		return nil, ErrInvalidMnemonic
	}
	return &Wallet{Mnemonic: mnemonic}, nil
}

// Deprecated: NewWallet is deprecated. Please Use NewWalletFromMnemonic instead.
func NewWallet(seedOrPhrase string) (*Wallet, error) {
	return NewWalletFromMnemonic(seedOrPhrase)
}

// Only support Polka keystore.
func NewWalletFromKeyStore(keyStoreJson string, password string) (*Wallet, error) {
	// check is valid keystore
	if !polka.IsValidKeystore(keyStoreJson, password) {
		return nil, ErrKeystore
	}
	return &Wallet{
		Keystore: keyStoreJson,
		password: password,
	}, nil
}

// Get or create the polka account with specified network.
func (w *Wallet) GetOrCreatePolkaAccount(network int) (*polka.Account, error) {
	if w.polkaAccounts == nil {
		w.polkaAccounts = make(map[int]*polka.Account)
	}

	cache := w.polkaAccounts[network]
	if cache != nil {
		return cache, nil
	}

	var account *polka.Account
	var err error
	if len(w.Mnemonic) > 0 {
		account, err = polka.NewAccountWithMnemonic(w.Mnemonic, network)
	} else if len(w.Keystore) > 0 {
		account, err = polka.NewAccountWithKeystore(w.Keystore, w.password, network)
	}
	if err != nil {
		return nil, err
	}
	// save to cache
	w.polkaAccounts[network] = account
	return account, err
}

// Get or create the bitcoin account with specified chainnet.
func (w *Wallet) GetOrCreateBitcoinAccount(chainnet string) (*btc.Account, error) {
	if w.bitcoinAccounts == nil {
		w.bitcoinAccounts = make(map[string]*btc.Account)
	}

	cache := w.bitcoinAccounts[chainnet]
	if cache != nil {
		return cache, nil
	}

	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := btc.NewAccountWithMnemonic(w.Mnemonic, chainnet)
	if err != nil {
		return nil, err
	}

	// save to cache
	w.bitcoinAccounts[chainnet] = account
	return account, err
}

// Get or create the ethereum account.
func (w *Wallet) GetOrCreateEthereumAccount() (*eth.Account, error) {
	cache := w.ethereumAccount
	if cache != nil {
		return cache, nil
	}
	if len(w.Mnemonic) <= 0 {
		return nil, ErrInvalidMnemonic
	}

	account, err := eth.NewAccountWithMnemonic(w.Mnemonic)
	if err != nil {
		return nil, err
	}
	// save to cache
	w.ethereumAccount = account
	return account, err
}

// Deprecated: CheckPassword is deprecated. Please use wallet.PolkaAccount(network).CheckPassword() instead
func (w *Wallet) CheckPassword(password string) (bool, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return false, err
	}
	err = account.CheckPassword(password)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Deprecated: Sign is deprecated. Please use wallet.PolkaAccount(network).Sign() instead
func (w *Wallet) Sign(message []byte, password string) (b []byte, err error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return nil, err
	}
	return account.Sign(message, password)
}

// Deprecated: SignFromHex is deprecated. Please use wallet.PolkaAccount(network).SignHex() instead
func (w *Wallet) SignFromHex(messageHex string, password string) ([]byte, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return nil, err
	}
	return account.SignHex(messageHex, password)
}

// Deprecated: GetPublicKey is deprecated. Please use wallet.PolkaAccount(network).PublicKey() instead
func (w *Wallet) GetPublicKey() ([]byte, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return nil, err
	}
	return types.HexDecodeString(account.PublicKey())
}

// Deprecated: GetPublicKeyHex is deprecated. Please use wallet.PolkaAccount(network).PublicKey() instead
func (w *Wallet) GetPublicKeyHex() (string, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return "", err
	}
	return account.PublicKey(), nil
}

// Deprecated: GetAddress is deprecated. Please use wallet.PolkaAccount(network).Address() instead
func (w *Wallet) GetAddress(network int) (string, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return "", err
	}
	return account.Address(), nil
}

// Deprecated: GetPrivateKeyHex is deprecated. Please use wallet.PolkaAccount(network).PrivateKey() instead
func (w *Wallet) GetPrivateKeyHex() (string, error) {
	account, err := w.GetOrCreatePolkaAccount(44)
	if err != nil {
		return "", err
	}
	return account.PrivateKey()
}

// 内置账号，主要用来给用户未签名的交易签一下名
// 然后给用户去链上查询手续费，保护用户资产安全
func mockWallet() *Wallet {
	mnemonic := "infant carbon above canyon corn collect finger drip area feature mule autumn"
	w, _ := NewWallet(mnemonic)
	return w
}

func Verify(publicKey [32]byte, msg []byte, signature []byte) bool {
	var sigs [64]byte
	copy(sigs[:], signature)
	sig := new(schnorrkel.Signature)
	if err := sig.Decode(sigs); err != nil {
		return false
	}
	publicKeyD := schnorrkel.NewPublicKey(publicKey)
	return publicKeyD.Verify(sig, schnorrkel.NewSigningContext([]byte("substrate"), msg))
}
