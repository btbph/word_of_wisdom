package data_structures

type ConcurrentSet[K comparable] struct {
	m ConcurrentMap[K, struct{}]
}

func NewConcurrentSet[K comparable]() *ConcurrentSet[K] {
	return &ConcurrentSet[K]{
		m: *NewConcurrentMap[K, struct{}](),
	}
}

func (c *ConcurrentSet[K]) Exist(key K) bool {
	_, ok := c.m.Get(key)
	return ok
}

func (c *ConcurrentSet[K]) Insert(key K) {
	c.m.Insert(key, struct{}{})
}
