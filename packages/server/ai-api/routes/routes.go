package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/ai-api/service"
)

type Route struct {
	Method      string
	Path        string
	TracingName string
	Handler     gin.HandlerFunc
}

// routes defines all available API routes
func routes(services *service.Services) []Route {
	return []Route{
		{
			Method:      "GET",
			Path:        "/query",
			TracingName: "GET /query",
			Handler:     query(services),
		},
	}
}

// RegisterRoutes registers all API routes and their middleware
func RegisterRoutes(ctx context.Context, r *gin.Engine, services *service.Services) {
	// Register each route with its middleware
	for _, route := range routes(services) {
		handlers := append(
			getCommonMiddleware(ctx, services, route.TracingName),
			route.Handler,
		)
		r.Handle(route.Method, route.Path, handlers...)
	}
}

// getCommonMiddleware returns the middleware stack used by all routes
func getCommonMiddleware(ctx context.Context, services *service.Services, tracingName string) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		tracing.TracingEnhancer(ctx, tracingName),
		security.ApiKeyCheckerHTTP(
			services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository,
			services.CommonServices.PostgresRepositories.AppKeyRepository,
			security.AI_API,
			security.WithCache(services.CommonServices.Cache),
		),
	}
}
