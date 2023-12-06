package base

import (
	"encoding/json"
	"math/big"

	"github.com/coming-chat/wallet-SDK/core/base/inter"
)

// exchange `Aniable Object` & `Any`
type Aniable interface {
	AsAny() *Any

	// You need to implement the following methods if your class name is Xxx
	// func AsXxx(a *Any) *Xxx

	// ======= template
	// func (a *Xxx) AsAny() *base.Any {
	// 	return &base.Any{Value: a}
	// }
	//	func AsXxx(a *base.Any) *Xxx {
	//		if r, ok := a.Value.(*Xxx); ok {
	//			return r
	//		}
	//		if r, ok := a.Value.(Xxx); ok {
	//			return &r
	//		}
	//		return nil
	//	}
}

// 如果需要自定义类型支持 Any, 需要遵循协议 Aniable
type Any struct {
	Value any
}

func NewAny() *Any {
	return &Any{}
}

// `Any` only support Marshal
func (a Any) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Value)
}

func (a *Any) JsonString() (*OptionalString, error) {
	return JsonString(a)
}

func (a *Any) SetString(v string)  { a.Value = v }
func (a *Any) SetBool(v bool)      { a.Value = v }
func (a *Any) SetInt(v int)        { a.Value = v }
func (a *Any) SetInt8(v int8)      { a.Value = v }
func (a *Any) SetInt16(v int16)    { a.Value = v }
func (a *Any) SetInt32(v int32)    { a.Value = v }
func (a *Any) SetInt64(v int64)    { a.Value = v }
func (a *Any) SetUInt8(v *BigInt)  { n := uint8(v.bigint.Uint64()); a.Value = &n }
func (a *Any) SetUInt16(v *BigInt) { n := uint16(v.bigint.Uint64()); a.Value = &n }
func (a *Any) SetUInt32(v *BigInt) { n := uint32(v.bigint.Uint64()); a.Value = &n }
func (a *Any) SetUInt64(v *BigInt) { n := v.bigint.Uint64(); a.Value = &n }
func (a *Any) SetBigInt(v *BigInt) { a.Value = v }

func (a *Any) GetString() string  { return a.Value.(string) }
func (a *Any) GetBool() bool      { return a.Value.(bool) }
func (a *Any) GetInt() int        { return a.Value.(int) }
func (a *Any) GetInt8() int8      { return a.Value.(int8) }
func (a *Any) GetInt16() int16    { return a.Value.(int16) }
func (a *Any) GetInt32() int32    { return a.Value.(int32) }
func (a *Any) GetInt64() int64    { return a.Value.(int64) }
func (a *Any) GetUInt8() *BigInt  { return &BigInt{new(big.Int).SetUint64(uint64(*a.Value.(*uint8)))} }
func (a *Any) GetUInt16() *BigInt { return &BigInt{new(big.Int).SetUint64(uint64(*a.Value.(*uint16)))} }
func (a *Any) GetUInt32() *BigInt { return &BigInt{new(big.Int).SetUint64(uint64(*a.Value.(*uint32)))} }
func (a *Any) GetUInt64() *BigInt { return &BigInt{new(big.Int).SetUint64(*a.Value.(*uint64))} }
func (a *Any) GetBigInt() *BigInt { return a.Value.(*BigInt) }

// MARK - AnyArray

type AnyArray struct {
	inter.AnyArray[*Any]
}

func NewAnyArray() *AnyArray {
	return &AnyArray{[]*Any{}}
}

// MARK - AnyMap

type AnyMap struct {
	inter.AnyMap[string, *Any]
}

func NewAnyMap() *AnyMap {
	return &AnyMap{map[string]*Any{}}
}

func (a *AnyMap) Keys() *StringArray {
	keys := inter.KeysOf(a.AnyMap)
	return &StringArray{keys}
}
