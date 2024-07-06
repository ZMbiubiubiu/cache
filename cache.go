// 实现并发控制的in-memory cache

package cache

import (
	"log"
	"sync"

	"cache/lru"
)

type Cache struct {
	mu       sync.Mutex
	lruCache *lru.Cache
}

func NewCache(maxBytes int64) *Cache {
	return &Cache{
		lruCache: lru.NewCache(maxBytes, nil),
	}
}

func (c *Cache) get(key string) (Byteview, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.lruCache.Get(key)
	if !ok {
		return Byteview{}, false
	}
	return v.(Byteview), true
}

func (c *Cache) set(key string, value Byteview) {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.lruCache.Set(key, value)
	if err != nil {
		log.Fatalf("cache set failed || key=%s, err=%v", key, err)
	}
}
