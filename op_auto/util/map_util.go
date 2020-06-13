package util

import (
	"sync"
)

//map 并发存取
type BeeMap struct {
	lock *sync.RWMutex
	bm   map[string]interface{}
}

func NewBeeMap() *BeeMap {
	return &BeeMap{
		lock: new(sync.RWMutex),
		bm:   make(map[string]interface{}),
	}
}

//Get from maps return the k's value
func (m *BeeMap) Get(k string) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

// Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *BeeMap) Set(k string, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

// Returns true if k is exist in the map.
func (m *BeeMap) Check(k string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.bm[k]; !ok {
		return false
	}
	return true
}

func (m *BeeMap) Delete(k string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
}
