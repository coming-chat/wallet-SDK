package base

import (
	"errors"
	"strconv"
)

/* [zh] 该方法会捕捉 panic 抛出的值，并转成一个 error 对象通过参数指针返回
 *      注意: 如果想要返回它抓住的 error, 必须使用命名返回值！！
 * [en] This method will catch the value thrown by panic, and turn it into an error object and return it through the parameter pointer
 *		Note: If you want to return the error it caught, you must use a named return value! !
 *  ```
 *  func actionWillThrowError(parameters...) (namedErr error, other...) {
 *      defer CatchPanicAndMapToBasicError(&namedErr)
 *      // action code ...
 *      return namedErr, other...
 *  }
 *  ```
 */
func CatchPanicAndMapToBasicError(errOfResult *error) {
	errOfPanic := recover()
	if errOfResult == nil {
		return
	}
	if errOfPanic != nil {
		*errOfResult = MapAnyToBasicError(errOfPanic)
	} else {
		*errOfResult = MapAnyToBasicError(*errOfResult)
	}
}

func MapAnyToBasicError(e any) error {
	if e == nil {
		return nil
	}

	err, ok := e.(error)
	if ok {
		return errors.New(err.Error())
	}

	msg, ok := e.(string)
	if ok {
		return errors.New("panic error: " + msg)
	}

	code, ok := e.(int)
	if ok {
		return errors.New("panic error: code = " + strconv.Itoa(code))
	}

	return errors.New("panic error: unexpected error.")
}
