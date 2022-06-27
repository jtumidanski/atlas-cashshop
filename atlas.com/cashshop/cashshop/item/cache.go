package item

import (
	"errors"
	"sync"
)

type cache struct {
	items map[uint32]Model
	lock  sync.RWMutex
}

var c *cache
var once sync.Once

func GetCache() *cache {
	once.Do(func() {
		c = &cache{
			items: make(map[uint32]Model, 0),
			lock:  sync.RWMutex{},
		}
	})
	return c
}

func (c *cache) Init() error {
	items, err := readItems()
	if err != nil {
		return err
	}

	c.lock.Lock()
	for _, q := range items {
		c.items[q.SerialNumber()] = q
	}
	c.lock.Unlock()
	return nil
}

func (c *cache) GetItem(serialNumber uint32) (Model, error) {
	c.lock.RLock()
	if val, ok := c.items[serialNumber]; ok {
		c.lock.RUnlock()
		return val, nil
	}
	return Model{}, errors.New("item not found")
}
