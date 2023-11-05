package caches

import (
	"encoding/json"
	"github.com/coocood/freecache"
	"sync"
)

const (
	KB         = 1024
	cache100KB = 100 * KB
	cache1MB   = 1 * 1024 * KB
	cache10MB  = 10 * 1024 * KB
)

const (
	expire5Min  = 5 * 60  // 5 minutes
	expire15Min = 15 * 60 // 15 minutes
	expire1Hour = 1 * 60 * 60
)

type UserDetail struct {
	UserId string   `json:"userId"`
	Tenant string   `json:"tenant"`
	Roles  []string `json:"roles"`
}

type Cache struct {
	mu              sync.RWMutex
	apiKeyCache     *freecache.Cache
	tenantCache     *freecache.Cache
	userDetailCache *freecache.Cache
}

func NewCommonCache() *Cache {
	return &Cache{
		apiKeyCache:     freecache.NewCache(cache100KB),
		tenantCache:     freecache.NewCache(cache1MB),
		userDetailCache: freecache.NewCache(cache10MB),
	}
}

// SetApiKey Method to add an API key to the cache
func (c *Cache) SetApiKey(app, apiKey string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	keyBytes := []byte(string(app)) // Use app as the key
	valueBytes := []byte(apiKey)    // Store apiKey as the value

	_ = c.apiKeyCache.Set(keyBytes, valueBytes, expire5Min)
}

// CheckApiKey Method to check if an API key is in the cache
func (c *Cache) CheckApiKey(app, apiKey string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keyBytes := []byte(string(app))
	valueBytes, err := c.apiKeyCache.Get(keyBytes)
	if err != nil {
		return false // Key not found in cache
	}

	return string(valueBytes) == apiKey // Check if the apiKey matches the one in the cache
}

func (c *Cache) AddTenant(tenant string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert strings to []byte
	keyBytes := []byte(tenant)
	valueBytes := []byte("1") // A simple marker value

	_ = c.tenantCache.Set(keyBytes, valueBytes, expire1Hour)
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

func (c *Cache) GetUserDetailsFromCache(username string) (string, string, []string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Convert strings to []byte
	keyBytes := []byte(username)

	valueBytes, err := c.userDetailCache.Get(keyBytes)
	if err != nil {
		return "", "", []string{}, false // Key not found in cache
	}

	var userDetail UserDetail
	err = json.Unmarshal(valueBytes, &userDetail)
	if err != nil {
		return "", "", []string{}, false // Key not found in cache
	}

	return userDetail.UserId, userDetail.Tenant, userDetail.Roles, true // Key found in cache
}

func (c *Cache) AddUserDetailsToCache(username string, userId string, tenant string, roles []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert strings to []byte
	keyBytes := []byte(username)

	userDetail := UserDetail{
		UserId: userId,
		Tenant: tenant,
		Roles:  roles,
	}

	valueBytes, _ := json.Marshal(userDetail)

	_ = c.userDetailCache.Set(keyBytes, valueBytes, expire15Min)
}
