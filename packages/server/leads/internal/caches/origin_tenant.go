package caches

import (
	"github.com/coocood/freecache"
)

const (
	oneMB       = 1 * 1024 * 1024
	expire1Hour = 60 * 60 // 1 hour in seconds
)

// OriginTenantCache wraps a freecache.Cache to store origin → tenant mappings.
type OriginTenantCache struct {
	cache *freecache.Cache
}

// NewOriginTenantCache creates a new OriginTenantCache with 1 MB of cache space.
func NewOriginTenantCache() *OriginTenantCache {
	return &OriginTenantCache{
		cache: freecache.NewCache(oneMB),
	}
}

// SetTenantForOrigin stores tenant (value) mapped by origin (key), expiring after 1 hour.
func (o *OriginTenantCache) SetTenantForOrigin(origin, tenant string) {
	_ = o.cache.Set([]byte(origin), []byte(tenant), expire1Hour)
}

// GetTenantForOrigin retrieves the tenant associated with the given origin key.
// Returns an empty string if not found or if an error occurs.
func (o *OriginTenantCache) GetTenantForOrigin(origin string) string {
	value, err := o.cache.Get([]byte(origin))
	if err != nil {
		return ""
	}
	return string(value)
}
