package routes

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/service"
)

const (
	SCRAPIN    string = "SCRAPIN"
	BRANDFETCH string = "BRANDFETCH"
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
			Path:        "/enrichPerson",
			TracingName: "GET /enrichPerson",
			Handler:     enrichPerson(services),
		},
		{
			Method:      "GET",
			Path:        "/findWorkEmail",
			TracingName: "GET /findWorkEmail",
			Handler:     findWorkEmail(services),
		},
		{
			Method:      "GET",
			Path:        "/enrichOrganization",
			TracingName: "GET /enrichOrganization",
			Handler:     enrichOrganization(services),
		},
		{
			Method:      "GET",
			Path:        "/scrapinOrganization",
			TracingName: "GET /enrichOrganization",
			Handler:     scrapinOrganization(services),
		},
		{
			Method:      "GET",
			Path:        "/snitcher",
			TracingName: "GET /snitcher",
			Handler:     snitcherData(services),
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
			security.ENRICHMENT_API,
			security.WithCache(services.CommonServices.Cache),
		),
	}
}
