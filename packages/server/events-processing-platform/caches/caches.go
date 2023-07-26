package caches

import (
	"github.com/coocood/freecache"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	"strings"
)

const (
	KB         = 1024
	cache10KB  = 10 * KB
	cache500KB = 500 * KB
)
const (
	expire30Days = 30 * 24 * 60 * 60 // 30 days
)

type Cache interface {
	SetIndustry(key, value string)
	GetIndustry(key string) (string, bool)
}

type cache struct {
	permanentIndustries map[string]string
	industryCache       *freecache.Cache

	marketCache *freecache.Cache
}

func InitCaches() Cache {
	result := cache{
		industryCache:       freecache.NewCache(cache10KB),
		permanentIndustries: data.IndustryValuesUpperCaseMap(),

		marketCache: freecache.NewCache(cache500KB),
	}

	return &result
}

// Industry cache
func (c *cache) SetIndustry(key, value string) {
	// Convert strings to []byte
	keyBytes := []byte(strings.ToUpper(key))
	valueBytes := []byte(value)

	_ = c.industryCache.Set(keyBytes, valueBytes, expire30Days)
}

func (c *cache) GetIndustry(key string) (string, bool) {
	upperKey := strings.ToUpper(key)
	if val, ok := c.permanentIndustries[key]; ok {
		return val, true
	}
	return c.get(c.industryCache, upperKey)
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
