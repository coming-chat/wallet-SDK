package polka

import "errors"

type Account struct {
	*RootAccount
	*Util
	Network int
}

func NewAccountWithMnemonic(mnemonic string, network int) (*Account, error) {
	root, err := NewRootAccountWithMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}
	return NewAccountFromRoot(root, network)
}

func NewAccountWithKeystore(keystoreString, password string, network int) (*Account, error) {
	root, err := NewRootAccountWithKeystore(keystoreString, password)
	if err != nil {
		return nil, err
	}
	return NewAccountFromRoot(root, network)
}

func NewAccountFromRoot(root *RootAccount, network int) (*Account, error) {
	if root == nil {
		return nil, errors.New("no root account")
	}
	util := NewUtilWithNetwork(network)
	account := &Account{root, util, network}
	account.rootUtil = util
	return account, nil
}
