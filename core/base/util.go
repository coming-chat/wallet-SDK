package base

import (
	"errors"
	"math/big"
	"strconv"
	"sync"
)

// This method will traverse the array concurrently and map each object in the array.
// @param list: [TYPE1], a list that all item is TYPE1
// @param limit: maximum number of tasks to execute, 0 means no limit
// @param maper: func(TYPE1) (TYPE2, error), a function that input TYPE1, return TYPE2
//                you can throw an error to finish task.
// @return : [TYPE2], a list that all item is TYPE2
// @example : ```
//     nums := []interface{}{1, 2, 3, 4, 5, 6}
//     res, _ := MapListConcurrent(nums, func(i interface{}) (interface{}, error) {
//         return strconv.Itoa(i.(int) * 100), nil
//     })
//     println(res) // ["100" "200" "300" "400" "500" "600"]
// ```
func MapListConcurrent(list []interface{}, limit int, maper func(interface{}) (interface{}, error)) ([]interface{}, error) {
	thread := 0
	max := limit
	wg := sync.WaitGroup{}

	mapContainer := newSafeMap()
	var firstError error
	for _, item := range list {
		if firstError != nil {
			continue
		}
		if max == 0 {
			wg.Add(1)
			// no limit
		} else {
			if thread == max {
				wg.Wait()
				thread = 0
			}
			if thread < max {
				wg.Add(1)
			}
		}

		go func(w *sync.WaitGroup, item interface{}, mapContainer *safeMap, firstError *error) {
			maped, err := maper(item)
			if *firstError == nil && err != nil {
				*firstError = err
			} else {
				mapContainer.writeMap(item, maped)
			}
			wg.Done()
		}(&wg, item, mapContainer, &firstError)
		thread++
	}
	wg.Wait()
	if firstError != nil {
		return nil, firstError
	}

	result := []interface{}{}
	for _, item := range list {
		result = append(result, mapContainer.Map[item])
	}
	return result, nil
}

// The encapsulation of MapListConcurrent.
func MapListConcurrentStringToString(strList []string, maper func(string) (string, error)) ([]string, error) {
	list := make([]interface{}, len(strList))
	for i, s := range strList {
		list[i] = s
	}
	temp, err := MapListConcurrent(list, 10, func(i interface{}) (interface{}, error) {
		return maper(i.(string))
	})
	if err != nil {
		return nil, err
	}

	result := make([]string, len(temp))
	for i, v := range temp {
		result[i] = v.(string)
	}
	return result, nil
}

// Return the more biger of the two numbers
func MaxBigInt(x, y *big.Int) *big.Int {
	if x.Cmp(y) > 0 {
		return x
	} else {
		return y
	}
}

// @note float64 should use `math.Max()`
func Max[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | string](x, y T) T {
	if x >= y {
		return x
	} else {
		return y
	}
}

// @note float64 should use `math.Min()`
func Min[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | string](x, y T) T {
	if x <= y {
		return x
	} else {
		return y
	}
}

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
	// first we have to recover()
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
