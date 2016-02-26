package router

import (
	"fmt"
	"sync"
)

//mapCache implements the Cacher interface
// caches the last size byte slices with a
// string key
type mapCache struct {
	sync.RWMutex
	elements map[string][]byte
	keys     []string
	iHead    int
	size     int
}

func NewMapCache(size int) (*mapCache, error) {
	return &mapCache{
		elements: make(map[string][]byte),
		keys:     make([]string, size, size),
		iHead:    0,
		size:     size,
	}, nil
}

func (c *mapCache) Get(key string) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()
	if val, ok := c.elements[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("No key '%s' cached", key)
}

func (c *mapCache) Add(key string, val []byte) error {
	c.Lock()
	defer c.Unlock()

	//remove old entries
	var oldKey string
	oldKey, c.keys[c.iHead] = c.keys[c.iHead], key
	delete(c.elements, oldKey)

	//move head
	c.iHead = (c.iHead + 1) % c.size

	//cache value
	c.elements[key] = val
	return nil
}
