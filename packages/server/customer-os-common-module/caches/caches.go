package caches

import (
	"encoding/json"
	"github.com/coocood/freecache"
	"strconv"
	"strings"
	"sync"
)

const (
	KB        = 1024
	cache1MB  = 1 * 1024 * KB
	cache10MB = 10 * 1024 * KB
	cache20MB = 20 * 1024 * KB
)

const (
	expire15Min    = 15 * 60             // 15 minutes
	expire1Hour    = 60 * 60             // 1 hour
	expire24Hours  = 24 * 60 * 60        // 24 hours
	expire9999Days = 9999 * 24 * 60 * 60 // 9999 days
)

type UserDetail struct {
	UserId string   `json:"userId"`
	Tenant string   `json:"tenant"`
	Roles  []string `json:"roles"`
}

type Cache struct {
	apiKeyCache                                *freecache.Cache
	tenantApiKeyCache                          *freecache.Cache
	tenantCache                                *freecache.Cache
	userDetailCache                            *freecache.Cache
	organizationWebsiteHostingUrlPatternsCache *freecache.Cache
	personalEmailProviderCache                 *freecache.Cache
}

func NewCommonCache() *Cache {
	return &Cache{
		apiKeyCache:       freecache.NewCache(cache1MB),
		tenantApiKeyCache: freecache.NewCache(cache1MB),
		tenantCache:       freecache.NewCache(cache1MB),
		userDetailCache:   freecache.NewCache(cache20MB),
		organizationWebsiteHostingUrlPatternsCache: freecache.NewCache(cache1MB),
		personalEmailProviderCache:                 freecache.NewCache(cache10MB),
	}
}

// SetApiKey adds an API key to the cache
func (c *Cache) SetApiKey(app, apiKey string) {
	keyBytes := []byte(app)
	valueBytes := []byte(apiKey)

	_ = c.apiKeyCache.Set(keyBytes, valueBytes, expire24Hours)
}

// CheckApiKey checks if an API key is in the cache
func (c *Cache) CheckApiKey(app, apiKey string) bool {
	keyBytes := []byte(app)
	valueBytes, err := c.apiKeyCache.Get(keyBytes)
	if err != nil {
		return false
	}
	return string(valueBytes) == apiKey
}

// CheckTenantApiKey checks if a tenant API key exists in the cache
func (c *Cache) CheckTenantApiKey(apiKey string) bool {
	keyBytes := []byte(apiKey)
	valueBytes, err := c.tenantApiKeyCache.Get(keyBytes)
	if err != nil {
		return false
	}
	return string(valueBytes) == apiKey
}

// SetTenantApiKey sets the tenant's API key in the cache
func (c *Cache) SetTenantApiKey(tenant, apiKey string) {
	keyBytes := []byte(apiKey)
	valueBytes := []byte(tenant)

	_ = c.tenantApiKeyCache.Set(keyBytes, valueBytes, expire24Hours)
}

// AddTenant adds a tenant to the cache
func (c *Cache) AddTenant(tenant string) {
	keyBytes := []byte(tenant)
	valueBytes := []byte("1")

	_ = c.tenantCache.Set(keyBytes, valueBytes, expire1Hour)
}

// CheckTenant verifies if a tenant exists in the cache
func (c *Cache) CheckTenant(tenant string) bool {
	keyBytes := []byte(tenant)
	_, err := c.tenantCache.Get(keyBytes)
	return err == nil
}

// GetUserDetailsFromCache retrieves user details from the cache
func (c *Cache) GetUserDetailsFromCache(username string) (string, string, []string, bool) {
	keyBytes := []byte(username)

	valueBytes, err := c.userDetailCache.Get(keyBytes)
	if err != nil {
		return "", "", []string{}, false
	}

	var userDetail UserDetail
	err = json.Unmarshal(valueBytes, &userDetail)
	if err != nil {
		return "", "", []string{}, false
	}

	return userDetail.UserId, userDetail.Tenant, userDetail.Roles, true
}

// AddUserDetailsToCache stores user details in the cache
func (c *Cache) AddUserDetailsToCache(username, userId, tenant string, roles []string) {
	keyBytes := []byte(username)

	userDetail := UserDetail{
		UserId: userId,
		Tenant: tenant,
		Roles:  roles,
	}
	valueBytes, _ := json.Marshal(userDetail)

	_ = c.userDetailCache.Set(keyBytes, valueBytes, expire15Min)
}

