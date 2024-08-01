package security

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

type App string

const (
	CUSTOMER_OS_API      App = "customer-os-api"
	CUSTOMER_OS_WEBHOOKS App = "customer-os-webhooks"
	FILE_STORE_API       App = "file-store-api"
	SETTINGS_API         App = "settings-api"
	VALIDATION_API       App = "validation-api"
	AI_API               App = "ai-api"
	PLATFORM_ADMIN_API   App = "platform-admin-api"
	USER_ADMIN_API       App = "user-admin-api"
)

const ApiKeyHeader = "X-Openline-API-KEY"
const TenantApiKeyHeader = "X-CUSTOMER-OS-API-KEY"

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
			keyResult := appKeyRepo.FindByKey(ctx, string(app), kh)

			if keyResult.Error != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message": fmt.Sprintf("Error while checking api key: %s", keyResult.Error.Error())}},
				})
				c.Abort()
				return
			}

			appKey := keyResult.Result.(*postgresEntity.AppKey)

			if appKey == nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message": "Invalid api key"}},
				})
				c.Abort()
				return
			}

			// If the API key is valid after database check, cache it
			if config.cache != nil && keyResult.Result != nil {
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
					"errors": []gin.H{{"message”": "Invalid api key"}},
				})
				c.Abort()
				return
			}

			if !apiKey.Enabled {
				c.JSON(http.StatusUnauthorized, gin.H{
					"errors": []gin.H{{"message”": "Api key disabled"}},
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
				"errors": []gin.H{{"message": "Api key is required"}},
			})
			tracing.TraceErr(span, errors.New("Api key is required"))
			c.Abort()
			return
		}
	}
}

func ApiKeyCheckerGRPC(ctx context.Context, appKeyRepo postgresRepository.AppKeyRepository, app App) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	kh := md.Get(ApiKeyHeader)
	if len(kh) == 1 {
		keyResult := appKeyRepo.FindByKey(ctx, string(app), kh[0])
		if keyResult.Error != nil {
			return false
		}
		appKey := keyResult.Result.(*postgresEntity.AppKey)
		return appKey != nil
	}
	return false
}
