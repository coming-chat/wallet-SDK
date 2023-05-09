package base

type TokenInfo struct {
	Name    string
	Symbol  string
	Decimal int16
}

type Token interface {
	Chain() Chain

	TokenInfo() (*TokenInfo, error)

	BalanceOfAddress(address string) (*Balance, error)
	BalanceOfPublicKey(publicKey string) (*Balance, error)
	BalanceOfAccount(account Account) (*Balance, error)

	BuildTransfer(sender, receiver, amount string) (txn Transaction, err error)
	// Before invoking this method, it is best to check `CanTransferAll()`
	CanTransferAll() bool
	BuildTransferAll(sender, receiver string) (txn Transaction, err error)
}
