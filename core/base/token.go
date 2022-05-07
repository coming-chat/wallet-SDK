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
}
