package security

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

const (
	KEY_TENANT_NAME           = "TenantName"
	KEY_AUTHENTICATED_USER_ID = "AuthenticatedUserId"
	KEY_USER_ID               = "UserId"
	KEY_USER_EMAIL            = "UserEmail"
	KEY_USER_ROLES            = "UserRoles"
)

const UsernameHeader = "X-Openline-USERNAME"
const TenantHeader = "X-Openline-TENANT"

func TenantUserContextEnhancer(cr *neo4jrepository.Repositories, opts ...CommonServiceOption) func(c *gin.Context) {
	// Apply the options to configure the middleware
	config := &Options{}
	for _, opt := range opts {
		opt(config)
	}

	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "TenantUserContextEnhancer")
		defer span.Finish()

		//if API call is made with a tenant api key, the apiKeyChecker validated it already
		tenantKh := c.GetHeader(TenantApiKeyHeader)
		if tenantKh != "" {
			return
		}

		tenantHeader := c.GetHeader(TenantHeader)
		usernameHeader := c.GetHeader(UsernameHeader)
		span.LogFields(
			log.String("header.tenant", tenantHeader),
			log.String("header.username", usernameHeader))

		if tenantHeader == "" && usernameHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []gin.H{{"message": "X-Openline-USERNAME AND X-Openline-TENANT must be specified"}},
			})
			c.Abort()
			return
		}

		//TODO remove this after 01.03.2025
		//fallback for missing tenant header
		// should work for all customers until 01.03.2025. should be removed after that
		if tenantHeader == "" {
			allUsers, err := cr.UserReadRepository.FindAllUsersWithRolesByEmail(ctx, usernameHeader)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message": fmt.Sprintf("failed to find user: %v", err)}},
				})
				c.Abort()
				return
			}

			if allUsers == nil || len(allUsers) != 1 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message": fmt.Sprintf("failed to find user: %v", err)}},
				})
				c.Abort()
				return
			}

			tenantHeader = allUsers[0].Tenant
		}

		authenticatedUserInTenant, err := checkUsernameHeader(c, tenantHeader, usernameHeader, cr, ctx, config.cache)
		if err != nil {
			return
		}

		c.Set(KEY_TENANT_NAME, authenticatedUserInTenant.Tenant)
		c.Set(KEY_AUTHENTICATED_USER_ID, authenticatedUserInTenant.AuthenticatedUserId)
		c.Set(KEY_USER_ID, authenticatedUserInTenant.UserId)
		c.Set(KEY_USER_EMAIL, usernameHeader)
		c.Set(KEY_USER_ROLES, authenticatedUserInTenant.Roles)

		c.Next()
	}
}

func checkUsernameHeader(c *gin.Context, tenant, username string, cr *neo4jrepository.Repositories, ctx context.Context, cache *caches.Cache) (*neo4jrepository.AuthenticatedUserInTenant, error) {
	if cache != nil {
		userDetails, found := cache.GetUserDetailsFromCache(tenant, username)
		if found {
			return userDetails, nil
		}
	}
	authenticatedUserInTenant, err := cr.UserReadRepository.FindFirstUserWithRolesByEmail(ctx, tenant, username)
	if err != nil || authenticatedUserInTenant == nil || authenticatedUserInTenant.UserId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": fmt.Sprintf("failed to find user: %v", err)}},
		})
		c.Abort()
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	if cache != nil {
		cache.AddUserDetailsToCache(tenant, username, authenticatedUserInTenant)
	}

	return authenticatedUserInTenant, nil
}
