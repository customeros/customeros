package caches

import (
	"github.com/coocood/freecache"
	"sync"
)

const (
	KB       = 1024
	cache5MB = 5 * 1024 * KB
)
const (
	days30       = 30
	expire30Days = days30 * 24 * 60 * 60 // 20 days
)
const delimiter = "--"

type Cache struct {
	mu    sync.RWMutex
	cache *freecache.Cache
}

func NewCache() *Cache {
	return &Cache{
		cache: freecache.NewCache(cache5MB),
	}
}

func (c *Cache) AddExternalSystem(tenant, externalSystem string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert strings to []byte
	keyBytes := []byte(tenant + delimiter + externalSystem)
	valueBytes := []byte("1") // A simple marker value

	_ = c.cache.Set(keyBytes, valueBytes, expire30Days)
}

func (c *Cache) CheckExternalSystem(tenant, externalSystem string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Convert strings to []byte
	keyBytes := []byte(tenant + delimiter + externalSystem)

	_, err := c.cache.Get(keyBytes)
	if err != nil {
		return false // Key not found in cache
	}

	return true // Key found in cache
}
