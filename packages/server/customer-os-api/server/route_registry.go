package server

import (
	"context"

	"github.com/gin-gonic/gin"
	commoncaches "github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	cosHandler "github.com/customeros/customeros/packages/server/customer-os-api/rest"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type RouteType string

const (
	RoutePublic   RouteType = "public"
	RouteInternal RouteType = "internal"
	RouteCustomer RouteType = "customer"
	RouteFiles    RouteType = "files"
)

type RouteConfig struct {
	method     string
	path       string
	handler    gin.HandlerFunc
	routeType  RouteType
	services   *cosapi_services.Services
	cache      *commoncaches.Cache
	grpcClient *grpc_client.Clients
}

func RegisterRestRoutes(ctx context.Context, r *gin.Engine, grpcClients *grpc_client.Clients, s *cosapi_services.Services, cache *commoncaches.Cache) {
	registerInternalRoutes(ctx, r, s)
	registerPublicRoutes(ctx, r, s)
	registerFileRoutes(ctx, r, s)

	registerBillingRoutes(ctx, r, s, cache)
	registerCustomerBaseRoutes(ctx, r, s, cache)
	registerEnrichRoutes(ctx, r, s, cache)
	registerFlowRoutes(ctx, r, s, cache)
	registerIDRoutes(ctx, r, s)
	registerMailstackRoutes(ctx, r, s, cache)
	registerOutreachRoutes(ctx, r, s, cache)
	registerRevealRoutes(ctx, r, s, cache)
	registerVerifyRoutes(ctx, r, s, cache)
}

func registerRoute(ctx context.Context, r *gin.Engine, config RouteConfig) {
	middlewares := []gin.HandlerFunc{
		tracing.TracingEnhancer(ctx, config.method+":"+config.path),
	}

	switch config.routeType {
	case RouteInternal:
		middlewares = append(middlewares,
			security.ApiKeyCheckerHTTP(
				config.services.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository,
				config.services.Repositories.PostgresRepositories.AppKeyRepository,
				security.PLATFORM_ADMIN_API,
				security.WithCache(config.services.Cache),
			))

	case RouteCustomer:
		middlewares = append(middlewares,
			security.ApiKeyCheckerHTTP(
				config.services.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository,
				config.services.Repositories.PostgresRepositories.AppKeyRepository,
				security.CUSTOMER_OS_API,
				security.WithCache(config.cache),
			),
			enrichContextMiddleware(constants.AppSourceCustomerOsApiRest),
			cosHandler.StatsSuccessHandler(config.method+":"+config.path, config.services),
		)

	case RouteFiles:
		middlewares = append(middlewares,
			config.services.JWTService.GetJWTTenantUserEnhancer(),
			security.TenantUserContextEnhancer(
				security.USERNAME_OR_TENANT,
				config.services.Repositories.Neo4jRepositories,
				security.WithCache(config.cache),
			))

	}

	middlewares = append(middlewares, config.handler)
	r.Handle(config.method, config.path, middlewares...)
}
