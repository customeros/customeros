package caches

import (
	"encoding/json"
	"github.com/coocood/freecache"
)

const (
	KB       = 1024
	cache5MB = 5 * 1024 * KB
)
const (
	expire20Days = 20 * 24 * 60 * 60 // 20 days
)

type UserType string

const UserType_User UserType = "user"
const UserType_Contact UserType = "contact"
const UserType_NonUser UserType = "non-user"

type SlackUser struct {
	UserType UserType `json:"type,omitempty"`
	Name     string   `json:"name,omitempty"`
}

type Cache interface {
	SetSlackUser(tenant, userId string, user SlackUser)
	GetSlackUser(tenant, userId string) (SlackUser, bool)
	SetSlackUserAsContactForOrg(orgId, userId, value string)
	GetSlackUserAsContactForOrg(orgId, userId string) (string, bool)
}

type cache struct {
	usersCache *freecache.Cache
}

func InitCaches() Cache {
	result := cache{
		usersCache: freecache.NewCache(cache5MB),
	}

	return &result
}

// User cache
func (c *cache) SetSlackUser(tenant, userId string, user SlackUser) {
	// Convert strings to []byte
	keyBytes := []byte(tenant + "-" + userId)
	userJson, err := json.Marshal(user)
	if err == nil {
		valueBytes := userJson
		_ = c.usersCache.Set(keyBytes, valueBytes, expire20Days)
	}
}

func (c *cache) GetSlackUser(tenant, userId string) (SlackUser, bool) {
	userJson, found := c.get(c.usersCache, tenant+"-"+userId)
	if !found {
		return SlackUser{}, false
	}
	slackUser := SlackUser{}
	err := json.Unmarshal([]byte(userJson), &slackUser)
	if err == nil {
		return slackUser, true
	}
	return SlackUser{}, false
}

// User as contact cache
func (c *cache) SetSlackUserAsContactForOrg(orgId, userId, value string) {
	// Convert strings to []byte
	keyBytes := []byte(orgId + "-" + userId)
	valueBytes := []byte(value)

	_ = c.usersCache.Set(keyBytes, valueBytes, expire20Days)
}

func (c *cache) GetSlackUserAsContactForOrg(orgId, userId string) (string, bool) {
	return c.get(c.usersCache, orgId+"-"+userId)
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
