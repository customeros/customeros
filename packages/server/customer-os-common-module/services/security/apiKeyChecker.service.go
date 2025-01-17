package security

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type App string

const (
	CUSTOMER_OS_API      App = "customer-os-api"
	CUSTOMER_OS_WEBHOOKS App = "customer-os-webhooks"
	MAILSHEPRA_API       App = "mailsherpa-api"
	PLATFORM_ADMIN_API   App = "platform-admin-api"
)

const (
	ApiKeyHeader       = "X-Openline-API-KEY"
	TenantApiKeyHeader = "X-CUSTOMER-OS-API-KEY"
)

func ApiKeyCheckerHTTP(tenantApiKeyRepo postgresRepository.TenantWebhookApiKeyRepository, appKeyRepo postgresRepository.AppKeyRepository, app App, opts ...CommonServiceOption) func(c *gin.Context) {
	// Apply the options to configure the middleware
	config := &Options{}
	for _, opt := range opts {
		opt(config)
	}

	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "ApiKeyCheckerHTTP")
		spanFinished := false
		defer func() {
			if !spanFinished {
				span.Finish()
			}
		}()
		span.LogFields(log.String("app", string(app)))

		kh := c.GetHeader(ApiKeyHeader)
		tenantKh := c.GetHeader(TenantApiKeyHeader)
		if kh != "" {
			// Check if the API key matches the cached value
			if config.cache != nil && config.cache.CheckApiKey(string(app), kh) {
				// Valid API key found in cache
				span.LogFields(log.Bool("cached", true))
				if !spanFinished {
					spanFinished = true
					span.Finish()
				}
				c.Next()
				return
			}
			span.LogFields(log.Bool("cached", false))
			appKey, err := appKeyRepo.FindByKey(ctx, string(app), kh)
			if err != nil {

				c.JSON(http.StatusUnauthorized, gin.H{
					"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
					"status":    "error",
					"message":   fmt.Sprintf("error while checking api key: %s", err.Error()),
				})
				c.Abort()
				return
			}

			if appKey == nil {

				c.JSON(http.StatusUnauthorized, gin.H{
					"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
					"status":    "error",
					"message":   "invalid API key",
				})
				c.Abort()
				return
			}

			// If the API key is valid after database check, cache it
			if config.cache != nil && appKey != nil {
				config.cache.SetApiKey(string(app), kh)
			}

			if !spanFinished {
				spanFinished = true
				span.Finish()
			}
			c.Next()
		} else if tenantKh != "" {
			// Check if the API key matches the cached value
			if config.cache != nil && config.cache.CheckTenantApiKey(tenantKh) {
				// Valid API key found in cache
				span.LogFields(log.Bool("cached", true))
				if !spanFinished {
					spanFinished = true
					span.Finish()
				}
				c.Next()
				return
			}
			span.LogFields(log.Bool("cached", false))

			apiKey, err := tenantApiKeyRepo.GetTenantForApiKey(ctx, tenantKh)
			if err != nil || apiKey == nil {

				c.JSON(http.StatusUnauthorized, gin.H{
					"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
					"status":    "error",
					"message":   "invalid API key",
				})
				c.Abort()
				return
			}

			if !apiKey.Enabled {

				c.JSON(http.StatusUnauthorized, gin.H{
					"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
					"status":    "error",
					"message":   "invalid API key disabled",
				})
				c.Abort()
				return
			}

			if config.cache != nil {
				config.cache.SetTenantApiKey(tenantKh, apiKey.Tenant)
			}

			c.Set(KEY_TENANT_NAME, apiKey.Tenant)
			c.Set(KEY_USER_ROLES, []string{"USER"})

			if !spanFinished {
				spanFinished = true
				span.Finish()
			}
			c.Next()
		} else {
			// illegal request, terminate the current process

			c.JSON(http.StatusUnauthorized, gin.H{
				"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
				"status":    "error",
				"message":   "API key is required",
			})
			span.LogFields(log.String("result", "Missing api key"))
			c.Abort()
			return
		}
	}
}
