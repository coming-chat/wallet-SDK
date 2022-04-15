package eth

import (
	"errors"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// 将MethodOpts 进行转化，由于端的限制，只能传入string字符
func OptsTobigInt(opts *CallMethodOpts) *CallMethodOptsBigInt {

	GasPrice, _ := new(big.Int).SetString(opts.GasPrice, 10)
	GasLimit, _ := strconv.Atoi(opts.GasLimit)
	MaxPriorityFeePerGas, _ := new(big.Int).SetString(opts.MaxPriorityFeePerGas, 10)
	Value, _ := new(big.Int).SetString(opts.Value, 10)

	return &CallMethodOptsBigInt{
		Nonce:                uint64(opts.Nonce),
		Value:                Value,
		GasPrice:             GasPrice,
		MaxPriorityFeePerGas: MaxPriorityFeePerGas,
		GasLimit:             uint64(GasLimit),
	}

}

// 私钥转地址
func PrivateKeyToAddress(privateKey string) (string, error) {
	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", err
	}
	return crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex(), nil
}

func IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

// 这个方法会并发遍历数组, 并映射数组中的每一个对象
// @param list : [TYPE1], a list that all item is TYPE1
// @param maper : func(TYPE1) (TYPE2, error), a function that input TYPE1, return TYPE2
//                you can throw an error to finish task.
// @return : [TYPE2], a list that all item is TYPE2
// @example : ```
//     nums := []interface{}{1, 2, 3, 4, 5, 6}
//     res, _ := MapListConcurrent(nums, func(i interface{}) (interface{}, error) {
//         return strconv.Itoa(i.(int) * 100), nil
//     })
//     println(res) // ["100" "200" "300" "400" "500" "600"]
// ```
func MapListConcurrent(list []interface{}, maper func(interface{}) (interface{}, error)) ([]interface{}, error) {
	thread := 0
	max := 10
	wg := sync.WaitGroup{}

	mapContainer := newSafeMap()
	var firstError error
	for _, item := range list {
		if firstError != nil {
			continue
		}
		if thread == max {
			wg.Wait()
			thread = 0
		}
		if thread < max {
			wg.Add(1)
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

// 对 MapListConcurrent 的封装, 映射对象不是该方法的重点，并发才是
func MapListConcurrentStringToString(strList []string, maper func(string) (string, error)) ([]string, error) {
	list := make([]interface{}, len(strList))
	for i, s := range strList {
		list[i] = s
	}
	temp, err := MapListConcurrent(list, func(i interface{}) (interface{}, error) {
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

// 将任何其他自定义类型的 error, 转为系统基础的 error 对象
func MapToBasicError(e error) error {
	if e == nil {
		return e
	}
	return errors.New(e.Error())
}
