package base

type Balance struct {
	Total  string
	Usable string
}

func EmptyBalance() *Balance {
	return &Balance{
		Total:  "0",
		Usable: "0",
	}
}