var cachedPersonalEmailProviders []string
var personalEmailProvidersMu sync.RWMutex

// GetPersonalEmailProviders retrieves personal email providers from the cache with local caching
func (c *Cache) GetPersonalEmailProviders() []string {
	personalEmailProvidersMu.RLock()
	if cachedPersonalEmailProviders != nil {
		providers := cachedPersonalEmailProviders
		personalEmailProvidersMu.RUnlock()
		return providers
	}
	personalEmailProvidersMu.RUnlock()

	personalEmailProvidersMu.Lock()
	defer personalEmailProvidersMu.Unlock()

	var allDomains []string
	keyIndex := 1

	for {
		key := strconv.Itoa(keyIndex)
		domainChunkBytes, err := c.personalEmailProviderCache.Get([]byte(key))
		if err != nil {
			break
		}

		var domainChunk []string
		if err = json.Unmarshal(domainChunkBytes, &domainChunk); err != nil {
			break
		}

		allDomains = append(allDomains, domainChunk...)
		keyIndex++
	}

	cachedPersonalEmailProviders = allDomains
	return cachedPersonalEmailProviders
}

// SetPersonalEmailProviders caches the list of personal email providers in chunks
func (c *Cache) SetPersonalEmailProviders(domains []string) {
	const chunkSize = 100
	for i, j := 0, chunkSize; i < len(domains); i, j = i+chunkSize, j+chunkSize {
		if j > len(domains) {
			j = len(domains)
		}

		domainChunk := domains[i:j]
		domainChunkBytes, err := json.Marshal(domainChunk)
		if err != nil {
			c.personalEmailProviderCache.Clear()
			return
		}

		key := strconv.Itoa(i/chunkSize + 1)
		err = c.personalEmailProviderCache.Set([]byte(key), domainChunkBytes, expire9999Days)
		if err != nil {
			c.personalEmailProviderCache.Clear()
		}
	}
}

// IsPersonalEmailProvider checks if a given domain is a personal email provider
func (c *Cache) IsPersonalEmailProvider(domain string) bool {
	domainLower := strings.ToLower(domain)
	for _, v := range c.GetPersonalEmailProviders() {
		if domainLower == strings.ToLower(v) {
			return true
		}
	}
	return false
}

var cachedOrganizationUrlPatterns []string
var organizationUrlPatternsMu sync.RWMutex

// GetOrganizationWebsiteHostingUrlPatters retrieves organization URL patterns from the cache with local caching
func (c *Cache) GetOrganizationWebsiteHostingUrlPatters() []string {
	organizationUrlPatternsMu.RLock()
	if cachedOrganizationUrlPatterns != nil {
		patterns := cachedOrganizationUrlPatterns
		organizationUrlPatternsMu.RUnlock()
		return patterns
	}
	organizationUrlPatternsMu.RUnlock()

	organizationUrlPatternsMu.Lock()
	defer organizationUrlPatternsMu.Unlock()

	var urlPatterns []string
	keyIndex := 1

	for {
		key := strconv.Itoa(keyIndex)
		chunkBytes, err := c.organizationWebsiteHostingUrlPatternsCache.Get([]byte(key))
		if err != nil {
			break
		}

		var chunk []string
		if err = json.Unmarshal(chunkBytes, &chunk); err != nil {
			break
		}

		urlPatterns = append(urlPatterns, chunk...)
		keyIndex++
	}

	cachedOrganizationUrlPatterns = urlPatterns
	return cachedOrganizationUrlPatterns
}

// SetOrganizationWebsiteHostingUrlPatters caches URL patterns in chunks
func (c *Cache) SetOrganizationWebsiteHostingUrlPatters(urlPatterns []string) {
	const chunkSize = 100
	for i, j := 0, chunkSize; i < len(urlPatterns); i, j = i+chunkSize, j+chunkSize {
		if j > len(urlPatterns) {
			j = len(urlPatterns)
		}

		chunk := urlPatterns[i:j]
		chunkBytes, err := json.Marshal(chunk)
		if err != nil {
			c.organizationWebsiteHostingUrlPatternsCache.Clear()
			return
		}

		key := strconv.Itoa(i/chunkSize + 1)
		err = c.organizationWebsiteHostingUrlPatternsCache.Set([]byte(key), chunkBytes, expire1Hour)
		if err != nil {
			c.organizationWebsiteHostingUrlPatternsCache.Clear()
		}
	}
}
