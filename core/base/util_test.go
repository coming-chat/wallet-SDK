package base

import (
	"strconv"
	"testing"
)

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

func TestMapConcurrent(t *testing.T) {
	nums := []interface{}{1, 2, 3, 4, 5, 6}
	// nums := []interface{}{"1", "2", "3", "4"}
	res, _ := MapListConcurrent(nums, 10, func(i interface{}) (interface{}, error) {
		return strconv.Itoa(i.(int) * 100), nil
	})
	t.Log(res)
}
