package base

import (
	"encoding/json"
	"sync"
)

type SDKEnumInt = int
type SDKEnumString = string

// Optional string for easy of writing iOS code
type OptionalString struct {
	Value string
}

// Optional bool for easy of writing iOS code
type OptionalBool struct {
	Value bool
}

type safeMap struct {
	sync.RWMutex
	Map map[interface{}]interface{}
}

func newSafeMap() *safeMap {
	return &safeMap{Map: make(map[interface{}]interface{})}
}

func (l *safeMap) readMap(key interface{}) (interface{}, bool) {
	l.RLock()
	value, ok := l.Map[key]
	l.RUnlock()
	return value, ok
}

func (l *safeMap) writeMap(key interface{}, value interface{}) {
	l.Lock()
	l.Map[key] = value
	l.Unlock()
}

type StringArray struct {
	Values []string
}

func (a *StringArray) Count() int {
	return len(a.Values)
}

func (a *StringArray) Append(value string) {
	a.Values = append(a.Values, value)
}

func (a *StringArray) Remove(index int) {
	a.Values = append(a.Values[:index], a.Values[index+1:]...)
}

func (a *StringArray) SetValue(value string, index int) {
	a.Values[index] = value
}

func (a *StringArray) ValueOf(index int) string {
	return a.Values[index]
}

func (a *StringArray) String() string {
	data, err := json.Marshal(a.Values)
	if err != nil {
		return "[]"
	}
	return string(data)
}
