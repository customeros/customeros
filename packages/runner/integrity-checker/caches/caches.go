package caches

import (
	"encoding/json"
	"github.com/coocood/freecache"
	"sync"
)

const (
	KB       = 1024
	cache1MB = 1 * 1024 * KB
)
const (
	expire48HoursInSeconds = 48 * 60 * 60
)

type Cache struct {
	mu                            sync.RWMutex
	previousNeo4jAlertMessages    *freecache.Cache
	previousPostgresAlertMessages *freecache.Cache
}

func NewCache() *Cache {
	cache := Cache{
		previousNeo4jAlertMessages:    freecache.NewCache(cache1MB),
		previousPostgresAlertMessages: freecache.NewCache(cache1MB),
	}
	return &cache
}

func (c *Cache) SetPreviousNeo4jAlertMessages(results []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(results)
	if err != nil {
		return err
	}

	err = c.previousNeo4jAlertMessages.Set([]byte("previousNeo4jAlertMessages"), data, expire48HoursInSeconds)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetPreviousNeo4jAlertMessages() ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := c.previousNeo4jAlertMessages.Get([]byte("previousNeo4jAlertMessages"))
	if err != nil {
		// Record not found, return empty slice
		return []string{}, nil
	}

	var results []string
	err = json.Unmarshal(data, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Cache) SetPreviousPostgresAlertMessages(results []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(results)
	if err != nil {
		return err
	}

	err = c.previousPostgresAlertMessages.Set([]byte("previousPostgresAlertMessages"), data, expire48HoursInSeconds)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetPreviousPostgresAlertMessages() ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := c.previousPostgresAlertMessages.Get([]byte("previousPostgresAlertMessages"))
	if err != nil {
		// Record not found, return empty slice
		return []string{}, nil
	}

	var results []string
	err = json.Unmarshal(data, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
