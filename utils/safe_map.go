package utils

import "sync"

type SafeMap struct {
	lock *sync.RWMutex
	sap  map[interface{}]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		sap:  make(map[interface{}]interface{}),
	}
}

func (m *SafeMap) Add(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Lock()

	if val, ok := m.sap[k]; !ok {
		m.sap[k] = v
	} else if val != v {
		m.sap[k] = v
	} else {
		return false
	}
	return true
}

func (m *SafeMap) Get(k interface{}) interface{} {

	m.lock.RLock()

	defer m.lock.RUnlock()

	if val, ok := m.sap[k]; ok {
		return val
	}

	return nil
}

func (m *SafeMap) Check(k interface{}) bool {

	m.lock.RLock()

	defer m.lock.RUnlock()

	_, ok := m.sap[k]

	return ok
}

func (m *SafeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.sap, k)
}

func (m *SafeMap) Items() map[interface{}]interface{} {
	m.lock.RLock()

	defer m.lock.RUnlock()

	r := make(map[interface{}]interface{})

	for k, v := range m.sap {
		r[k] = v
	}

	return r
}

func (m *SafeMap) Count() int {

	m.lock.RLock()

	defer m.lock.RUnlock()

	return len(m.sap)
}
