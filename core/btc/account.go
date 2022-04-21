package btc

import "errors"

type Account struct {
	*RootAccount
	*Util
	Chainnet string
}

func NewAccountWithMnemonic(mnemonic, chainnet string) (*Account, error) {
	root, err := NewRootAccountWithMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	return NewAccountFromRoot(root, chainnet)
}

func NewAccountFromRoot(root *RootAccount, chainnet string) (*Account, error) {
	if root == nil {
		return nil, errors.New("no root account")
	}
	util, err := NewUtilWithChainnet(chainnet)
	if err != nil {
		return nil, err
	}

	// re-encode address
	account := &Account{root, util, chainnet}
	account.address, err = account.EncodePublicKeyToAddress(account.publicKey)
	if err != nil {
		return nil, err
	}

	return account, nil
}
