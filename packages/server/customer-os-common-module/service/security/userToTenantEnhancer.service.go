package security

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"net/http"
)

type HeaderAllowance string

const (
	USERNAME           HeaderAllowance = "USERNAME"
	TENANT             HeaderAllowance = "TENANT"
	USERNAME_OR_TENANT HeaderAllowance = "USERNAME_OR_TENANT"
)

const (
	KEY_TENANT_NAME = "TenantName"
	KEY_USER_ID     = "UserId"
	KEY_USER_EMAIL  = "UserEmail"
	KEY_USER_ROLES  = "UserRoles"
	KEY_APP_SOURCE  = "AppSource"
)

const UsernameHeader = "X-Openline-USERNAME"
const TenantHeader = "X-Openline-TENANT"

func TenantUserContextEnhancer(headerAllowance HeaderAllowance, cr *neo4jrepository.Repositories, opts ...CommonServiceOption) func(c *gin.Context) {
	// Apply the options to configure the middleware
	config := &Options{}
	for _, opt := range opts {
		opt(config)
	}

	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "TenantUserContextEnhancer")
		spanFinished := false
		defer func() {
			if !spanFinished {
				span.Finish()
			}
		}()

		//if API call is made with a tenant api key, the apiKeyChecker validated it already
		tenantKh := c.GetHeader(TenantApiKeyHeader)
		if tenantKh != "" {
			return
		}

		tenantHeader := c.GetHeader(TenantHeader)
		usernameHeader := c.GetHeader(UsernameHeader)
		span.LogFields(
			log.String("header.tenant", tenantHeader),
			log.String("header.username", usernameHeader),
			log.String("headerAllowance", string(headerAllowance)))

		var (
			tenantExists bool
			userId       string
			tenantName   string
			roles        []string
			err          error
		)

		switch headerAllowance {
		case TENANT:
			tenantExists, err = checkTenantHeader(c, tenantHeader, cr, ctx, config.cache)
			if err != nil {
				return
			}
			c.Set(KEY_TENANT_NAME, tenantHeader)

		case USERNAME:
			userId, tenantName, roles, err = checkUsernameHeader(c, usernameHeader, cr, ctx, config.cache)
			if err != nil {
				return
			}
			c.Set(KEY_TENANT_NAME, tenantName)
			c.Set(KEY_USER_ID, userId)
			c.Set(KEY_USER_EMAIL, usernameHeader)
			c.Set(KEY_USER_ROLES, roles)

		case USERNAME_OR_TENANT:
			if tenantHeader == "" && usernameHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message": "User or Tenant must be specified"}},
				})
				c.Abort()
				return
			}

			if tenantHeader != "" {
				tenantExists, err = checkTenantHeader(c, tenantHeader, cr, ctx, config.cache)
				if err != nil {
					return
				}
				if tenantExists {
					c.Set(KEY_TENANT_NAME, tenantHeader)
				}
			}
			if usernameHeader != "" {
				userId, tenantName, roles, err = checkUsernameHeader(c, usernameHeader, cr, ctx, config.cache)
				if err != nil {
					return
				}
				c.Set(KEY_TENANT_NAME, tenantName)
				c.Set(KEY_USER_ID, userId)
				c.Set(KEY_USER_EMAIL, usernameHeader)
				c.Set(KEY_USER_ROLES, roles)
			}
			if !spanFinished {
				span.Finish()
				spanFinished = true
			}
			c.Next()
			return
		default:
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []gin.H{{"message": "unknown header Allowance"}},
			})
			c.Abort()
			return
		}

		if !spanFinished {
			span.Finish()
			spanFinished = true
		}
		c.Next()
	}
}

func checkTenantHeader(c *gin.Context, tenantHeader string, cr *neo4jrepository.Repositories, ctx context.Context, cache *caches.Cache) (bool, error) {
	if tenantHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": "missing tenant header"}},
		})
		c.Abort()
		return false, fmt.Errorf("missing tenant header")
	}

	if cache != nil && cache.CheckTenant(tenantHeader) {
		return true, nil
	}
	exists, err := cr.TenantReadRepository.TenantExists(ctx, tenantHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": fmt.Sprintf("failed to check tenant existence: %v", err)}},
		})
		c.Abort()
		return false, fmt.Errorf("failed to check tenant existence: %v", err)
	}
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": "tenant does not exist"}},
		})
		c.Abort()
		return false, fmt.Errorf("tenant does not exist")
	}

	if cache != nil {
		cache.AddTenant(tenantHeader)
	}
	return true, nil
}

func checkUsernameHeader(c *gin.Context, usernameHeader string, cr *neo4jrepository.Repositories, ctx context.Context, cache *caches.Cache) (string, string, []string, error) {
	if usernameHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": "missing username header"}},
		})
		c.Abort()
		return "", "", []string{}, fmt.Errorf("missing username header")
	}

	if cache != nil {
		userId, tenantName, roles, found := cache.GetUserDetailsFromCache(usernameHeader)
		if found {
			return userId, tenantName, roles, nil
		}
	}
	userId, tenantName, roles, err := cr.UserReadRepository.FindFirstUserWithRolesByEmail(ctx, usernameHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": fmt.Sprintf("failed to find user: %v", err)}},
		})
		c.Abort()
		return "", "", []string{}, fmt.Errorf("failed to find user: %v", err)
	}
	if tenantName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"message": "user has no associated tenant"}},
		})
		c.Abort()
		return "", "", []string{}, fmt.Errorf("user has no associated tenant")
	}

	if cache != nil {
		cache.AddUserDetailsToCache(usernameHeader, userId, tenantName, roles)
	}

	return userId, tenantName, roles, nil
}
