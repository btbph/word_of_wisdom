package data_structures

import "sync"

type ConcurrentMap[K comparable, V any] struct {
	data  map[K]V
	mutex sync.RWMutex
}

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return &ConcurrentMap[K, V]{
		data:  map[K]V{},
		mutex: sync.RWMutex{},
	}
}

func (c *ConcurrentMap[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func (c *ConcurrentMap[K, V]) Insert(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

func (c *ConcurrentMap[K, V]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, key)
}
