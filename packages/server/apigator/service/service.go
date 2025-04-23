package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coocood/freecache"
	entities "github.com/customeros/customeros/packages/server/apigator/entity"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
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

func (s *Service) GetTenantByUser(ctx context.Context, username string) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "Service.GetTenantByUser")
	defer spans.Finish()
	spans.LogKV("username", username)

	foundTenant, err := s.userRepo.GetCurrentTenantByUserEmail(s.ctx, username)
	if err != nil {
		err = fmt.Errorf("failed to get current tenant by user email: %w", err)
		spans.TraceError(err)
		return "", err
	}

	spans.LogKV("result.foundTenant", foundTenant)
	return foundTenant, err
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

	key := []byte(fmt.Sprintf("%s:%s", tenant, username))

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
		unmarshalErr := json.Unmarshal(cached, &result)
		if unmarshalErr != nil {
			return nil, fmt.Errorf("failed to unmarshal user details: %w", unmarshalErr)
		}
	}

	return result, err
}
