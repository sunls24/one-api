package cache

import "sync"

type Cache[T any] struct {
	data   map[string]T
	create func(string) T
	m      sync.Mutex
}

func NewCache[T any](create func(string) T) *Cache[T] {
	return &Cache[T]{create: create, data: make(map[string]T)}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	if v, ok := c.data[key]; ok {
		return v, ok
	}
	var def T
	if c.create != nil {
		def = c.create(key)
	}
	c.data[key] = def
	return def, false
}

func (c *Cache[T]) Set(key string, value T) {
	c.m.Lock()
	defer c.m.Unlock()
	c.data[key] = value
}
