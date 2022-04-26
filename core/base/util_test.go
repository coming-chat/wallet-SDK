package base

import "testing"

func TestCatchPanic(t *testing.T) {
	i, err := dangerousCode()
	t.Log(i, err)

	if err != nil {
		t.Log("err = ", err)
	} else {
		t.Log("suc = ", i)
	}
}

func dangerousCode() (i int, e error) {
	defer CatchPanicAndMapToBasicError(&e)

	// runtime error: invalid memory address or nil pointer dereference
	var a Account
	println("......", a.Address())

	// panic(3432434)

	return 13, e
}
