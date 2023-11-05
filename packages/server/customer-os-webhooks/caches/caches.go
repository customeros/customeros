package caches

import (
	"encoding/json"
	"github.com/coocood/freecache"
	"strconv"
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
	expire1Hour  = 60 * 60
)
const delimiter = "--"

type Cache struct {
	mu                         sync.RWMutex
	externalSystemCache        *freecache.Cache
	tenantCache                *freecache.Cache
	perosnalEmailProviderCache *freecache.Cache
}

func NewCache() *Cache {
	return &Cache{
		externalSystemCache:        freecache.NewCache(cache5MB),
		tenantCache:                freecache.NewCache(cache1MB),
		perosnalEmailProviderCache: freecache.NewCache(cache5MB),
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

func (c *Cache) SetPersonalEmailProviders(domains []string) {
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
			c.perosnalEmailProviderCache.Clear() // Clear the cache
			return
		}

		// Generate a key based on the index
		key := strconv.Itoa(i/chunkSize + 1) // Convert the integer to a string

		// Store the chunk in the cache
		err = c.perosnalEmailProviderCache.Set([]byte(key), domainChunkBytes, expire1Hour)
		if err != nil {
			c.perosnalEmailProviderCache.Clear()
		}
	}
}

func (c *Cache) GetPersonalEmailProviders() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var allDomains []string
	keyIndex := 1

	for {
		// Generate the key based on index
		key := strconv.Itoa(keyIndex)

		// Attempt to get the domains chunk from the cache
		domainChunkBytes, err := c.perosnalEmailProviderCache.Get([]byte(key))
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
