package caches

import (
	"github.com/coocood/freecache"
	"sync"
)

const (
	KB       = 1024
	cache5MB = 5 * 1024 * KB
	cache1MB = 1 * 1024 * KB
)
const (
	expire30Days = 30 * 24 * 60 * 60
	expire1Day   = 24 * 60 * 60
)
const delimiter = "--"

type Cache struct {
	mu                  sync.RWMutex
	externalSystemCache *freecache.Cache
	tenantCache         *freecache.Cache
}

func NewCache() *Cache {
	return &Cache{
		externalSystemCache: freecache.NewCache(cache5MB),
		tenantCache:         freecache.NewCache(cache1MB),
	}
}

func (c *Cache) AddExternalSystem(tenant, externalSystem string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert strings to []byte
	keyBytes := []byte(tenant + delimiter + externalSystem)
	valueBytes := []byte("1") // A simple marker value

	_ = c.externalSystemCache.Set(keyBytes, valueBytes, expire30Days)
}

func (c *Cache) CheckExternalSystem(tenant, externalSystem string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Convert strings to []byte
	keyBytes := []byte(tenant + delimiter + externalSystem)

	_, err := c.externalSystemCache.Get(keyBytes)
	if err != nil {
		return false // Key not found in cache
	}

	return true // Key found in cache
}

func (c *Cache) AddTenant(tenant string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert strings to []byte
	keyBytes := []byte(tenant)
	valueBytes := []byte("1") // A simple marker value

	_ = c.tenantCache.Set(keyBytes, valueBytes, expire1Day)
}

func (c *Cache) CheckTenant(tenant string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Convert strings to []byte
	keyBytes := []byte(tenant)

	_, err := c.tenantCache.Get(keyBytes)
	if err != nil {
		return false // Key not found in cache
	}

	return true // Key found in cache
}
