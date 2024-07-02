package caches

import (
	"encoding/json"
	"github.com/coocood/freecache"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"strconv"
	"sync"
)

const (
	KB         = 1024
	cache500MB = 500 * 1024 * KB
)
const (
	expire9999Days = 9999 * 24 * 60 * 60
	expire30Days   = 30 * 24 * 60 * 60
	expire1Day     = 24 * 60 * 60
	expire1Hour    = 60 * 60
)
const delimiter = "--"

type Cache struct {
	mu                         sync.RWMutex
	organizationsByTenantCache map[string]*freecache.Cache
}

func NewCache() *Cache {
	cache := Cache{
		organizationsByTenantCache: make(map[string]*freecache.Cache),
	}
	return &cache
}

func (c *Cache) SetOrganizations(tenant string, organizations []*model.Organization) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.organizationsByTenantCache[tenant] = freecache.NewCache(cache500MB)

	const chunkSize = 100 // Size of each domain chunk
	for i, j := 0, chunkSize; i < len(organizations); i, j = i+chunkSize, j+chunkSize {
		// This ensures we don't go past the end of the slice
		if j > len(organizations) {
			j = len(organizations)
		}

		// Get the current chunk and marshal it
		organizationsChunk := organizations[i:j]
		organizationsChunkBytes, err := json.Marshal(organizationsChunk)
		if err != nil {
			return
		}

		// Generate a key based on the index
		key := strconv.Itoa(i/chunkSize + 1) // Convert the integer to a string

		// Store the chunk in the cache
		err = c.organizationsByTenantCache[tenant].Set([]byte(key), organizationsChunkBytes, expire9999Days)
		if err != nil {
			return
		}
	}
}

func (c *Cache) GetOrganizations(tenant string) []*model.Organization {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var all []*model.Organization
	keyIndex := 1

	for {
		// Generate the key based on index
		key := strconv.Itoa(keyIndex)

		// Attempt to get the domains chunk from the cache
		domainChunkBytes, err := c.organizationsByTenantCache[tenant].Get([]byte(key))
		if err != nil {
			break // If a key is not found, assume no more chunks are available
		}

		var domainChunk []*model.Organization
		err = json.Unmarshal(domainChunkBytes, &domainChunk)
		if err != nil {
			// If there is an error unmarshalling, decide how to handle it
			// For simplicity, we stop and return what we have so far
			break
		}

		// Append this chunk of domains to the allDomains slice
		all = append(all, domainChunk...)

		keyIndex++ // Increment key index for next iteration
	}

	return all // Return the combined list of all domains
}
