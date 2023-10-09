package testcase

import "strconv"

type Amount struct {
	Amount   float64
	Multiple float64
}

func (a Amount) Int64() int64 {
	return int64(a.Amount * a.Multiple)
}
func (a Amount) Uint64() uint64 {
	return uint64(a.Amount * a.Multiple)
}
func (a Amount) String() string {
	return strconv.FormatUint(uint64(a.Amount*a.Multiple), 10)
}
