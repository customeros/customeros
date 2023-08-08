package caches

import (
	"github.com/coocood/freecache"
)

const (
	KB         = 1024
	cache500KB = 1024 * KB
)
const (
	expire20Days = 20 * 24 * 60 * 60 // 20 days
)

type Cache interface {
	SetUser(key, value string)
	GetUser(key string) (string, bool)
}

type cache struct {
	usersCache *freecache.Cache
}

func InitCaches() Cache {
	result := cache{
		usersCache: freecache.NewCache(cache500KB),
	}

	return &result
}

// User cache
func (c *cache) SetUser(key, value string) {
	// Convert strings to []byte
	keyBytes := []byte(key)
	valueBytes := []byte(value)

	_ = c.usersCache.Set(keyBytes, valueBytes, expire20Days)
}

func (c *cache) GetUser(key string) (string, bool) {
	return c.get(c.usersCache, key)
}

func (c *cache) get(cache *freecache.Cache, key string) (string, bool) {
	keyBytes := []byte(key)
	value, err := cache.Get(keyBytes)

	var strValue string
	if err != nil {
		return strValue, false
	}
	strValue = string(value)
	if strValue == "" {
		return strValue, false
	}
	return strValue, true
}
