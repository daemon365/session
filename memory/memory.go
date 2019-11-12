package memory

import (
	"session/errors"
	"sync"
)

type MemorySession struct {
	data   map[string]interface{}
	id     string
	rwlock sync.RWMutex
}

func (m *MemorySession) Set(key string, values interface{}) (err error) {
	m.rwlock.Lock()
	defer m.rwlock.Unlock()

	m.data[key] = values
	return
}

func (m *MemorySession) Get(key string) (value interface{}, err error) {
	m.rwlock.RLock()
	defer m.rwlock.RUnlock()

	value, ok := m.data[key]
	if !ok {
		err = errors.ErrKeyNotInSession
		return
	}
	return
}

func (m *MemorySession) Del(key string) (err error) {
	m.rwlock.Lock()
	m.rwlock.Unlock()

	delete(m.data, key)
	return
}

func (m *MemorySession) Save() (err error) {
	return
}