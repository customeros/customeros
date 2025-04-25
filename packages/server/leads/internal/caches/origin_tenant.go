package caches

import (
	"encoding/json"

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

type OriginData struct {
	Tenant       string `json:"tenant"`
	WebTrackerID string `json:"webTrackerId"`
}

// SetDataForOrigin stores both tenant and webTrackerID mapped by origin, expiring after 1 hour.
func (o *OriginTenantCache) SetDataForOrigin(origin, tenant, webTrackerID string) error {
	data := OriginData{
		Tenant:       tenant,
		WebTrackerID: webTrackerID,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return o.cache.Set([]byte(origin), jsonData, expire1Hour)
}

// GetDataForOrigin retrieves both tenant and webTrackerID associated with the given origin key.
// Returns empty strings if not found or if an error occurs.
func (o *OriginTenantCache) GetDataForOrigin(origin string) (tenant string, webTrackerID string, err error) {
	value, err := o.cache.Get([]byte(origin))
	if err != nil {
		return "", "", err
	}

	var data OriginData
	if err := json.Unmarshal(value, &data); err != nil {
		return "", "", err
	}

	return data.Tenant, data.WebTrackerID, nil
}
