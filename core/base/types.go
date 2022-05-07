package base

import "sync"

type SDKEnumInt = int
type SDKEnumString = string

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
