package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
	"net/http"
)

type App string

const (
	CUSTOMER_OS_API      App = "customer-os-api"
	CUSTOMER_OS_WEBHOOKS App = "customer-os-webhooks"
	FILE_STORE_API       App = "file-store-api"
	SETTINGS_API         App = "settings-api"
	VALIDATION_API       App = "validation-api"
	ANTHROPIC_API        App = "anthropic-api"
	OPENAI_API           App = "openai-api"
)

const ApiKeyHeader = "X-Openline-API-KEY"

func ApiKeyCheckerHTTP(appKeyRepo repository.AppKeyRepository, app App, opts ...CommonServiceOption) func(c *gin.Context) {
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

			appKey := keyResult.Result.(*entity.AppKey)

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
			// illegal request, terminate the current process
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errors": []gin.H{{"message": "Api key is required"}},
			})
			tracing.TraceErr(span, errors.New("Api key is required"))
			c.Abort()
			return
		}
	}
}

func ApiKeyCheckerGRPC(ctx context.Context, appKeyRepo repository.AppKeyRepository, app App) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	kh := md.Get(ApiKeyHeader)
	if kh != nil && len(kh) == 1 {
		keyResult := appKeyRepo.FindByKey(ctx, string(app), kh[0])
		if keyResult.Error != nil {
			return false
		}
		appKey := keyResult.Result.(*entity.AppKey)
		if appKey == nil {
			return false
		}
		return true
	}
	return false
}
