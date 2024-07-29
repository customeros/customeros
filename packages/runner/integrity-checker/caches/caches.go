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
	mu                    sync.RWMutex
	previousAlertMessages *freecache.Cache
}

func NewCache() *Cache {
	cache := Cache{
		previousAlertMessages: freecache.NewCache(cache1MB),
	}
	return &cache
}

func (c *Cache) SetPreviousAlertMessages(results []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(results)
	if err != nil {
		return err
	}

	err = c.previousAlertMessages.Set([]byte("previousAlertMessages"), data, expire48HoursInSeconds)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) GetPreviousAlertMessages() ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := c.previousAlertMessages.Get([]byte("previousAlertMessages"))
	if err != nil {
		return nil, err
	}

	var results []string
	err = json.Unmarshal(data, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
