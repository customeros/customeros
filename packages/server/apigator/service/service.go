package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coocood/freecache"
	entities "github.com/customeros/customeros/packages/server/apigator/entity"
	neo4jRepos "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type ApiKeyRepository interface {
	GetTenantForApiKey(ctx context.Context, apiKey string) (*postgres_entity.TenantWebhookApiKey, error)
}
type UserRepository interface {
	GetCurrentTenantByUserEmail(ctx context.Context, email string) (string, error)
	FindFirstUserWithRolesByEmail(ctx context.Context, tenant string, email string) (*neo4jRepos.AuthenticatedUserInTenant, error)
}

type Cache struct {
	tenantsByUsers   *freecache.Cache
	tenantsByApiKeys *freecache.Cache
	userDetails      *freecache.Cache
}

type Service struct {
	ctx        context.Context
	cache      *Cache
	apiKeyRepo ApiKeyRepository
	userRepo   UserRepository
}

func (s *Service) Init(
	apiKeyRepo ApiKeyRepository,
	userRepo UserRepository,
) *Service {
	s.cache = &Cache{
		tenantsByUsers:   freecache.NewCache(5 * 1024 * 1024),
		tenantsByApiKeys: freecache.NewCache(5 * 1024 * 1024),
		userDetails:      freecache.NewCache(5 * 1024 * 1024),
	}
	s.apiKeyRepo = apiKeyRepo
	s.userRepo = userRepo

	return s
}

func (s *Service) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *Service) GetTenantByUser(username string) (string, error) {
	var result string
	var err error

	cached, cacheErr := s.cache.tenantsByUsers.Get([]byte(username))
	if cacheErr != nil {
		retrieved, err := s.userRepo.GetCurrentTenantByUserEmail(s.ctx, username)
		if err != nil {
			err = fmt.Errorf("failed to get current tenant by user email: %w", err)
		}
		if retrieved != "" {
			result = retrieved
			s.cache.tenantsByUsers.Set([]byte(username), []byte(retrieved), 24*60*60)
		}
	} else {
		result = string(cached)
	}

	return result, err
}

func (s *Service) GetTenantByApiKey(apiKey string) (string, error) {
	var result string
	var err error

	cached, cacheErr := s.cache.tenantsByApiKeys.Get([]byte(apiKey))
	if cacheErr != nil {
		retrieved, err := s.apiKeyRepo.GetTenantForApiKey(s.ctx, apiKey)
		if err != nil {
			err = fmt.Errorf("failed to get tenant by api key: %w", err)
		}
		if retrieved != nil {
			result = retrieved.Tenant
			s.cache.tenantsByApiKeys.Set([]byte(apiKey), []byte(result), 24*60*60)
		}
	} else {
		result = string(cached)
	}

	return result, err
}

func (s *Service) GetUserDetails(tenant string, username string) (*entities.UserDetails, error) {
	var result *entities.UserDetails
	var err error

	key := fmt.Appendf(nil, "%s:%s", tenant, username)

	cached, cacheErr := s.cache.userDetails.Get(key)
	if cacheErr != nil {
		retrieved, err := s.userRepo.FindFirstUserWithRolesByEmail(s.ctx, tenant, username)
		if err != nil {
			err = fmt.Errorf("failed to find user with roles by email: %w", err)
		}
		if retrieved != nil {
			result = &entities.UserDetails{
				Id:        retrieved.UserId,
				Email:     retrieved.UserPrimaryEmail,
				FirstName: retrieved.UserFirstname,
				LastName:  retrieved.UserLastname,
				Roles:     retrieved.Roles,
			}
			marshaled, err := result.Marshal()
			if err != nil {
				return nil, fmt.Errorf("failed to marshal user details: %w", err)
			}
			s.cache.userDetails.Set(key, marshaled, 24*60*60)
		}
	} else {
		err := json.Unmarshal(cached, &result)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal user details: %w", err)
		}
	}

	return result, err
}
