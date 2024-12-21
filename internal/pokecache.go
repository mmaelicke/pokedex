package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}
type Cache struct {
	mu   *sync.Mutex
	data map[string]cacheEntry
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.data[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		t := <-ticker.C
		for key, entry := range c.data {
			if t.After(entry.createdAt.Add(interval)) {
				c.mu.Lock()
				delete(c.data, key)
				c.mu.Unlock()
			}
		}
	}
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		mu:   &sync.Mutex{},
		data: map[string]cacheEntry{},
	}
	// start the reaper
	go cache.reapLoop(interval)
	return cache
}
