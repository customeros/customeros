package server

import (
	"context"

	commoncaches "github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/clients/grpc_client"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	rest_handlers "github.com/customeros/customeros/packages/server/customer-os-api/rest"
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

func RegisterRestRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerInternalRoutes(ctx, r, s, h)
	registerPublicRoutes(ctx, r, s, h)
	registerFileRoutes(ctx, r, s, h)

	registerBillingRoutes(ctx, r, s, h)
	registerCustomerBaseRoutes(ctx, r, s, h)
	registerEnrichRoutes(ctx, r, s, h)
	registerFlowRoutes(ctx, r, s, h)
	registerIDRoutes(ctx, r, s, h)
	registerMailstackRoutes(ctx, r, s, h)
	registerOutreachRoutes(ctx, r, s, h)
	registerRevealRoutes(ctx, r, s, h)
	registerVerifyRoutes(ctx, r, s, h)
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
			rest_handlers.StatsSuccessHandler(config.method+":"+config.path, config.services),
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
