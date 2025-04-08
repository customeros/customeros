package security

import (
	"net/http"

	postgresRepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

const (
	ApiKeyHeader       = "X-Openline-API-KEY"
	TenantApiKeyHeader = "X-CUSTOMER-OS-API-KEY"
)

func ApiKeyCheckerHTTP(
	tenantApiKeyRepo postgresRepository.TenantWebhookApiKeyRepository,
	appKey string,
	opts ...CommonServiceOption,
) func(c *gin.Context) {
	// Apply the options to configure the middleware
	config := &Options{}
	for _, opt := range opts {
		opt(config)
	}

	return func(c *gin.Context) {
		spans, ctx := telemetry.StartSpan(c.Request.Context(), "ApiKeyCheckerHTTP")
		spanFinished := false
		defer func() {
			if !spanFinished {
				spans.Finish()
			}
		}()

		kh := c.GetHeader(ApiKeyHeader)
		tenantKh := c.GetHeader(TenantApiKeyHeader)
		if kh != "" {
			if appKey != kh {
				spans.LogKV("result", "Invalid app API key")
				c.JSON(http.StatusUnauthorized, gin.H{
					"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
					"status":    "error",
					"message":   "invalid API key",
				})
				c.Abort()
				return
			}
			spans.LogKV("result", "Valid app API key")
			tenant := c.GetHeader(TenantHeader)
			c.Set(KEY_TENANT_NAME, tenant)
			c.Set(KEY_USER_ROLES, []string{"USER"})
			spanFinished = true
			spans.Finish()
			c.Next()
		} else if tenantKh != "" {
			// Check if the API key matches the cached value
			if config.cache != nil && config.cache.CheckTenantApiKey(tenantKh) {
				// Valid API key found in cache
				spans.LogKV("result", "Valid tenant API key from cache")
				if !spanFinished {
					spanFinished = true
					spans.Finish()
				}
				c.Next()
				return
			}
			spans.LogKV("cached", false)

			apiKey, err := tenantApiKeyRepo.GetTenantForApiKey(ctx, tenantKh)
			if err != nil || apiKey == nil {
				spans.LogKV("result", "Invalid tenant API key")
				c.JSON(http.StatusUnauthorized, gin.H{
					"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
					"status":    "error",
					"message":   "invalid API key",
				})
				c.Abort()
				return
			}

			if !apiKey.Enabled {
				spans.LogKV("result", "Disabled tenant API key")
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

			spans.LogKV("result", "Valid tenant API key")
			if !spanFinished {
				spanFinished = true
				spans.Finish()
			}
			c.Next()
		} else {
			// illegal request, terminate the current process
			spans.LogKV("result", "Missing API key")
			c.JSON(http.StatusUnauthorized, gin.H{
				"requestId": utils.GenerateNanoIdWithPrefix("api", 16),
				"status":    "error",
				"message":   "API key is required",
			})
			c.Abort()
			return
		}
	}
}
