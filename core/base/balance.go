package base

import "strconv"

type Balance struct {
	Total  string
	Usable string
}

// Deprecated: use `NewBalance("0")`
func EmptyBalance() *Balance {
	return &Balance{
		Total:  "0",
		Usable: "0",
	}
}

func NewBalance(amount string) *Balance {
	return &Balance{Total: amount, Usable: amount}
}

func NewBalanceWithInt(amount int64) *Balance {
	a := strconv.FormatInt(amount, 10)
	return &Balance{Total: a, Usable: a}
}
