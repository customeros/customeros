package security

import (
	"net/http"

	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
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

func ApiKeyCheckerHTTP(tenantApiKeyRepo postgresRepository.TenantWebhookApiKeyRepository, appKey string, opts ...CommonServiceOption) func(c *gin.Context) {
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

		kh := c.GetHeader(ApiKeyHeader)
		tenantKh := c.GetHeader(TenantApiKeyHeader)
		if kh != "" {
			if appKey != kh {
				c.JSON(http.StatusUnauthorized, gin.H{
					"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
					"status":    "error",
					"message":   "invalid API key",
				})
				c.Abort()
				return
			}
			spanFinished = true
			span.Finish()
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
