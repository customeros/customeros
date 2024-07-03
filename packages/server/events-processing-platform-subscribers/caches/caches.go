package caches

import (
	"encoding/json"
	"github.com/coocood/freecache"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"strconv"
	"strings"
	"sync"
)

const (
	KB         = 1024
	cache10KB  = 10 * KB
	cache500KB = 500 * KB
	cache5MB   = 5 * 1024 * KB
)
const (
	expire20Days = 20 * 24 * 60 * 60 // 20 days
	expire1Hour  = 60 * 60           // 1 hour
)

type Cache interface {
	SetIndustry(key, value string)
	GetIndustry(key string) (string, bool)
	SetPersonalEmailProviders(domains []string)
	GetPersonalEmailProviders() []string
}

type cache struct {
	mu                         sync.RWMutex
	permanentIndustries        map[string]string
	industryCache              *freecache.Cache
	personalEmailProviderCache *freecache.Cache
	marketCache                *freecache.Cache
}

func InitCaches() Cache {
	result := cache{
		industryCache:              freecache.NewCache(cache5MB),
		marketCache:                freecache.NewCache(cache10KB),
		personalEmailProviderCache: freecache.NewCache(cache5MB),
	}
	result.permanentIndustries = data.IndustryValuesUpperCaseMap()
	// add brandfetch industries
	result.permanentIndustries = utils.MergeMaps(result.permanentIndustries, data.BrandfetchIndustryUpperCasedMap())
	// add other industries
	result.permanentIndustries = utils.MergeMaps(result.permanentIndustries, data.OtherIndustryUpperCasedMap())

	return &result
}

// Industry cache
func (c *cache) SetIndustry(key, value string) {
	// Convert strings to []byte
	keyBytes := []byte(strings.ToUpper(key))
	valueBytes := []byte(value)

	_ = c.industryCache.Set(keyBytes, valueBytes, expire20Days)
}

func (c *cache) GetIndustry(key string) (string, bool) {
	upperKey := strings.ToUpper(key)
	if val, ok := c.permanentIndustries[upperKey]; ok {
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

func (c *cache) SetPersonalEmailProviders(domains []string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	const chunkSize = 100 // Size of each domain chunk
	for i, j := 0, chunkSize; i < len(domains); i, j = i+chunkSize, j+chunkSize {
		// This ensures we don't go past the end of the slice
		if j > len(domains) {
			j = len(domains)
		}

		// Get the current chunk and marshal it
		domainChunk := domains[i:j]
		domainChunkBytes, err := json.Marshal(domainChunk)
		if err != nil {
			c.personalEmailProviderCache.Clear() // Clear the cache
			return
		}

		// Generate a key based on the index
		key := strconv.Itoa(i/chunkSize + 1) // Convert the integer to a string

		// Store the chunk in the cache
		err = c.personalEmailProviderCache.Set([]byte(key), domainChunkBytes, expire1Hour)
		if err != nil {
			c.personalEmailProviderCache.Clear()
		}
	}
}

func (c *cache) GetPersonalEmailProviders() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var allDomains []string
	keyIndex := 1

	for {
		// Generate the key based on index
		key := strconv.Itoa(keyIndex)

		// Attempt to get the domains chunk from the cache
		domainChunkBytes, err := c.personalEmailProviderCache.Get([]byte(key))
		if err != nil {
			break // If a key is not found, assume no more chunks are available
		}

		var domainChunk []string
		err = json.Unmarshal(domainChunkBytes, &domainChunk)
		if err != nil {
			// If there is an error unmarshalling, decide how to handle it
			// For simplicity, we stop and return what we have so far
			break
		}

		// Append this chunk of domains to the allDomains slice
		allDomains = append(allDomains, domainChunk...)

		keyIndex++ // Increment key index for next iteration
	}

	return allDomains // Return the combined list of all domains
}
